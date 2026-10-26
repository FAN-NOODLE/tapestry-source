<?php

namespace App\Micros\Work\Activity;

use CrsSdk\Sdk;
use Exception;

class Scratch
{
    /**
     * 券列表
     *
     * @param array $data
     *
     * @return array|mixed
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:44
     */
    public function Info(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/info")->post($data);
    }

    /**
     * 券列表
     *
     * @param array $data
     *
     * @return array|mixed
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:44
     */
    public function ScratchPrize(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/scratch-prize")->post($data);
    }

    /**
     * 券列表
     *
     * @param array $data
     *
     * @return array|mixed
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:44
     */
    public function CouponList(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/coupon-list")->post($data);
    }

    /**
     * 核销券
     *
     * @param array $data
     *
     * @return array|mixed
     * @throws Exception
     * @author  user <qiuapeng921@vchangyi.com>
     * @date    2022/8/12 16:43
     */
    public function WriteOffCoupon(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/write-off-coupon")->post($data);
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
        return Sdk::work()->setApi("service/v1/activity/scratch/user-view")->post($data);
    }

    /**
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function CheckCoupon(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/check-coupon")->post($data);
    }

    /**
     * @param array $data
     * @return array|mixed
     * @throws Exception
     */
    public function ShareInfo(array $data)
    {
        return Sdk::work()->setApi("service/v1/activity/scratch/share-info")->post($data);
    }
}