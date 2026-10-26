<?php
/**
 * 售后单用户填入物流信息后审核订单
 * Created by PhpStorm.
 * User: cjp
 * Date: 2019/6/3
 * Time: 18:26
 */
namespace App\Http\Controllers\V1\Refund;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;

class HandRefundController extends BaseController
{
    public function index(Request $request)
    {
        $input = $request->only(['order_no', 'status']);

        $rules = [
            'order_no' => 'required',
            'status' => 'required|in:0,1',
        ];

        $message = [
            'order_no.required' => '请填写订单号',
            'status.required' => '请填写审核状态',
            'status.in' => '审核状态类型不正确'
        ];

        $this->validate($request, $rules, $message);

        $this->writeOperatorLog("后台审核售后单是否同意退款");

        return (new Order())->refundOrder($input);
    }
}