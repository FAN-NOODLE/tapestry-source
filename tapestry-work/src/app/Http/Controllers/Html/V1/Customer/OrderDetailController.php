<?php

namespace App\Http\Controllers\Html\V1\Customer;

use App\Http\Controllers\Html\V1\BaseController;
use App\Services\OrderService;
use CrsSdk\Open\Logistics;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class OrderDetailController extends BaseController
{
    /**
     * 订单详情
     *
     * @param Request $request
     * @param Order $order
     * @return array|void
     * @throws \Exception
     */
    public function index(Request $request, Order $order)
    {
        // 订单id
        $orderNo = $request->get('order_no', 0);
        if (empty($orderNo)) {
            notice(200017);
        }

        $orderDetail = OrderService::getOrderDetail($orderNo);
        return success(['data' => $orderDetail]);

        /*
        switchAppHostByDomain($this->domain);
        $_SERVER['enterprise_hash'] = $this->enterprise['hash'];
        $result = $order->orderDetail(['order_no' => $orderNo]);
        $result['pay_data'] = $result['main_order']['pay_data'];
        unset($result['main_order']['pay_data']);
        return success(['data' => $result]);
        */
    }
}
