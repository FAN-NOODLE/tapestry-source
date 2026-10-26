<?php

namespace App\Http\Controllers\V1\Order;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class DetailController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        $input = $this->request->all();

        $this->validate($this->request, [
            'order_id'        => 'required|integer',
            'refund_apply_id' => 'integer'
        ], [
            'order_id.required'       => 'order_id必填',
            'order_id.integer'        => 'order_id需为integer',
            'refund_apply_id.integer' => 'refund_apply_id需为integer',
        ]);

        // 调用Shop SDK获取订单详情
        $data = (new Order())->orderDetail([
            'enterprise_id'   => $this->enterpriseId,
            'order_id'        => $input['order_id'],
            'refund_apply_id' => $input['refund_apply_id'] ?? 0
        ]);

        // 过滤不需要的字段
        foreach ($data['order_remark'] as $key => $value) {
            unset(
                $data['order_remark'][$key]['id'],
                $data['order_remark'][$key]['enterprise_id'],
                $data['order_remark'][$key]['order_no'],
                $data['order_remark'][$key]['admin_id']
            );
        }
        foreach ($data['order_logistic'] as $key => $value) {
            unset(
                $data['order_logistic'][$key]['id'],
                $data['order_logistic'][$key]['enterprise_id'],
                $data['order_logistic'][$key]['order_no'],
                $data['order_logistic'][$key]['refund_apply_no'],
                $data['order_logistic'][$key]['return_order_no'],
                $data['order_logistic'][$key]['sku_ids']
            );
        }

        return success($data);
    }
}
