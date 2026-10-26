<?php
/**
 * CoachCrm class
 */

namespace App\Sdk\Active\CRM;

use App\Exceptions\Error;
use App\Sdk\Active\CRM\CoachCore\UserRequestSwService;
use Exception;
use ReflectionException;
use ReflectionMethod;
use Throwable;

class SwCrm extends BaseCrm
{

    private $enterpriseId;
    private $userRequestSwService;

    public function __construct($enterpriseConfig)
    {

        $this->enterpriseId = $enterpriseConfig['enterprise_id'];
        $this->userRequestSwService = new UserRequestSwService();

    }

    public function getToken()
    {
        return $this->userRequestSwService->getToken($this->enterpriseId);
    }

    /**
     * 用户是否是会员接口
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function checkMember(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        $res = $service->queryInfoByUnionId();
        if (empty($res)) {
            return parent::reqResponse($user);
        }
        if (empty($res['records'])) {
            return parent::reqMemberInfoNo($user);
        }

        $info = $res['records'][0];

        return [
            'status' => true,
            'union_id' => $info['union_id']
        ];
    }

    /**
     * 通用查询会员基本信息
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerInfo(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        $res = $service->getCustomerInfo();
        if (empty($res)) {
            return parent::reqResponse($user);
        }
        if (empty($res['records'])) {
            return parent::reqMemberInfoNo($user);
        }
        return $this->handleCustomerInfo($res['records'][0]);

    }

    public function handleCustomerInfo($info): array
    {
        // 会员等级对应
        $codeDescription = [
            '' => '非会员',
            'Z' => '钻石「履」行家',
            'B' => '钻石「履」行家',
            'C' => '金牌「履」行家',
            'W' => '新星「履」行家',
        ];
        $cardDescription = [
            '' => '0',
            'Z' => '1',
            'B' => '2',
            'C' => '3',
            'W' => '4',
        ];

        return [
            'status' => true,
            'id' => $info['id'],              //String	记录编号
            'customer_code' => $info['customer_code'],   //String	会员编号，用于识别会员
            'customer_type1' => $info['customer_type'],   //String	会员等级
            'customer_level' => $info['customer_type'],   //String	会员等级
//            'customer_type'            => $codeDescription[$info['customer_type']] ?? '',   //String	会员等级
            'customer_type' => '',   //String	会员等级
            'customer_card' => $cardDescription[$info['customer_type']] ?? 0,   //String	会员等级类型
            'customer_card_str' => $codeDescription[$info['customer_type']] ?? '',   //String	会员等级名称
            'union_id' => $info['union_id'],        //String	微信unionId
            'mobile' => !empty($info['mobile']) ? $info['mobile'] : '',          //String	会员手机号
            'name' => $info['name'],            //String	姓名
            'birth_year' => $info['birth_year'],      //Int		YYYY 出生年份
            'birth_month' => $info['birth_month'],     //Int		MM: 出生月份，1~12
            'birth_day' => $info['birth_day'],       //Int		DD: 出生日，1~31
            'birthday' => ($info['birth_year'] . '-' . $info['birth_month'] . '-' . $info['birth_day']),
            'province' => $info['province'],        //String	省
            'city' => $info['city'],            //String	市
            'expiry_date1' => $info['expiry_date'],     //Int		会员有效期
            'expiry_date' => handleDateTimeISO($info['expiry_date']),
            'upgrade_date_time1' => $info['upgrade_date_time'],   //Date	最近一次升级/降级时间
            'upgrade_date_time' => handleDateTimeISO($info['upgrade_date_time']),
            'upgrade_need_points' => $info['upgrade_need_points'], //Float	升级需要积分
            'last_12m_purchase_amount' => $info['last_12m_purchase_amount'] ?? 0.00,   //Float	过去12个月消费总额
            'last_36m_purchase_amount' => $info['last_36m_purchase_amount'] ?? 0.00,   //Float	过去36个月消费总额
            'membership'               => $info['membership'],      //String	主卡卡号
            'client_type'              => $info['client_type'],     //String	CTL 类型
            'total_valid_points'       => $info['total_valid_points'] ?? 0.00,     //Float	积分
            'total_inta_valid_points'  => (int)isset($info['total_inta_valid_points']) ? $info['total_inta_valid_points'] : 0, // 可兑换积分
            'pts_need_to_upgrade'      => (float)isset($info['pts_need_to_upgrade']) ? $info['pts_need_to_upgrade'] : 0.00, // 升级所需消费金额
            'pts_need_to_renew'        => (float)isset($info['pts_need_to_renew']) ? $info['pts_need_to_renew'] : 0.00, // 保级所需消费金额
            'msg'                      => '会员信息返回成功'
        ];
    }

    /**
     * 直接调用接口，不按CRM的父类格式返回数据
     *
     * @param $method
     * @param $params
     *
     * @return array
     * @throws ReflectionException
     */
    public function __call($method, $params)
    {
        $user = $params[0];

        $service = $this->userRequestSwService->setUser($user);
        if (method_exists($service, $method)) {
            $method = new ReflectionMethod($service, $method);
            $method->setAccessible(true);
            unset($params[0]);
            return success(['data' => $method->invoke($service, $params)]);
        }

        throw new Error(10000, class_basename($service) . "::" . $method . " not exist");
    }

