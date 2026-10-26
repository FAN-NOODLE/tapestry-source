<?php

namespace App\Micros\Work\Activity;

use CrsSdk\Sdk;
use Exception;

class ReserveCIIE
{
    /**
     * 活动信息
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function Info(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/info")->post($data);
    }

    /**
     * 创建预约
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function create(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/create")->post($data);
    }

    /**
     * 取消预约
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function cancel(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/cancel")->post($data);
    }

    /**
     * 编辑预约
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function edit(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/edit")->post($data);
    }

    /**
     * 用户进入活动记录上报
     *
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function UserView(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/user-view")->post($data);
    }

    /**
     * 预约列表
     *
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function ReserveLists(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/lists")->post($data);
    }


    /**
     * 核销券
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function writeOffCoupon(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/write-off-coupon")->post($data);
    }

    /**
     * 白名单
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function WhiteLists(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/white-list")->post($data);
    }

    /**
     * 导购的核销列表
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function CheckLists(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/check-lists")->post($data);
    }


    /**
     * 预约列表
     *
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function TimeLists(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/time-lists")->post($data);
    }

    /**
     * 配置活动预约时间段
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function TimeConfig(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/time-config")->post($data);
    }

    /**
     * 配置活动预约时间段
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function GetPhoneNumber(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/reserve/get-phone-number")->post($data);
    }


}