<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class RefundInfoDataController extends BaseController
{
    /**
     * 获取售后信息弹窗数据
     * {domain}/app/v1/order/refund-info-data
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->only(['order_no']);

        $this->validate($this->request, [
            'order_no' => 'required'
        ], [
            'order_no.required' => '[参数]订单号order_no必须'
        ]);

        $response = (new Order())->refundInfoData($input);

        return success($response);
    }
}
