<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class LogisticsDeleteController extends BaseController
{
    /**
     * 删除订单物流
     * {domain}/app/v1/order/logistics-delete
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only(['order_no', 'logistic_no']);
        
        $this->validate($this->request, [
            'order_no'      => 'required',
            'logistic_no'   => 'required'
        ], [
            'order_no.required'      => '[参数]订单号order_no必填',
            'logistic_no.required'   => '[参数]物流单号logistic_no必填'
        ]);

        $input['enterprise_id'] = $this->enterpriseId;
        $response = (new Order())->logisticsDelete($input);

        return success($response);
    }
}
