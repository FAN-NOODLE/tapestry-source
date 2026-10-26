<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2019/11/15
 * Time: 2:11 PM
 */

namespace App\StateMachine\ActivityCenter\Haggle;

use App\Caches\ActivityCenter\Common\ActivityInitiateCache;
use App\Models\ActivityCenter\Haggle\HaggleInitiate;
use App\Models\ActivitySku;
use App\Models\ActivityStatus;
use App\Models\Error;
use App\StateMachine\ActivityCenter\Machine;
use Illuminate\Support\Facades\DB;

class InitiateTranstion
{
    const ACTION_PAY = 'pay'; // 支付
    const ACTION_CANCEL = 'cancel'; // 取消
    const ACTION_CREATE = 'create'; // 生成
    const ACTION_DOWN = 'down'; // 下架
    const ACTION_OVER = 'over'; // 过期
    const ACTION_LAST_TIMES = 'last_times'; // 最后一刀
    const ACTION_TICKET = 'ticket'; // 生成票据
    const ACTION_ORDER_ERROR = 'order_error'; // 生成订单异常
    const ACTION_CHECKING = 'order_paying'; // 订单支付中间状态

    private static function transtion()
    {
        $transtion = [
            HaggleInitiate::STATUS_ING => [
                self::ACTION_LAST_TIMES => HaggleInitiate::STATUS_SUCCESS,
                self::ACTION_DOWN       => HaggleInitiate::STATUS_FAIL,
                self::ACTION_OVER       => HaggleInitiate::STATUS_FAIL,
                self::ACTION_TICKET     => HaggleInitiate::STATUS_TICKET_ING,
            ],

            HaggleInitiate::STATUS_SUCCESS => [
                self::ACTION_TICKET => HaggleInitiate::STATUS_TICKET_ING,
                self::ACTION_OVER        => HaggleInitiate::STATUS_OVER_TIME,
            ],

            HaggleInitiate::STATUS_TICKET_ING => [
                self::ACTION_CREATE      => HaggleInitiate::STATUS_ORDER_NOT_PAY,
                self::ACTION_OVER        => HaggleInitiate::STATUS_OVER_TIME,
                self::ACTION_ORDER_ERROR => HaggleInitiate::STATUS_ING,
            ],

            HaggleInitiate::STATUS_ORDER_NOT_PAY => [
                self::ACTION_CHECKING => HaggleInitiate::STATUS_PAY_CHECKING,
                self::ACTION_CANCEL   => HaggleInitiate::STATUS_ORDER_CANCEL,
                self::ACTION_PAY   => HaggleInitiate::STATUS_ORDER_PAY,
            ],

            HaggleInitiate::STATUS_PAY_CHECKING => [
                self::ACTION_PAY => HaggleInitiate::STATUS_ORDER_PAY,
            ],
        ];

        return $transtion;
    }

    /**
     * @param $enterpriseId int 企业id
     * @param $id int 主键
     * @param $action string 动作
     * @param $to
     *
     * @return mixed
     */
    public static function run($enterpriseId, $initiateId, $action, $to, $extUpdate = [])
    {
        $transtions = self::transtion();

        $history = ActivityStatus::getHistory($enterpriseId, 'HaggleInitiate', $initiateId);

        $fsm = Machine::create($transtions, $history);

        $assertEqualsTo = $fsm->doAction($action);

        // 成功
        if ($to != $assertEqualsTo) {
            notice('group_initiate_status_abnormal');
        }

        $activityInitiateCache = new ActivityInitiateCache();

        try {
            DB::beginTransaction();
            // 保存历史
            ActivityStatus::setHistory($enterpriseId, 'HaggleInitiate', $initiateId, $fsm->getHistory());

            HaggleInitiate::query()->where('id', $initiateId)->update(['status' => $to]);

            $update = array_merge(['status' => $to], $extUpdate);
            $activityInitiateCache->updateInitiateInfoCache($enterpriseId, $initiateId, ActivitySku::TYPE_HAGGLE, $update);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            $activityInitiateCache->deleteInitiateInfoCache($enterpriseId, $initiateId, ActivitySku::TYPE_HAGGLE);
            Error::log('更新砍价发起状态异常', $e);
            notice('machine_error');
        }

        return $to;
    }
}
