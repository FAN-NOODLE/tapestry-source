<?php
/**
 * 微信先享卡
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019/12/25
 * Time: 17:51
 */

namespace App\Handlers\WxPayApi;

use App\Models\ActivityInteract;

class WxCardHandler extends WxPayClient
{
    const HOST = 'https://api.mch.weixin.qq.com/v3';
    /**
     * 获取约惠卡订单
     * get
     * %s = out_order_no
     * 更新约惠卡订单
     * patch
     * %s = out_order_no
     */
    const WX_DISCOUNT_CARD_ORDERS = '/discount-card/orders/%s';

    /**
     * 获取约惠卡订单
     *
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getDiscountCardOrders($data = [])
    {
        $url      = sprintf(self::WX_DISCOUNT_CARD_ORDERS, $data['out_order_no']);
        $response = $this->wxHttpRequest($url, $data, 'GET', self::HOST);
        return $response;
    }

    /**
     * 更新约惠卡订单
     *
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function updateDiscountCardOrders($data = [])
    {
        $requestId = ActivityInteract::request('更新约惠卡订单', $data);
        $url       = sprintf(self::WX_DISCOUNT_CARD_ORDERS, $data['out_order_no']);
        $response  = $this->wxHttpRequest($url, $data, 'PATCH', self::HOST);
        ActivityInteract::response($requestId, $response);
        return $response;
    }

}
