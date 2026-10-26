<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/13
 * Time: 17:36
 */

namespace App\StateMachine\Refund\Actions;

use App\Exceptions\Error;
use App\Models\Order;
use App\Models\RefundApply;
use App\Models\ReturnOrder;
use App\Models\ReturnOrderTrace;
use App\StateMachine\Refund\RefundMachineConfig;
use App\Services\RefundApplyService;
use App\Traits\RefundMachineTools;
use Illuminate\Support\Facades\Redis;
use Yansongda\Pay\Log;

class CreateRefundApply
{
    use RefundMachineTools;
    const LOCK_TIME = 1;

    protected $refundApplyService;

    public function run($event, $data)
    {
        // 订单模型
        $orderObj = Order::query()->find($data['order_id']);

        // 订单不存在
        if (empty($orderObj)) {
            throw new Error(codeNum('order_not_exists'), codeMsg('order_not_exists'), 210);
        }

        // 订单信息
        $orderArr = $orderObj->toArray();

        // 防多次申请售后
        $redisKey = redisKey('PREVENT_MULTI_REFUND_APPLY', ['order_id' => $data['order_id']]);
        $this->preventRepeatApply($redisKey, $data['order_id']);

        $toStatus   = RefundMachineConfig::ALLOW_EVENTS[$event]['toStatus'];
        $fromStatus = RefundMachineConfig::STATUS_INIT;
        $update     = $data;
        unset($update['order_id']);

        try {
            // 售后表
            list($applyId, $refund_apply_no) = (new RefundApplyService())->insertRefundApply($orderArr, $data['refund_fee'], $toStatus, $data);

            if (!empty($update['return_goods_info'])) {
                $insertArr = [];
                foreach ($update['return_goods_info'] as $return_good) {
                    $return_good['refund_apply_id'] = $applyId;
                    $insertArr[]                    = $return_good;
                }
                (new RefundApplyService())->insertAll($insertArr);

                //创建退货单
                $return_order_no = ReturnOrder::insertReturnOrder($orderArr['enterprise_id'], $orderArr['order_no'], $refund_apply_no, $data);

                // 退货单追踪表
                $attributes = [
                    'event'             => 'create_return_order',
                    'event_desc'        => '创建退货单',
                    'from_apply_status' => 0,
                    'to_apply_status'   => 1,
                    'remark'            => '',
                ];
                ReturnOrderTrace::insertReturnOrderTrace($return_order_no,$attributes);
            }
            // 售后追踪表
            (new RefundApplyService())->insertApplyTrace($applyId, $fromStatus, $toStatus, $event, $update);
            //通知open
            $data['order_no'] = $orderArr['order_no'];
            switch ($event) {
                case RefundMachineConfig::EVENT_CUSTOMER_APPLY_REFUND_OF_NOT_SEND:
                    $data['refund_apply_id'] = $applyId;
                    if($data['third_order_id']!=0){
                        $this->only_refund_money($data);
                    }
                    break;
                default:
                    break;
            }
        } catch (\Exception $e) {
            Log::error(codeMsg('refund_apply_event_fail', [$event, $e->getMessage()]));

            throw new Error(codeNum('refund_apply_event_fail'), codeMsg('refund_apply_event_fail', [$event, $e->getMessage()]), 210);
        }
    }

    /**
     * 防止多次发售后
     */
    protected function preventRepeatApply($redisKey, $orderId)
    {
        // 存在，不允许申请售后
        if (getOrSetData($redisKey)) {
            throw new Error(codeNum('order_apply_refund_repeat'), codeMsg('order_apply_refund_repeat'), 210);
        }
        // 加锁，不允许其他用户操作
        getOrSetData($redisKey, 1, self::LOCK_TIME);

        $refundObj = RefundApply::query()->where('order_id', $orderId)->orderBy('id', 'desc')->first();
        if (!empty($refundObj) && !in_array($refundObj->status, RefundMachineConfig::FINAL_STATES)) {
            throw new Error(codeNum('order_apply_refund_repeat'), codeMsg('order_apply_refund_repeat'), 210);
        }
    }

}