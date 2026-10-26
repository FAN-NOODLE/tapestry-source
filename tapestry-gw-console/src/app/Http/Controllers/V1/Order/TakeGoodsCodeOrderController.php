<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class TakeGoodsCodeOrderController extends BaseController
{
    /**
     * @return JsonResponse|void
     * @throws ValidationException
     */
    public function index()
    {
        $input = $this->request->all();

        $this->validate($this->request, [
            'take_goods_code' => 'required',
        ], [
            'take_goods_code.required' => '自提码必填',
        ]);

        //根据自提码查询订单，返回多个提示错误，返回一个，再查询订单详情
        $order_nos = (new Order())->searchByTakeCode([
            'enterprise_id'     => $this->enterpriseId,
            'take_goods_code'   => $input['take_goods_code'],
            'take_goods_name'   => $input['take_goods_name'] ?? '',
            'take_goods_mobile' => $input['take_goods_mobile'] ?? '',
            'order_no'          => $input['order_no'] ?? '',
        ])['order_nos'];

        if (empty($order_nos)) {
            notice(100000, ['未找到您搜索的自提订单']);
        }
        if (count($order_nos) != 1) {
            notice(100000, ['无法识别唯一订单，请补充订单其它信息帮助定位订单']);
        }

        // 调用Shop SDK获取订单详情
        $data = (new Order())->orderDetail([
            'enterprise_id'   => $this->enterpriseId,
            'order_no'        => $order_nos[0],
            'refund_apply_id' => 0,
            'api_type'        => 1 // 小程序版本：0 老版，1 新版
        ]);

        return success($data);
    }
}
