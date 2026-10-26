<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019/6/18
 * Time: 17:37
 */

namespace App\StateMachine\Refund\Actions;

use App\Exceptions\Error;
use App\Models\OrderLogistic;
use App\Models\RefundApply;
use App\StateMachine\Refund\RefundMachineConfig;
use App\Services\RefundApplyService;
use App\StateMachine\Action;
use App\Traits\RefundMachineTools;
use Illuminate\Support\Facades\Redis;
use Yansongda\Pay\Log;

class MoveRefundApply implements Action
{
    use RefundMachineTools;

    const LOCK_TIME = 1;

    protected $refundApplyService;

    public function handle($from, $event, $to, $data, $applyObj, $eventDesc)
    {

    }

    public function run($event, $post)
    {
        $refundApplyService = new RefundApplyService();

        // 加锁(防并发)
        $redisKey = redisKey('PREVENT_OPERATE_SAME_REFUND_APPLY', ['id' => $post['id']]);
        $this->preventConcurrent($redisKey);

        $applyObj = RefundApply::query()->find($post['id']);

        // 售后单不存在
        if (empty($applyObj)) {
            throw new Error(codeNum('refund_apply_not_exists'), codeMsg('refund_apply_not_exists'), 210);
        }

        $applyArr      = $applyObj->toArray();
        $currentStatus = $applyArr['status'];

        // 状态验证
        $this->validateState($event, $currentStatus);

        $eventDesc  = RefundMachineConfig::ALLOW_EVENTS[$event]['desc'];
        $toStatus   = RefundMachineConfig::ALLOW_EVENTS[$event]['toStatus'];
        $fromStatus = $currentStatus;
        $update     = $post;
        unset($post['id']);
        $update['status'] = $toStatus;

        try {
            // 更新售后申请状态
            $applyObj->update($update);

            // 售后追踪表
            $refundApplyService->insertApplyTrace($applyArr['id'], $fromStatus, $toStatus, $event, $update);

            //通知open
            $data['order_no'] = $applyArr['order_no'];
            switch ($event) {
                case RefundMachineConfig::EVENT_CONSOLE_AGREE:
                    //后台同意退款或者退款退货类型的售后
                    if ($applyArr['refund_apply_type'] == RefundApply::TYPE_REFUND_ONLY_MONEY) {
                        $this->only_refund_money($data);
                    } else {
                        $this->console_agree_refund_money_goods($data);
                    }
                    break;
                case RefundMachineConfig::EVENT_CUSTOMER_CANCEL:
                    //用户取消仅退款或退货退款类型的售后
                    if ($applyArr['refund_apply_type'] == RefundApply::TYPE_REFUND_ONLY_MONEY) {
                        $this->cancel_only_refund_money($data);
                    } else {
                        $this->cancel_refund_money_goods($data);
                    }
                    break;
                case RefundMachineConfig::EVENT_CUSTOMER_DELIVER:
                    // 用户寄回商品
                    OrderLogistic::setData($applyObj->enterprise_id, $update['logistic_identity'], $update['logistics_no'], [
                        'order_no'        => $applyObj->order_no,
                        'refund_apply_no' => $applyObj->out_refund_no,
                        'type'            => OrderLogistic::TYPE_IN
                    ]);

                    $this->customer_deliver($data);
                    break;
                default:
                    break;
            }

        } catch (\Exception $e) {
            Log::notice(codeMsg('refund_apply_event_fail', [$eventDesc, $e->getMessage()]));

            throw new Error(codeNum('refund_apply_event_fail'), codeMsg('refund_apply_event_fail', [$eventDesc, $e->getMessage()]), 210);
        }

    }

    /**
     * 状态验证
     */
    protected function validateState($event, $currentStatus)
    {
        // 被允许操作的状态
        $allowStatus = RefundMachineConfig::ALLOW_EVENTS[$event]['allowStatus'];
        $eventDesc   = RefundMachineConfig::ALLOW_EVENTS[$event]['desc'];

        if (!in_array($currentStatus, (array)$allowStatus)) {
            throw new Error(codeNum('refund_apply_status_abnormal'), codeMsg('refund_apply_status_abnormal', [$eventDesc]), 210);
        }
    }

    /**
     * 防并发
     */
    protected function preventConcurrent($redisKey)
    {
        // 存在，不允许当前用户操作
        if (getOrSetData($redisKey)) {
            throw new Error(codeNum('refund_apply_operate_by_other'), codeMsg('refund_apply_operate_by_other'), 210);
        }

        // 加锁，不允许其他用户操作
        getOrSetData($redisKey, 1, self::LOCK_TIME);
    }
}