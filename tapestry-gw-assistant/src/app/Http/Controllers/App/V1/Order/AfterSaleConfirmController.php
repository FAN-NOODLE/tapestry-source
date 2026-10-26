<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/11 0011
 * Time: 14:41
 */

namespace App\Http\Controllers\App\V1\Order;


use App\Common\Constant;
use CrsSdk\Shop\NewOrder;
use CrsSdk\Shop\Refund;
use Illuminate\Support\Facades\Log;

class AfterSaleConfirmController extends OrderBaseController
{
    /**
     * 售后审核状态流转
     *
     * @var array[]
     */
    private $refund
        = [
            1 => [
                'desc'   => '同意售后',
                'field'  => [
                    'id', 'store_id', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::REFUND_CONSOLE_AGREE,
                'object' => 'refund',
                'other'  => []
            ],
            2 => [
                'desc'   => '待售后驳回申请不符', // 状态25使用
                'field'  => [
                    'id', 'store_id', 'reject_reason', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::MANUAL_REFUSE_APPLY_REFUND,
                'object' => 'newOrder',
                'other'  => []
            ],
            3 => [
                'desc'   => '订单状态变更待收货',
                'field'  => [
                    'id', 'store_id', 'type', 'logistic_identity', 'logistic_no', 'reject_reason', 'refund_apply_id'
                ],
                'event'  => Constant::MANUAL_SEND_OUT_ORDER,
                'object' => 'newOrder',
                'other'  => []
            ],
            4 => [
                'desc'   => '订单状态变更已完成已退货',
                'field'  => [
                    'id', 'store_id', 'reject_reason', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::ORDER_MANUAL_COMPLETE,
                'object' => 'newOrder',
                'other'  => [
                    'complete_type' => 'is_refund',
                    'reject_reason' => '已退款'
                ]
            ],
            5 => [
                'desc'   => '订单状态变更已完成换货',
                'field'  => [
                    'id', 'store_id', 'reject_reason', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::ORDER_MANUAL_COMPLETE,
                'object' => 'newOrder',
                'other'  => [
                    'complete_type' => 'change_goods',
                    'reject_reason' => '换货处理'
                ]
            ],
            6 => [
                'desc'   => '订单状态变更已完成其他',
                'field'  => [
                    'id', 'store_id', 'reject_reason', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::ORDER_MANUAL_COMPLETE,
                'object' => 'newOrder',
                'other'  => [
                    'complete_type' => 'other',
                    'reject_reason' => '其他'
                ]
            ],
            7 => [
                'desc'   => '同意退款',
                'field'  => [
                    'id', 'store_id', 'type', 'refund_apply_id'
                ],
                'event'  => Constant::REFUND_MANUAL_AGREE_REFUND,
                'object' => 'newOrder',
                'other'  => [
                    'reject_reason' => '退款'
                ]
            ],
            8 => [
                'desc'   => '待退款的驳回申请不符',
                'field'  => [
                    'id', 'store_id', 'type', 'refund_apply_id', 'reject_reason'
                ],
                'event'  => 'manual_refuse_apply_refund',
                'object' => 'newOrder',
                'other'  => []
            ],
            9 => [
                'desc'   => '同意退款退货',
                'field'  => [
                    'id', 'store_id', 'type', 'refund_apply_id'
                ],
                'event'  => 'open_agree_return',
                'object' => 'refund',
                'other'  => []
            ],
            10 => [
                'desc'   => '待售后驳回申请不符', // 状态21使用
                'field'  => [
                    'id', 'store_id', 'reject_reason', 'type', 'refund_apply_id','logistic_no','logistic_identity','complete_type'
                ],
                'event'  => 'console_refuse_apply_refund',
                'object' => 'newOrder',
                'other'  => []
            ]
        ];


    /**
     * 售后单处理
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules = [
            'type'     => 'required|in:' . implode(',', array_keys($this->refund)),
            'id'       => 'required',
        ];

        $messages    = [
            'type.required'     => '售后操作类型不能为空',
            'type.in'           => '售后操作类型不合法',
            'id.required'       => '售后单ID不能为空'
        ];
        $input       = $this->validate($this->request, $rules, $messages);
        $params      = $this->refund[$input['type']];
        $requestData = $this->request->only($params['field']);
        try {
            $this->checkRequestParams($params, $requestData);
            $requestData          = array_merge($requestData, $params['other']);
            $requestData['event'] = $params['event'];

            // 增加操作的其他信息，管理员信息、门店信息
            $requestData['enterprise_id']    = $this->enterpriseId;
            $requestData['operate_store_id'] = $input['store_id']??0;
            $requestData['operate_id']       = $this->adminId;
            $requestData['operate_name']     = $this->mobile;
            unset($requestData['type']);
            switch ($params['object']) {
                case 'refund':
                    $requestData['id'] = $requestData['refund_apply_id'];
                    unset($requestData['refund_apply_id']);
                    (new Refund())->machineRun($requestData);
                    break;
                case 'newOrder':
                    unset($requestData['refund_apply_id']);
                    (new NewOrder())->changeOrderState($requestData);
                    break;
            }
            Log::info("【门店助手】售后单[{$params['desc']}]处理成功", $requestData);
        } catch (\Exception $e) {
            Log::error("【门店助手】售后单[{$params['desc']}]处理失败, 请求参数：", $requestData);
            // erp 比人工快的情况下单独处理
            if ($e->getCode() == 6006404) {
                notice('order_status_change');
            }else{
                error($e->getCode(), $e->getMessage());
            }
        }

        return success();
    }

    /**
     * 判断字段值
     *
     * @param $defaultParam
     * @param $params
     * @throws \Exception
     */
    private function checkRequestParams($defaultParam, $params)
    {
        foreach ($defaultParam['field'] as $field) {
            if (!isset($params[$field])) {
                throw new \Exception($defaultParam['desc'] . "中{$field}不能为空", 654014);
            }
        }
    }
}