    /**
     * 根据customerId查询会员积分记录
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerPointHistory(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        try {
            $customerInfo = $this->getMemberInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if ($customerInfo['status'] == false) {
            return $customerInfo;
        }

        // point_date_time 积分生成时间查询
        $data = $service->queryCustomerPointHistoryByCustomerId($customerInfo['id'], $user['startTime'] ?? date("Y-m-d", strtotime("-1 year")), $user['endTime'] ?? date("Y-m-d", time()));

        if (is_array($data['records']) && !empty($data['records'])) {
            foreach ($data['records'] as &$v) {
                $v['point_date_time'] = handleDateTimeISO($v['point_date_time']);
                $v['creation_date_time'] = handleDateTimeISO($v['creation_date_time']);
                $v['expiry_start_date'] = handleDateTimeISO($v['expiry_start_date']);
                $v['expiry_end_date'] = handleDateTimeISO($v['expiry_end_date']);
            }
        }

        return $data;

    }

    /**
     * 查询会员基本信息
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getMemberInfo(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        $res = $service->queryInfoByUnionId();
        if (empty($res)) {
            return parent::reqResponse($user);
        }
        if (empty($res['records'])) {
            return parent::reqMemberInfoNo($user);
        }

        return $this->handleCustomerInfo($res['records'][0]);

    }

    /**
     * 根据customerId查询会员购买记录
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerPurchase(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        try {
            $customerInfo = $this->getMemberInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if ($customerInfo['status'] == false) {
            return $customerInfo;
        }

        return $service->queryCustomerPurchaseByCustomerId($customerInfo['id'], $user['startDateTime'], $user['endDateTime']);


    }

    /**
     * 根据customerId查询会员标签记录
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerTag(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        try {
            $customerInfo = $this->getMemberInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if ($customerInfo['status'] == false) {
            return $customerInfo;
        }

        $data = $service->queryCustomerTagByCustomerId($customerInfo['id'], $user['source']);
        if (is_array($data['records']) && !empty($data['records'])) {
            foreach ($data['records'] as &$v) {
                $v['update_date_time'] = handleDateTimeISO($v['update_date_time']);
            }
        }

        return $data;

    }


    /**
     * Bind SA by unionId
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function saAssignment(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        return $service->customerSaAssignment($user['params']);

    }

    /**
     * SA Force Change
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function saUnAssignment(array $user): array
    {
        $service = $this->userRequestSwService->setUser($user);

        return $service->customerSaUnAssignment($user['params']);

    }


}