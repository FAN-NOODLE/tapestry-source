<?php

namespace App\Http\Controllers\V1\NewOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class PushOrderToErpController extends BaseController
{
    /**
     * 手动-补推ERP查不到的订单
     * {domain}/gw-console/v1/new-order/push-order-to-erp
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index()
    {
        $orderNo = $this->request->post('order_no', '');
        if (empty($orderNo)) {
            notice(320264);
        }

        try {
            (new Order())->pushOrderToErp(['order_no' => strval($orderNo)]);
        } catch (\Exception $e) {
        }

        return success();
    }
}
