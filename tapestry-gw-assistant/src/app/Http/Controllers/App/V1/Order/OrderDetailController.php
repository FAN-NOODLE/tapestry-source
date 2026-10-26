<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/4/8 0008
 * Time: 18:19
 */

namespace App\Http\Controllers\App\V1\Order;


use App\Common\Constant;
use CrsSdk\Shop\Order;
use CrsSdk\Open\Logistics;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Log;

class OrderDetailController extends OrderBaseController
{
    /**
     * 获取订单详情
     *
     * @param Order       $order
     * @param Logistics   $logistics
     *
     * @return \Illuminate\Http\JsonResponse|void
     * @throws \Exception
     */
    public function index(Order $order, Logistics $logistics)
    {
        $orderNo = $this->request->get('order_no');
        $storeId = $this->request->get('store_id');
        if (empty($orderNo)) {
            error(654004);
        }
        $enterpriseId            = Auth::user()['enterprise_id'];
        $params['enterprise_id'] = $enterpriseId;
        $params['order_no']      = $orderNo;

        //获取店员id
        $workerDetail = $this->storeWorkerDetail($enterpriseId, $storeId);

        try {
            $data = $order->orderDetail($params);
            // 物流信息
            $orderLogistic    = $data['order_logistic'];
            $logisticTraceEnd = $logisticTraceEnds = [];
            if (!empty($orderLogistic)) {
                unset($data['order_logistic']);
                // 查询物流信息
                $logisticTrace = $logistics->query([
                    'data' => [
                        [
                            'logistic_no'  => $orderLogistic[0]['logistic_no'],
                            'logistic_key' => $orderLogistic[0]['logistic_identity']
                        ]
                    ]
                ]);
                if (isset($logisticTrace['data']['success_list']) &&
                    !empty($logisticTrace['data']['success_list'])) {
                    $logisticTraceEnd = $logisticTrace['data']['success_list'][0];
                }else{
                    $logisticTrace['data']['fail_list'][0]['traces'] = [];
                    $logisticTraceEnd                                = $logisticTrace['data']['fail_list'][0];
                }

                /***----------------------多物流显示部分开始-------------------------***/
                // 多快递查询
                foreach ($orderLogistic as $value) {
                    // 查询物流信息
                    $logisticTrace = $logistics->query([
                        'data' => [
                            [
                                'logistic_no'  => $value['logistic_no'],
                                'logistic_key' => $value['logistic_identity']
                            ]
                        ]
                    ]);
                    if (isset($logisticTrace['data']['success_list']) &&
                        !empty($logisticTrace['data']['success_list'])) {
                        $logisticTraceEnds[] = $logisticTrace['data']['success_list'][0];
                    }else{
                        $logisticTrace['data']['fail_list'][0]['traces'] = [];
                        $logisticTraceEnds[]                             = $logisticTrace['data']['fail_list'][0];
                    }

                }
                /***----------------------多物流显示部分结束-------------------------***/
            }
            $data['order_logistic'] = $logisticTraceEnd;
            // 多物流列表
            $data['order_logistics'] = $logisticTraceEnds;

            //  重新结算订单总额和优惠金额
            $data['coupon_discount_amount']               = $data['discount_fee'];
            $data['main_order']['coupon_discount_amount'] = $data['discount_fee'];
            $data['order_total_amount']                   = $data['cash_fee'] + $data['preferential_price'] - $data['total_refund_fee'];
            $data['main_order']['order_total_amount']     = $data['cash_fee'] + $data['preferential_price'] - $data['total_refund_fee'];

            $serviceOrderType = 0;
            // 没人有推广人则推广订单
            if ($data['store_info']['service_staff_id'] == 0) {
                $serviceOrderType = Constant::RECEIVED_ORDER_TYPE;
            }
            // 如果当前用户和服务人是同一个人
            if (!empty($workerDetail['id']) && $workerDetail['id'] ?? '' == $data['store_info']['service_staff_id']) {
                $serviceOrderType = Constant::MY_ORDER_TYPE;
            }

            // 如果是店长或者企业管理角色
            if (
                in_array($this->role, [
                    Constant::LOGIN_ROLE_STORE_MANAGER,
                    Constant::LOGIN_ROLE_ENTERPRISE_ADMIN
                ])
                && !empty($data['store_info']['service_staff_id'])) {
                $serviceOrderType = Constant::MY_ORDER_TYPE;;
            }

            $data['service_order_type'] = $serviceOrderType;

            $data['parcel_button'] = false;
            // 需配送显示添加包裹按钮(类型为快递,状态已付款,有服务人)
            if ($data['delivery_type'] == 1 && $data['status'] == 60 && $data['service_staff_id'] > 0) {
                $data['parcel_button'] = true;
            }

            $orderRemark = '';
            foreach ($data['order_remark'] ?? [] as $item) {
                if ($item['type'] == 2) {
                    $orderRemark = $item;
                }
            }
            $data['order_remark'] = $orderRemark;
        } catch (\Exception $exception) {
            Log::error("获取订单详情失败，错误信息：{$exception->getMessage()}");
            error($exception->getCode(), $exception->getMessage());
        }

        $data['is_refund_authority'] = (int)config('app.is_refund_authority', true);

        return success($data);
    }
}
