<?php

namespace App\Http\Controllers\App\V1\Order;

use CrsSdk\Shop\Order;
use Illuminate\Support\Facades\Log;

class AfterSaleController extends OrderBaseController
{
    /**
     * 售后详情
     *
     * @param Order $order
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function index(Order $order)
    {
        $id = $this->request->get('order_no');
        if (empty($id)) {
            error(654004);
        }

        $refundApplyId = $this->request->get('refund_apply_id');
        if (empty($refundApplyId)) {
            error(654013);
        }

        $params['enterprise_id']   = $this->user['enterprise_id'];
        $params['order_no']        = $id;
        $params['refund_apply_id'] = $refundApplyId;

        try {
            $data = $order->orderDetail($params);
        } catch (\Exception $e) {
            Log::error("【门店助手】售后订单详情查询失败,错误信息: {$e->getMessage()}, 请求参数：", $params);
            error($e->getCode(), $e->getMessage());
        }

        $data['is_refund_authority'] = (int)config('app.is_refund_authority', true);
        return success($data);
    }
}
