<?php


namespace App\Services;

namespace App\Services\Crm;

use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CrmFactory;
use App\Sdk\Active\CrmHelper;
use App\Services\LogService;
use Exception;
use Illuminate\Support\Facades\Log;

class CrmService
{

    /**
     * 企业对接第三方信息
     *
     * @var |null
     */
    private $enterpriseConfig;

    /**
     * CRM SDK实例
     */
    private $sdk;

    public function __construct($enterpriseId)
    {
        $this->enterpriseConfig = CrmHelper::getEnterpriseConfig($enterpriseId);

        $this->sdk = CrmFactory::getInstance($this->enterpriseConfig);
//        $this->sdk = (new CoachCrm($this->enterpriseConfig));
//        $this->sdk = (new KsCrm($this->enterpriseConfig));
    }

    /**
     * 测试获取token
     *
     * @return mixed
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getToken()
    {
        return $this->sdk->getToken();
    }

    /**
     * 通用获取第三方会员信息接口
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2023/3/17
     */
    public function getCustomerInfo($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->getCustomerInfo($input);

            LogService::info('获取第三方会员信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('获取第三方会员信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);
            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;
    }

    /**
     * 获取第三方会员信息接口
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getMemberInfo($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->getMemberInfo($input);

            LogService::info('获取第三方会员信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('获取第三方会员信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        // 响应数据
        return $result;
    }


    /**
     * cdpCommon
     */
    public function cdpCommon($input)
    {
        $method = $input['method'];
        unset($input['method']);
        try {
            // 调用SDK方法
            return $this->sdk->{$method}($input);
        } catch (Exception $e) {
            LogService::error('CDP接口', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            error("三方接口请求报错" . $e->getMessage());
        }

    }

    /**
     * 查询会员积分记录
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerPointHistory($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->getCustomerPointHistory($input);

            LogService::info('查询会员积分记录信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询会员积分记录信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }

    /**
     * 查询会员优惠卷记录
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerCoupon($input): array
    {
        $input['page'] = isset($input['page']) && $input['page'] > 0 ? $input['page'] : 1;
        $input['size'] = isset($input['size']) && $input['size'] > 0 ? $input['size'] : 99;
        $result = [];

        try {
            // 调用SDK方法
            $result = $this->sdk->getCustomerCoupon($input);

            LogService::info('查询会员优惠卷记录信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询会员优惠卷记录信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }

    /**
     * 查询会员购买记录\订单
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerPurchase($input): array
    {
        $result = [];

        $input['startDateTime'] = $input['startDateTime'] ?? 0;
        $input['endDateTime'] = $input['endDateTime'] ?? time() * 1000;

        $input['page'] = isset($input['page']) && $input['page'] > 0 ? $input['page'] - 1 : 0;
        $input['size'] = isset($input['size']) && $input['size'] > 0 ? $input['size'] : 999;
        try {
            // 调用SDK方法
            $result = $this->sdk->getCustomerPurchase($input);

            LogService::info('查询会员购买记录信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询会员购买记录信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        return $result;
    }

    /**
     * 查询会员标签记录
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function getCustomerTag($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->getCustomerTag($input);

            LogService::info('查询会员标签记录信息', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询会员标签记录信息异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        return $result;

    }

    /**
     * 用户是否是会员
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function checkMember($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->checkMember($input);

            LogService::info('查询用户是否会员', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询用户是否会员异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        return $result;

    }

    /**
     * coach发放优惠券
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareGive($input): array
    {
        $result = [];
        if ($this->enterpriseConfig['name'] != 'coach') {
            return $result;
        }
        $input['campaignSID'] = Config::campaignSID;
        try {
            // 调用SDK方法
            $result = $this->sdk->actionWelfareGive($input);

            LogService::info('查询用户优惠卷', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('优惠卷异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        return $result;

    }

    /**
     * coach是否发放过优惠券
     *
     * @param     $input
     * @param int $i
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareHaved($input, $i = 3): array
    {
        $result = [];
        if ($this->enterpriseConfig['name'] != 'coach' && $this->enterpriseConfig['name'] != 'ks') {
            return $result;
        }
        try {
            // 调用SDK方法
            if ($this->enterpriseConfig['name'] == 'coach') {
                $result = $this->sdk->actionWelfareHaved($input);
            } elseif ($this->enterpriseConfig['name'] == 'ks') {
                $result = $this->sdk->userJoin($input);
            } else {
                return $result;
            }

            LogService::info('查询用户优惠卷', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {
            if ($e->getCode() == 502 && $i > 0) {//请求接口502的状态添加重试机制
                $i = $i - 1;
                $result = $this->actionWelfareHaved($input, $i);
            } else {
                LogService::error('优惠卷异常', [
                    'enterprise' => $this->enterpriseConfig['name'],
                    'message' => $e->getMessage(),
                    'line' => $e->getLine(),
                    'file' => $e->getFile(),
                ]);

                if ($e->getMessage() == "参数有误") {
                    return $result;
                }
                error("三方接口请求报错" . $e->getMessage());
            }

        }

        return $result;

    }

    /**
     * ks是否发放过优惠券
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionUserJoin($input): array
    {
        $result = [];
        if ($this->enterpriseConfig['name'] != 'ks') {
            return $result;
        }
        try {
            // 调用SDK方法
            $result = $this->sdk->userJoin($input);

            LogService::info('查询用户优惠卷', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('优惠卷异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }

        return $result;

    }

    /**
     * 导购添加顾客
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function saAssignment($input): array
    {
        $result = [];
        Log::info("好友关系同步CDP请求参数" . json_encode($input));
        try {
            // 调用SDK方法
            $result = $this->sdk->saAssignment($input);
            
            LogService::info('导购添加顾客', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('导购添加顾客异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }

    /**
     * 导购删除顾客
     *
     * @param $input
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/24
     */
    public function saUnAssignment($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->saUnAssignment($input);

            LogService::info('导购删除顾客', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('导购删除顾客异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }

    /**
     * 会员聚合筛选条件
     *
     * @param $input
     *
     * @return array
     */
    public function customerOption($input): array
    {
        return $this->sdk->customerOption($input);
    }

    /**
     * 根据门店或门店导购批量查询会员信息
     *
     * @param $input
     *
     * @return array
     */
    public function customerListWeclienteling($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->customerListWeclienteling($input);

            LogService::info('查询聚合会员', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询聚合会员异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }

    /**
     * 会员详情
     *
     * @param $input
     *
     * @return array
     */
    public function customerDetail($input): array
    {
        $result = [];
        try {
            // 调用SDK方法
            $result = $this->sdk->customerDetail($input);

            LogService::info('查询会员详情', [
                '请求数据' => $input,
                '响应数据' => $result
            ]);
        } catch (Exception $e) {

            LogService::error('查询会员详情异常', [
                'enterprise' => $this->enterpriseConfig['name'],
                'message' => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
            ]);

            if ($e->getMessage() == "参数有误") {
                return $result;
            }
            error("三方接口请求报错" . $e->getMessage());
        }
        return $result;

    }
}