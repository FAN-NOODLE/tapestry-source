<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;

class AfterSaleAddController extends BaseController
{
    /**
     * 创建手工售后单
     * {domain}/app/v1/order/after-sale-add
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        // 参数验证
        $this->validate($this->request, [
            'order_no'           => 'required',
            'refund_apply_type'  => 'required | in:1,2',
            'order_goods_status' => 'required',
            'user_problem'       => 'required | max:100',
            'return_goods_info'  => 'required | array'
        ], [
            'order_no.required'           => '[参数]订单号order_no必填',
            'refund_apply_type.required'  => '[参数]售后方式refund_apply_type必填',
            'refund_apply_type.in'        => '[参数]售后方式refund_apply_type值不合法',
            'order_goods_status.required' => '[参数]货物状态order_goods_status必填',
            'user_problem.required'       => '请填写售后原因',
            'user_problem.max'            => '[参数]售后原因最大字符数100',
            'return_goods_info.required'  => '[参数]return_goods_info必填',
            'return_goods_info.array'     => '[参数]return_goods_info格式不正确'
        ]);

        $input = $this->request->only([
            'order_no',
            'refund_apply_type',
            'order_goods_status',
            'user_problem',
            'return_goods_info'
        ]);

        // 前端传的参数格式不对，前端说无法处理，所以做此转化
        $returnGoodsInfo = $input['return_goods_info'];
        array_walk($returnGoodsInfo, function (&$value) {
            $value = json_decode($value, true);
        });
        $input['return_goods_info'] = $returnGoodsInfo;

        $response = (new Order())->afterSaleAdd($input);

        return success($response);
    }
}
