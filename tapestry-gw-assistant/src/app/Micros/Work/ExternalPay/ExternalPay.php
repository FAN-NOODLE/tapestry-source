<?php

namespace App\Micros\Work\ExternalPay;

use CrsSdk\Sdk;

class ExternalPay
{
    /**
     * 物流公司
     * @param $data
     * @return array|mixed
     * @throws \Exception
     */
    public function logisticsCompany($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/logistics-company")->post($data);
    }

    /**
     * 导购端订单列表
     * @param $data
     * @return array|mixed
     * @throws \Exception
     */
    public function orderList($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/order-list")->post($data);
    }

    /**
     * 导购退款订单列表
     * @param $data
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function refundOrderList($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/refund-order-list")->post($data);
    }

    /**
     * 订单发货
     * @param $data
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function deliverOrder($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/deliver-order")->post($data);
    }

    /**
     * 订单详情
     * @return array|mixed
     * @throws \Exception
     */
    public function orderDetail($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/order-detail")->post($data);
    }

    /**
     * 退款详情
     * @return void
     */
    public function refundOrderDetail($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/refund-order-detail")->post($data);
    }

    /**
     * 物流信息
     * @param $data
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function logisticsTrack($data):array
    {
        return Sdk::work()->setApi("service/v1/external-pay/order-logistics_track")->post($data);
    }

    /**
     * 用户订单列表
     * @return void
     */
    public function userOrderList($data):array
    {
        return Sdk::work()->setApi("service/v1/user-external-pay/order-list")->post($data);
    }

    /**
     * 用户退款列表
     * @param $data
     *
     * @return array
     * @throws \Exception
     */
    public function userRefundOrderList($data):array
    {
        return Sdk::work()->setApi("service/v1/user-external-pay/refund-order-list")->post($data);
    }


    /**
     * 用户订单详情
     * @return array|mixed
     * @throws \Exception
     */
    public function userOrderDetail($data):array
    {
        return Sdk::work()->setApi("service/v1/user-external-pay/order-detail")->post($data);
    }

    /**
     * 用户退款详情
     * @return void
     */
    public function userRefundOrderDetail($data):array
    {
        return Sdk::work()->setApi("service/v1/user-external-pay/refund-order-detail")->post($data);
    }
}