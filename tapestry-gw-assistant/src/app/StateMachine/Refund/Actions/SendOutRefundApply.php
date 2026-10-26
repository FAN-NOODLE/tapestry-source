<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/13
 * Time: 17:36
 */

namespace App\StateMachine\Refund\Actions;

use App\Exceptions\Error;
use App\Filters\Common\DeleteOtherData;
use App\Filters\Order\Validate\ValidateOrderSplit;
use App\Filters\RefundApply\Machine\BuildRefundApplyGoods;
use App\Filters\RefundApply\Machine\BuildRefundApplyInfo;
use App\Filters\RefundApply\Machine\CheckEnterpriseSplitConfig;
use App\Filters\RefundApply\Machine\CheckRefundApplyGoods;
use App\Models\RefundApply;
use App\Models\RefundApplyGoods;
use App\Models\RefundApplyTrace;
use App\Models\ReturnOrder;
use App\Models\ReturnOrderTrace;
use App\StateMachine\Refund\RefundMachineConfig;
use App\Traits\RefundMachineTools;
use Illuminate\Pipeline\Pipeline;

class SendOutRefundApply
{
    use RefundMachineTools;
    const LOCK_TIME = 1;

    protected $refundApplyService;

    public function run($event, $data)
    {
        // 防多次申请售后
        $redisKey = redisKey('PREVENT_MULTI_REFUND_APPLY', ['order_id' => $data['order_id']]);
        $this->preventRepeatApply($redisKey, $data['order_id']);

        $toStatus   = RefundMachineConfig::ALLOW_EVENTS[$event]['toStatus'];
        $fromStatus = RefundMachineConfig::STATUS_INIT;

        //业务检查和组织数据
        $pipes = [
            CheckRefundApplyGoods::class, //检查是否整单全退，售后申请的商品是否已退完
            CheckEnterpriseSplitConfig::class, //检查企业是否支持拆单
            ValidateOrderSplit::class, //检查订单是否支持拆单
            BuildRefundApplyInfo::class, //构建售后单信息
            BuildRefundApplyGoods::class, //构建售后商品信息
            DeleteOtherData::class, //删除other信息
        ];

        $data = app(Pipeline::class)
            ->send($data)
            ->through($pipes)
            ->then(function ($content) {
                return $content;
            });

        //插入售后单
        list($applyId, $refund_apply_no) = RefundApply::insertRefundApply($data);

        if (!empty($data)) {
            foreach ($data['return_goods_info'] as &$return_good) {
                $return_good['refund_apply_id'] = $applyId;
            }
        }

        if ($event == RefundMachineConfig::EVENT_SEND_OUT_APPLY_ONLY_REFUND) {
            if ($data['isAllRefund'] == RefundApply::ALL_REFUND_FALSE && !empty($data['return_goods_info'])) {
                RefundApplyGoods::insertRefundApplyGoods($data);
                //创建退货单
                $return_order_no = ReturnOrder::insertReturnOrder($data['refund_apply_info']['enterprise_id'], $data['refund_apply_info']['order_no'], $refund_apply_no, $data);

                // 退货单追踪表
                $attributes = [
                    'event'             => 'create_return_order',
                    'event_desc'        => '创建退货单',
                    'from_apply_status' => 0,
                    'to_apply_status'   => 1,
                    'remark'            => '',
                ];
                ReturnOrderTrace::insertReturnOrderTrace($return_order_no, $attributes);
            }
        } else {
            if (!empty($data['return_goods_info'])) {
                if($event == RefundMachineConfig::EVENT_MANUAL_APPLY_RETURN_GOODS && $data['refund_apply_info']['refund_apply_type']==1 && $data['isAllRefund'] == RefundApply::ALL_REFUND_TRUE){

                }else{
                    RefundApplyGoods::insertRefundApplyGoods($data);
                    //创建退货单
                    $return_order_no = ReturnOrder::insertReturnOrder($data['refund_apply_info']['enterprise_id'], $data['refund_apply_info']['order_no'], $refund_apply_no, $data);

                    // 退货单追踪表
                    $attributes = [
                        'event'             => 'create_return_order',
                        'event_desc'        => '创建退货单',
                        'from_apply_status' => 0,
                        'to_apply_status'   => 1,
                        'remark'            => '',
                    ];
                    ReturnOrderTrace::insertReturnOrderTrace($return_order_no, $attributes);
                }

            }
        }

        // 售后追踪表
        RefundApplyTrace::insertRefundApplyTrace($applyId, $fromStatus, $toStatus, $event, $data);
    }

    /**
     * 防止多次发售后
     *
     * @param $redisKey
     * @param $orderId
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