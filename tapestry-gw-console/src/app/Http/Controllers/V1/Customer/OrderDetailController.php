<?php

namespace App\Http\Controllers\V1\Customer;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use Exception;

class OrderDetailController extends BaseController
{
    /**
     * 客户订单详情
     *
     * @param Request $request
     * @return JsonResponse|void
     * @throws Exception
     */
    public function index(Request $request)
    {
        // 订单id
        $orderId = $request->get('order_id', 0);
        if (!$orderId) {
            error(320251);
        }

        // 调用sdk
        $detail = (new Order())->orderDetail(['order_id' => $orderId]);

        // 初始化待返数据
        $data = [];
        // 数据格式化
        if ($detail) {
            // 订单商品
            $goodsData = $this->formatGoods($detail['goods_data']);
            // 交易流水号
            $thirdPaymentNo = $detail['pay_data']['third_payment_no'] ?? '';

            // 组装待返数据
            $data = [
                'order_id'           => $detail['id'],
                'order_no'           => $detail['order_no'],
                'created_at'         => intval($detail['created_at']),
                'pay_at'             => $detail['pay_at'],
                'status'             => $detail['status'],
                'payment_channel'    => $detail['payment_channel'],
                'third_payment_no'   => $thirdPaymentNo,
                'total_amount'       => $detail['total_amount'],
                'preferential_price' => $detail['preferential_price'],
                'final_amount'       => $detail['final_amount'],
                'goods_data'         => $goodsData
            ];
        }

        return success($data);
    }

    /**
     * 订单商品数据格式化
     *
     * @param array $goods 待格式化订单商品列表
     *
     * @return array
     */
    public function formatGoods($goods)
    {
        // 初始化待返数据
        $data = [];
        if (!$goods) {
            return $data;
        }

        // 组装待返数据
        foreach ($goods as $good) {
            // 待返数据赋值
            $data[] = [
                'spu_id'        => $good['spu_id'],
                'spu_name'      => $good['spu_name'],
                'spu_picture'   => $good['spu_picture'],
                'sell_price'    => $good['sell_price'],
                'quantity'      => $good['quantity'],
                'sub_price'     => $good['sub_price'],
                'category_name' => implode('/', $good['category_name'])
            ];
        }

        return $data;
    }
}
