<?php

namespace App\Http\Controllers\V1\Order;

use CrsSdk\Shop\Order;
use App\Http\Controllers\V1\BaseController;

class CloseController extends BaseController
{

    /**
     * 关闭订单并且退款
     *
     * @return array
     */
    public function index()
    {
        // 数据验证
        $this->validate($this->request, [
            'order_no' => 'required'
        ], [
            'order_no.required' => '订单号必填'
        ]);

        // 获取请求数据
        $input = $this->request->input();

        $order = new Order();

        // 先将订单状态修改为申请退款中
        $data = $order->manualRefundStatus([
            'order_no' => $input['order_no']
        ]);
        if ($data['data']['status']) {
            // 调Shop直接进行退款
            $data = $order->refundOrder([
                'order_no' => $input['order_no'],
                'status'   => true
            ]);
        }

        return success(['data' => $data]);
    }

}
