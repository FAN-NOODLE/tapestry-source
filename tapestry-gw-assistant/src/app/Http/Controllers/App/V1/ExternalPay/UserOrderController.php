<?php

namespace App\Http\Controllers\App\V1\ExternalPay;

use App\Http\Controllers\Controller;
use App\Micros\Work\ExternalPay\ExternalPay;
use Illuminate\Http\JsonResponse;

class UserOrderController extends Controller
{
    /**
     * @var $inherit
     */
    private $order;

    public function __construct()
    {
        parent::__construct();
        $this->order = new ExternalPay();
    }

    /**
     * 订单列表
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function orderList(): JsonResponse
    {
        $tokenData               = request()->get('tokenData');
        $params                  = request()->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $result = $this->order->userOrderList($params);
        return success($result);
    }

    /**
     * 退款列表
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function refundOrderList(): JsonResponse
    {
        $tokenData               = request()->get('tokenData');
        $params                  = request()->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $result = $this->order->userRefundOrderList($params);
        return success($result);
    }

    /**
     * @return void
     * 物流信息
     */
    public function logisticsTrack(): JsonResponse
    {
        $tokenData               = request()->get('tokenData');
        $params                  = request()->all();
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $result = (new ExternalPay())->logisticsTrack($params);

        return success($result);
    }


    /**
     * 订单详情
     *
     * @return void
     */
    public function orderDetail(): JsonResponse
    {
        $tokenData               = request()->get('tokenData');
        $params                  = request()->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];
        $result                  = $this->order->userOrderDetail($params);
        return success($result);
    }

    /**
     * 退款详情
     *
     * @return void
     */
    public function refundOrderDetail(): JsonResponse
    {
        $tokenData               = request()->get('tokenData');
        $params                  = request()->all();
        $params['union_id']      = $tokenData['union_id'];
        $params['enterprise_id'] = $tokenData['enterprise_id'];

        $result = $this->order->userRefundOrderDetail($params);
        return success($result);
    }
}