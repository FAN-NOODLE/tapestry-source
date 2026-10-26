<?php

namespace App\Sdk\Active\CRM\KsCore;

use App\Exceptions\Error;
use App\Sdk\Active\CRM\KsCore\Customer\CustomerQueryCoupons;
use App\Sdk\Active\CRM\KsCore\Customer\CustomerQueryInfo;
use App\Sdk\Active\CRM\KsCore\Customer\CustomerQueryInfoByMobile;
use App\Sdk\Active\CRM\KsCore\Customer\CustomerQueryInfoByOpenId;
use App\Sdk\Active\CRM\KsCore\User\UserCheck;
use App\Sdk\Active\CRM\KsCore\User\UserGetCouponList;
use App\Sdk\Active\CRM\KsCore\User\UserGetInfo;
use App\Sdk\Active\CRM\KsCore\User\UserGetInfoCommon;
use App\Sdk\Active\CRM\KsCore\User\UserGetOrderList;
use App\Sdk\Active\CRM\KsCore\User\UserJoin;

class UserRequestService
{
    private $stringCode;
    private $stringValue;
    private $customerCode;

    private $mobile;
    private $openId;
    private $unionId;
    public  $enterpriseId;

    /**
     * @var Request
     */
    public $request;

    public function __construct()
    {
    }

    public function setUser(array $user): UserRequestService
    {

        if (isset($user['open_id']) && !empty($user['open_id'])) {
            $this->openId      = $user['open_id'];
            $this->stringCode  = 'OpenId';
            $this->stringValue = $user['open_id'];
        }
        if (isset($user['mobile']) && !empty($user['mobile'])) {
            $this->mobile      = $user['mobile'];
            $this->stringCode  = 'Mobile';
            $this->stringValue = $user['mobile'];
        }
        if (isset($user['union_id']) && !empty($user['union_id'])) {
            $this->unionId     = $user['union_id'];
            $this->stringCode  = 'unionId';
            $this->stringValue = $user['union_id'];
        }
        if (isset($user['customer_code']) && !empty($user['customer_code'])) {
            $this->customerCode = $user['customer_code'];
            $this->stringCode   = 'userCode';
            $this->stringValue  = $user['customer_code'];
        }

        $this->enterpriseId = $user['enterprise_id'];
        $this->request      = new Request();

        return $this;
    }

    /**
     * 根据mobile获取会员信息
     *
     * @return array|null
     */
    public function queryInfoByMobile(): ?array
    {
        if ($this->mobile) {
            $request  = new CustomerQueryInfoByMobile($this->mobile);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 根据open_id获取会员信息
     *
     * @return array|null
     */
    public function queryInfoByOpenId(): ?array
    {
        if ($this->openId) {
            $request  = new CustomerQueryInfoByOpenId($this->openId);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData();
            }
        }
        return [];

    }

    /**
     * 根据mobile/open_id/union_id获取会员信息
     *
     * @return array|null
     */
    public function queryInfo(): ?array
    {
        $request  = new CustomerQueryInfo($this->mobile, $this->openId, $this->unionId);
        $response = $this->request->requestApi($request);
        if ($response && $response->isOk()) {
            return $response->getData();
        }
        return [];
    }


    /**
     * 获取会员优惠券列表信息
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/22
     */
    public function queryCustomerCoupons(): array
    {
        $request  = new CustomerQueryCoupons($this->mobile, $this->openId);
        $response = $this->request->requestApi($request);
        if ($response && $response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    /**
     * 通用会员信息查询接口
     *
     * @return array|null
     */
    public function getUserInfoCommon(): ?array
    {
        if ($this->stringCode && $this->stringValue) {
            $request  = new UserGetInfoCommon($this->stringCode, $this->stringValue);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 根据union_id获取会员信息
     *
     * @return array|null
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/22
     */
    public function getUserInfo(): ?array
    {
        if ($this->unionId) {
            $request  = new UserGetInfo($this->unionId);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 通过unionid获取用户订单列表信息
     *
     * @param $index
     * @param $size
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/22
     */
    public function getUserOrderListByUnionId($index, $size): array
    {
        if ($this->unionId) {
            $request  = new UserGetOrderList($this->unionId, $index, $size);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData() ?? [];
            }
        }
        return [];

    }

    /**
     * 通过unionid获取用户优惠券列表信息
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/22
     */
    public function getUserCouponList(): array
    {
        if ($this->unionId) {
            $request  = new UserGetCouponList($this->unionId);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getData() ?? [];
            }
        }
        return [];

    }

    /**
     * 用户是否是会员接口
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/29
     */
    public function checkMember(): array
    {
        if ($this->unionId) {
            $request  = new UserCheck($this->unionId);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
                return $response->getDataReq() ?? [];
            }
        }

        return ['data' => false];

    }

    public function userJoin($campaignCode, $shopCode, $clerkCode): array
    {
        $request  = new UserJoin($this->unionId);
        $param    = [
            'form_params' =>
                [
                    'UnionId'      => $this->unionId,
                    'CampaignCode' => $campaignCode,
                    'ShopCode'     => $shopCode,
                    'ClerkCode'    => $clerkCode
                ]
        ];
        $response = $this->request->requestApi($request, "POST", $param);
        if ($response) {
            return $response->getDataReq() ?? [];
        }

        return ['data' => false];

    }

}