<?php

namespace App\Micros\Work;

use App\Micros\Common;

class Work
{
    const RedisCacheTime = 5 * 60;

    public static function addWorkerLoginLog($params)
    {

//        return Common::Work()->setApi('service/v1/department/lists')->post($params);

    }

    /**
     * 部门列表
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function thirdIdsList($params)
    {
        return Common::Work()->setApi('service/v1/department/list-by-scrm')->post($params);
    }

    /**
     * 会员聚合查询条件
     *
     * @param $params
     *
     * @return array|mixed
     */
    public static function customerOption($params)
    {
        return Common::Work()->setApi('service/crm/customer-option')->post($params);
    }

    public static function customerList($params)
    {
        return Common::Work()->setApi('service/crm/customer-list')->post($params);
    }

    /**
     * 获取 crm 顾客信息
     *
     * @param     $enterpriseId
     * @param     $unionId
     * @param int $euId
     *
     * @return array|bool|mixed
     * @Time: 2022/11/21
     */
    public function getCrmUserInfo($enterpriseId, $unionId, int $euId = 0, $customerCode = "")
    {
        // 查询crm
        $userReq['enterprise_id'] = $enterpriseId;
        $userReq['union_id']      = $unionId;
        $userReq['eu_id']         = $euId;
        $userReq['customer_code'] = $customerCode;
//        $cacheSuffix              = $unionId;
//        if (!empty($customerCode) && empty($unionId)) {
//            $cacheSuffix = $customerCode;
//        }
//
//        // redis 缓存用户数据
//        $redisKey = "CRM:" . $userReq['enterprise_id'] . ":user:info:" . $cacheSuffix;
//        $userCrm  = getOrSetData($redisKey);
//        if (!$userCrm) {

            $userCrm = $this->crmMemberInfo($userReq);
//            getOrSetData($redisKey, $userCrm, Work::RedisCacheTime);
//        }

        return $userCrm;

    }

    /**
     * crm 顾客信息
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function crmMemberInfo($params)
    {
        return Common::Work()->setApi('service/crm/customer-info')->post($params);
    }

    public function getCrmUserOrderList($params)
    {
        // redis 缓存用户数据
//        $redisKey = "CRM:" . $params['enterprise_id'] . ":user:order_list:" . $params['union_id'];
//        $userCrm  = getOrSetData($redisKey);
//        if (!$userCrm) {

            $userCrm = $this->crmOrderList($params);
//            getOrSetData($redisKey, $userCrm, Work::RedisCacheTime);
//        }

        return $userCrm;

    }

    /**
     * crm 顾客订单信息
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function crmOrderList($params)
    {
        return Common::Work()->setApi('service/crm/get-customer-purchase')->post($params);
    }

    public function getCrmUserCouponsList($params)
    {
        // redis 缓存用户数据
//        $redisKey = "CRM:" . $params['enterprise_id'] . ":user:coupons_list:" . $params['union_id'];
//        $userCrm  = getOrSetData($redisKey);
//        if (!$userCrm) {
//
            $userCrm = $this->crmCouponsList($params);
//            getOrSetData($redisKey, $userCrm, Work::RedisCacheTime);
//        }

        return $userCrm;

    }

    /**
     * crm 顾客卡劵信息
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function crmCouponsList($params)
    {
        return Common::Work()->setApi('service/crm/get-customer-coupon')->post($params);
    }


    /**
     * SW crm 顾客卡劵信息
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function swCrmCouponsList($params)
    {
        return Common::Work()->setApi('service/crm/get-sw-customer-coupon')->post($params);
    }


    /**
     * crm 顾客卡劵信息
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function customerDetailByCoachCrm($params)
    {
        return Common::Work()->setApi('service/crm/customer-detail-by-coach-crm')->post($params);
    }


    /**
     * cdp 推荐商品
     *
     * @param $params
     *
     * @return array|mixed
     * @Time: 2022/11/21
     */
    public static function getCustomerInterestedProducts($params)
    {
        return Common::Work()->setApi('service/cdp/get-customer-interested-products')->post($params);
    }
}

