<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Exceptions\Error;
use App\Sdk\Active\CRM\CoachCore\Auth\TokenCache;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerCouponSearch;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerDetail;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerListWeclienteling;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerPointHistory;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerPurchaseSearch;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSaAssignment;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSaUnAssignment;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSearchInfo;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSearchInfoByMobile;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSearchInfoByUnionId;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerSearchInfoV1;
use App\Sdk\Active\CRM\CoachCore\Customer\CustomerTagSearch;
use Exception;

class UserRequestService
{
    public $enterpriseId;
    /**
     * @var Request
     */
    public $request;
    private $stringCode;
    private $stringValue;
    private $mobile;
    private $unionId;

    public function __construct()
    {
    }

    public function setUser(array $user): UserRequestService
    {
        if (isset($user['mobile']) && !empty($user['mobile'])) {
            $this->mobile = $user['mobile'];
            $this->stringCode = 'mobile';
            $this->stringValue = $user['mobile'];
        }
        // coach 针对计划任务 内部传参是 customer_code 对接cdp 传参是 membership
        if (isset($user['customer_code']) && !empty($user['customer_code'])) {
            $this->stringCode = 'membership';
            $this->stringValue = $user['customer_code'];
        }
        if (isset($user['union_id']) && !empty($user['union_id'])) {
            $this->unionId = $user['union_id'];
            $this->stringCode = 'union_id';
            $this->stringValue = $user['union_id'];
        }
//        if (isset($user['membership']) && !empty($user['membership'])) {
//            $this->stringCode = 'membership';
//            $this->stringValue = $user['membership'];
//        }
        if (empty($this->stringValue) && empty($user['params'])) {
            error("参数有误");
        }

        $this->enterpriseId = $user['enterprise_id'];
        $this->request = new Request($this->enterpriseId);

        return $this;
    }

    /**
     * 获取token
     *
     * @return mixed|string|null
     */
    public function getToken($enterpriseId)
    {
        return (new TokenCache($enterpriseId))->getAccessToken();
    }

    /**
     * 通用查询会员信息接口
     *
     * @return array|null
     * @throws Exception
     */
    public function getCustomerInfo(): ?array
    {
        if ($this->stringCode && $this->stringValue) {

            $request = new CustomerSearchInfo($this->stringCode, $this->stringValue);
            $response = $this->request->requestApi($request);
            if ($response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 通用查询会员信息接口
     *
     * @return array|null
     * @throws Exception
     */
    public function getCustomerInfoV1(): ?array
    {
        if ($this->stringCode && $this->stringValue) {

            $request = new CustomerSearchInfoV1($this->stringCode, $this->stringValue);
            $response = $this->request->requestApi($request);
            if ($response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 根据union_id获取会员信息
     *
     * @return array|null
     * @throws Exception
     */
    public function queryInfoByUnionId(): ?array
    {
        if ($this->unionId) {
            $request = new CustomerSearchInfoByUnionId($this->unionId);
            $response = $this->request->requestApi($request);
            if ($response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 根据mobile获取会员信息
     *
     * @return array|null
     * @throws Exception
     */
    public function queryInfoByMobile(): ?array
    {
        if ($this->mobile) {
            $request = new CustomerSearchInfoByMobile($this->mobile);
            $response = $this->request->requestApi($request);
            if ($response->isOk()) {
                return $response->getData();
            }
        }
        return [];
    }

    /**
     * 根据customerId查询会员积分记录
     *
     * @param $customerId
     * @param $startTime
     * @param $endTime
     *
     * @return array
     * @throws Exception
     */
    public function queryCustomerPointHistoryByCustomerId($customerId, $startTime, $endTime): array
    {
        $request = new CustomerPointHistory($customerId, $startTime, $endTime);
        $response = $this->request->requestApi($request);
        if ($response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    /**
     * 根据customerInfo查询会员优惠卷记录
     *
     * @param array $customerInfo
     * @param int $status
     * @param string $channel_type
     * @param int $source
     * @param        $page
     * @param        $size
     *
     * @return array
     * @throws Exception
     */
    public function queryCustomerCouponByCustomerInfo(array $customerInfo, int $status, string $channel_type, int $source, $page, $size): array
    {
        $request = new CustomerCouponSearch($customerInfo, $status, $channel_type, $source, $page, $size);
        $response = $this->request->requestApi($request);
        if ($response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    /**
     * 根据customerId查询会员购买记录
     *
     * @param $customerId
     * @param $startDateTime
     * @param $endDateTime
     * @param $page
     * @param $size
     * @return array
     * @throws Exception
     */
    public function queryCustomerPurchaseByCustomerId($customerId, $startDateTime, $endDateTime, $page, $size): array
    {
        $request = new CustomerPurchaseSearch($customerId, $startDateTime, $endDateTime, $page, $size);
        $response = $this->request->requestApi($request);
        if ($response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    /**
     * 根据customerId查询会员标签记录
     *
     * @param $customerId
     * @param $source
     *
     * @return array
     * @throws Exception
     */
    public function queryCustomerTagByCustomerId($customerId, $source): array
    {
        $request = new CustomerTagSearch($customerId, $source);
        $response = $this->request->requestApi($request);
        if ($response->isOk()) {
            return $response->getData() ?? [];
        }
        throw new Error($response->getCode(), $response->getMsg());

    }

    /**
     * Bind SA by unionId
     *
     * @param $params
     *
     * @return array
     * @throws Exception
     * @Time: 2022/11/18
     */
    public function customerSaAssignment($params): array
    {
        if (empty($params['storecode']) || $params['storecode'] == 'ESA') {
            return [];
        }
        $request = new CustomerSaAssignment($params);
        $response = $this->request->requestApi($request);
        if ($response->isOkCode()) {
            return $response->getResData() ?? [];
        }
        throw new Error($response->getResCode(), $response->getMsg());

    }

    /**
     * SA Force Change
     *
     * @param $params
     *
     * @return array
     * @throws Exception
     * @Time: 2022/11/18
     */
    public function customerSaUnAssignment($params): array
    {
        $request = new CustomerSaUnAssignment($params);
        $response = $this->request->requestApi($request);
        if ($response->isOkCode()) {
            return $response->getResData() ?? [];
        }
        throw new Error($response->getResCode(), $response->getMsg());

    }

    /**
     * xml格式请求
     * 会员列表聚合查询
     *
     * @param $params
     *
     * @return array
     * @throws Exception
     */
    public function customerListWeclienteling($params): array
    {
        $request = new CustomerListWeclienteling($params);
        $response = $this->request->requestApi($request, 'xml');
        if ($response->isSuccessCode()) {
            return $response->getResData() ?? [];
        }
        throw new Error($response->getResSuccessCode(), $response->getMsg());
    }

    /**
     * xml格式请求
     * 会员详情
     *
     * @param $params
     *
     * @return array
     * @throws Exception
     */
    public function customerDetail($params): array
    {
        $request = new CustomerDetail($params);
        $response = $this->request->requestApi($request, 'xml');
        if ($response->isSuccessCode()) {
            return $response->getResData() ?? [];
        }
        throw new Error($response->getResSuccessCode(), $response->getMsg());
    }

}