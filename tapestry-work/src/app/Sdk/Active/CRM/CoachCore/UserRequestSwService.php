<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Exceptions\Error;
use App\Models\Scrm\Members;
use App\Sdk\Active\CRM\CoachCore\Auth\TokenCache;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerCouponSearch;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerPointHistory;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerPurchaseSearch;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerSaAssignment;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerSaUnAssignment;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerSearchInfo;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerSearchInfoByMobile;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerSearchInfoByUnionId;
use App\Sdk\Active\CRM\CoachCore\CustomerSw\CustomerTagSearch;
use Exception;
use Illuminate\Support\Facades\Log;

class UserRequestSwService
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

    public function setUser(array $user): UserRequestSwService
    {
        if (isset($user['mobile']) && !empty($user['mobile'])) {
            $this->mobile = $user['mobile'];
            $this->stringCode = 'mobile';
            $this->stringValue = $user['mobile'];
        }
        if (isset($user['union_id']) && !empty($user['union_id'])) {
            $this->unionId = $user['union_id'];
            $this->stringCode = 'union_id';
            $this->stringValue = $user['union_id'];
        }
        // sw 针对计划任务 内部传参是 customer_code 对接cdp 传参是 customer_code
        if (isset($user['customer_code']) && !empty($user['customer_code'])) {
            $this->stringCode = 'customer_code';
            $this->stringValue = $user['customer_code'];
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
     * @throws Exception
     */
    public function queryInfoByUnionId(): ?array
    {
        if ($this->unionId) {
            $request = new CustomerSearchInfoByUnionId($this->unionId);
            $response = $this->request->requestApi($request);
            if ($response && $response->isOk()) {
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
            if ($response && $response->isOk()) {
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
        if ($response && $response->isOk()) {
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
     *
     * @return array
     * @throws Exception
     */
    public function queryCustomerCouponByCustomerInfo(array $customerInfo, int $status, string $channel_type, int $source): array
    {
        $request = new CustomerCouponSearch($customerInfo, $status, $channel_type, $source);
        $response = $this->request->requestApi($request);
        if ($response && $response->isOk()) {
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
     *
     * @return array
     * @throws Exception
     */
    public function queryCustomerPurchaseByCustomerId($customerId, $startDateTime, $endDateTime): array
    {
        $request = new CustomerPurchaseSearch($customerId, $startDateTime, $endDateTime);
        $response = $this->request->requestApi($request);
        if ($response && $response->isOk()) {
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
        if ($response && $response->isOk()) {
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
     */
    public function customerSaAssignment($params): array
    {
        Log::info("sw绑定好友请求过来的参数" . json_encode($params));
        if (empty($params['storecode'])) {
            return [];
        }
        if (!empty($params['saAccount'])) {
            //sw需要获取sa的真实工号
            $saCode = Members::query()->where("enterprise_id",$this->enterpriseId)
                ->where("userid",$params['saAccount'])
                ->value('open_userid');
            if (!empty($saCode)) $params['saAccount'] = $saCode;
        }
        $request = new CustomerSaAssignment($params);
        $response = $this->request->requestApi($request);
        Log::info("请求SW的绑定好友接口" . json_encode($response->getBody()));
        if ($response->isOk()) {
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
     */
    public function customerSaUnAssignment($params): array
    {
        Log::info("sw取消绑定请求过来的参数" . json_encode($params));
        if (!empty($params['saAccount'])) {
            //sw需要获取sa的真实工号
            $saCode = Members::query()->where("enterprise_id", $this->enterpriseId)
                ->where("userid",$params['saAccount'])
                ->value('open_userid');
            if (!empty($saCode)) $params['saAccount'] = $saCode;
        }
        $request = new CustomerSaUnAssignment($params);
        $response = $this->request->requestApi($request);
        Log::info("请求SW的取消绑定好友接口" . json_encode($response->getBody()));
        if ($response->isOk()) {
            return $response->getResData() ?? [];
        }
        throw new Error($response->getResCode(), $response->getMsg());

    }

}