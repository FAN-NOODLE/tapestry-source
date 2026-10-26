<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2019/11/15
 * Time: 2:11 PM
 */

namespace App\StateMachine\ActivityCenter\DiscountCard;

use App\Models\ActivityStatus;
use App\Models\ActivityTicket;
use App\Models\ActivityTicketLog;
use App\Models\Error;
use App\StateMachine\ActivityCenter\Machine;
use Illuminate\Support\Facades\DB;

class TicketTranstion
{
    const ACTION_CREATE = 'create';  // 生成订单
    const ACTION_CANCEL = 'cancel'; // 取消订单
    const ACTION_OVER = 'over'; // 过期
    const ACTION_PAY = 'pay'; // 支付订单

    private static function transtion()
    {
        $transtion = [
            ActivityTicket::STATUS_INIT => [
                self::ACTION_CREATE => ActivityTicket::STATUS_WAITING,
                self::ACTION_OVER   => ActivityTicket::STATUS_OVER,
            ],

            ActivityTicket::STATUS_WAITING => [
                self::ACTION_PAY    => ActivityTicket::STATUS_PAYED,
                self::ACTION_CANCEL => ActivityTicket::STATUS_CANCELED,
            ],
        ];

        return $transtion;
    }

    /**
     * @param $enterpriseId int 企业id
     * @param $id int 主键
     * @param $action string 动作
     * @param $to
     * @param $extAdd
     *
     * @return mixed
     */
    public static function run($enterpriseId, $id, $action, $to, $extAdd = [])
    {
        $transtions = self::transtion();

        $history = ActivityStatus::getHistory($enterpriseId, 'DiscountCardTicket', $id);

        $fsm = Machine::create($transtions, $history);

        $assertEqualsTo = $fsm->doAction($action);

        // 成功
        if ($to != $assertEqualsTo) {
            notice('group_initiate_status_abnormal');
        }
        try {
            DB::beginTransaction();
            // 保存历史
            ActivityStatus::setHistory($enterpriseId, 'DiscountCardTicket', $id, $fsm->getHistory());

            ActivityTicket::query()->where('id', $id)->update(['status' => $to]);

            $insertLog = [
                'activity_ticket_id' => $id,
                'status' => $to,
                'enterprise_id' => $enterpriseId,
            ];
            $insertLog = array_merge($insertLog, $extAdd);

            ActivityTicketLog::query()->create($insertLog);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            Error::log('更新票据状态', $e);
            notice('machine_error');
        }

        return $to;
    }
}
