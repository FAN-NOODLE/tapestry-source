<?php
/**
 * CoachCrm class
 */

namespace App\Sdk\Active\CRM;

use App\Exceptions\Error;
use App\Models\StoreEmployee;
use App\Sdk\Active\CRM\CoachCore\Config;
use App\Sdk\Active\CRM\CoachCore\UserRequestService;
use App\Sdk\Active\CRM\CoachCore\WelfareRequestService;
use DateTime;
use Exception;
use ReflectionException;
use ReflectionMethod;
use Throwable;

class CoachCrm extends BaseCrm
{

    private $enterpriseId;
    private $userRequestService;
    private $welfareRequestService;

    public function __construct($enterpriseConfig)
    {

        $this->enterpriseId = $enterpriseConfig['enterprise_id'];
        $this->userRequestService = new UserRequestService();
        $this->welfareRequestService = new WelfareRequestService();

    }

    public function getToken()
    {
        return $this->userRequestService->getToken($this->enterpriseId);
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
        $service = $this->userRequestService->setUser($user);

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

        $service = $this->userRequestService->setUser($user);
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
        $service = $this->userRequestService->setUser($user);

        try {
            $customerInfo = $this->getMemberInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if (!$customerInfo['status']) {
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

    /*
     * 处理会员信息返回格式
     */

    public function getMemberInfo(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        $res = $service->queryInfoByUnionId();
        if (empty($res)) {
            return parent::reqResponse($user);
        }
        if (empty($res['records'])) {
            return parent::reqMemberInfoNo($user);
        }

        $result = $this->handleCustomerInfo($res['records'][0]);
        if ($result['status']) {
            // 查询导购门店code
            $memberWhere = ['mobile' => $result['mobile'],
                'enterprise_id' => $this->enterpriseId];
            $storeCode = StoreEmployee::query()
                    ->where($memberWhere)
                    ->value('third_store_code') ?? '1';

            $result['vipStatus'] = [];
            if (!empty($result['membership'])) {
                //如果membership不为空获取vip信息
                $paramsDetail = [
                    'enterprise_id' => $this->enterpriseId,
                    'params' =>
                        ['storeCode' => $storeCode, 'memberId' => $result['membership']]
                ];
                try {
                    // 获取会员vip信息
                    $customerDetail = $this->customerDetail($paramsDetail);
                    if (isset($customerDetail['resultCode']) && $customerDetail['resultCode'] == 'SUCCESS') {
                        $result['vipStatus'] = $customerDetail['vipStatus'];
                    }
                } catch (Throwable $e) {

                }

            }
        }

        return $result;
    }

    public function handleCustomerInfo($info): array
    {
        $codeDescription = ['O' => '非会员',
            'PROSPECT' => '注册会员',
            'C' => '普通会员',
            'D' => '商场会员',
            'GOLD' => '金卡会员',
            'SILVER' => '银卡会员',
            'DIAMOND' => '钻卡会员'];
        $cardDescription = ['O' => '0',
            'PROSPECT' => '1',
            'C' => '2',
            'D' => '3',
            'GOLD' => '4',
            'SILVER' => '5',
            'DIAMOND' => '6'];

        return [
            'status' => true,
            'id' => $info['id'],              //String	记录编号
            'customer_code' => $info['customer_code'],   //String	会员编号，用于识别会员
            'customer_type1' => $info['customer_type'],   //String	会员等级
            'customer_level' => $info['customer_type'],   //String	会员等级
            'customer_type' => $codeDescription[$info['customer_type']] ?? '',   //String	会员等级
            'customer_card_str' => $codeDescription[$info['customer_type']] ?? '',   //String	会员等级名字
            'customer_card' => $cardDescription[$info['customer_type']] ?? '0',   //String	会员等级类型
            'union_id' => $info['union_id'],        //String	微信unionId
            'mobile' => !empty($info['mobile']) ? $info['mobile'] : '',   //String    会员手机号
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
            'membership' => $info['membership'] ?? '',      //String	主卡卡号
            'client_type' => $info['client_type'] ?? '',     //String	CTL 类型
            'total_valid_points' => $info['total_valid_points'] ?? 0.00,     //Float	积分
            'pts_need_to_upgrade' => (float)isset($info['pts_need_to_upgrade']) ? $info['pts_need_to_upgrade'] : 0.00, // 升级所需消费金额
            'pts_need_to_renew' => (float)isset($info['pts_need_to_renew']) ? $info['pts_need_to_renew'] : 0.00, // 保级所需消费金额
            'msg' => '会员信息返回成功'
        ];

    }

    public function customerDetail(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        return $service->customerDetail($user['params']);

    }

    /**
     * 根据customerInfo查询会员优惠卷记录
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerCoupon(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        try {
//            $customerInfo = $this->getMemberInfo($user);
            $customerInfo = $this->getCustomerInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }
        if (!$customerInfo['status']) {
            return $customerInfo;
        }

        $data = $service->queryCustomerCouponByCustomerInfo($customerInfo, $user['status'], $user['channel_type'], $user['source'], $user['page'], $user['size']);

        $channelTypeArr = ['RETAIL' => '零售店', 'OUTLET' => '奥莱店'];
        $typeArr = ['1' => '折扣券', '2' => '代金券', '3' => '礼品兑换券'];
        $statusArr = ['1' => '未使用', '2' => '已使用', '3' => '过期', '4' => '冻结'];
        if (is_array($data['records']) && !empty($data['records'])) {
            foreach ($data['records'] as &$v) {

                $v['channel_type_str'] = $channelTypeArr[$v['channel_type']] ?? '其他券';
                $v['type_str'] = $typeArr[$v['type']] ?? '';
                $v['status_str'] = $statusArr[$v['status']] ?? '';
            }
        }

        return $data;

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
        $service = $this->userRequestService->setUser($user);

        $res = $service->getCustomerInfo();
        if (empty($res)) {
            return parent::reqResponse($user);
        }
        if (empty($res['records'])) {
            return parent::reqMemberInfoNo($user);
        }
        return $this->handleCustomerInfo($res['records'][0]);

    }

    /**
     * 通用查询会员基本信息（原始版，不做数据处理）
     *
     * @param array $user
     *
     * @return array
     * @throws Exception
     */
    public function getCustomerInfoBySearchV2(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        $res = $service->getCustomerInfo();
        if (empty($res) || empty($res['records'])) {
            return [];
        }
        $return = $res['records'][0];
        if ($this->checkCoachTime() && in_array($return['customer_type'], ['PROSPECT', 'C'])) {
            $return['expiry_date'] = 0;
        }
        $v1Res = $service->getCustomerInfoV1();
        if (!empty($v1Res['records']) && !empty($v1Res['records'][0])) {
            $return['registration_date'] = $v1Res['records'][0]['registration_date'] ?? '';
            $return['first_channel_code'] = $v1Res['records'][0]['first_channel_code'] ?? '';
        }

        return $return;

    }

    /**
     * 会员详情等级信息
     *
     * @param $membership
     * @return array
     */
    public function getCustomerVip($membership): array
    {
        try {
            // 获取会员vip信息
            $customerDetail = $this->userRequestService->customerDetail(['storeCode' => '1', 'memberId' => $membership]);
            if (isset($customerDetail['resultCode']) && $customerDetail['resultCode'] == 'SUCCESS') {
                return $customerDetail['vipStatus'];
            }
        } catch (Exception $e) {
        }
        return [];
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
        $service = $this->userRequestService->setUser($user);

        try {
//            $customerInfo = $this->getMemberInfo($user);
            $customerInfo = $this->getCustomerInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if (!$customerInfo['status']) {
            return $customerInfo;
        }

        $result = $service->queryCustomerPurchaseByCustomerId($customerInfo['id'], $user['startDateTime'], $user['endDateTime'], $user['page'], $user['size']);

        $L12MStartTime = strtotime('-12month');
        $L12MPurchaseAmount = 0;
        $L12MPurchaseCount = 0;
        $lastPurchaseTime = '';
        $data = [];
        if (isset($result['records']) && count($result['records']) > 0) {
            $lastPurchaseTime = date('Y-m-d', $result['records'][0]['purchase_date'] / 1000);
            foreach ($result['records'] as $item) {
                // 近12个月消费金额和数量
                if ($item['purchase_date'] / 1000 >= $L12MStartTime) {
                    $L12MPurchaseAmount += $item['amount'];
                    $L12MPurchaseCount += 1;
                }
                if (empty($item['purchase_skus'])) {
                    continue;
                }
                foreach ($item['purchase_skus'] as $value) {
                    $v['mobile'] = $item['mobile'] ?? '';
                    $v['membership'] = $item['membership'] ?? '';
                    $v['customer_id'] = $item['customer_id'] ?? '';
                    $v['store_name'] = $item['store_name'] ?? '';
                    $v['store_code'] = $item['store_code'] ?? '';
                    $v['ssociate'] = $item['ssociate'] ?? '';
                    $v['ssociate_name'] = $item['ssociate_name'] ?? '';
                    $v['purchase_id'] = $item['purchase_id'] ?? '';
                    $v['order_no'] = $item['order_no'] ?? '';
                    $v['amount'] = $value['amount'] ?? '';
                    $v['sku_code'] = $value['sku_code'] ?? '';
                    $v['sku_name'] = $value['sku_name'] ?? '';
                    $v['color'] = $value['color'] ?? '';
                    $v['style'] = $value['style'] ?? '';
                    $v['purchase_date_time'] = date('Y-m-d', $item['purchase_date'] / 1000) ?? '';
                    $v['purchase_date'] = $item['purchase_date'] / 1000;

                    // 拼接image-url 格式
                    // 首先去掉头尾空格
                    $str = str_replace('/', '', trim($v['style'] . '_' . $v['color']));
                    $v['sku_image'] = Config::getSkuImageUrl() . '/' . strtolower(str_replace(' ', '', $str)) . '_a0';
                    $data['records'][] = $v;
                }
            }
        }
        $data['pageNum'] = $result['pageNum'] ?? 0;
        $data['pageSize'] = $result['pageSize'] ?? 0;
        $data['isFirstPage'] = $result['isFirstPage'] ?? '';
        $data['isLastPage'] = $result['isLastPage'] ?? '';
        $data['hasNextPage'] = $result['hasNextPage'] ?? '';
        $data['total'] = $result['total'] ?? 0;
        $data['pages'] = $result['pages'] ?? 0;
        $data['L12MPurchaseAmount'] = (string)$L12MPurchaseAmount;
        $data['L12MPurchaseCount'] = $L12MPurchaseCount;
        $data['lastPurchaseTime'] = $lastPurchaseTime;

        return $data;

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
        $service = $this->userRequestService->setUser($user);

        try {
            $customerInfo = $this->getMemberInfo($user);

        } catch (Throwable $e) {
            throw new Error(10000, "查询会员信息失败");
        }

        if (!$customerInfo['status']) {
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
     * @param array $user
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareHaved(array $user): array
    {
        $res = [];
        if (empty($user['union_id'])) {
            return $res;
        }
        $service = $this->welfareRequestService->setOption($user);

        $havedInfo = $service->actionWelfareHaved(Config::campaignSID, $user['union_id']);
        if (is_array($havedInfo) && !empty($havedInfo)) {
            $res['welfare_info'] = $havedInfo['WelfareInfo'];
            $res['haved'] = $havedInfo['Haved'];
        }
        if ($res['haved'] == false && empty($res['welfare_info'])) {
            // 如果没发过劵，直接调用发劵接口
            $data = $service->actionWelfareGive(Config::campaignSID, $user['union_id'], $user['extUserid']);

            $res['welfare_info'] = $data['Comm_PrizeUserInfo'];
            $res['haved'] = false;
        }

        return $res;

    }

    /**
     * @param array $user
     *
     * @return array
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/7
     */
    public function actionWelfareGive(array $user): array
    {
        $service = $this->welfareRequestService->setOption($user);

        return $service->actionWelfareGive($user['campaignSID'], $user['union_id'], $user['extUserid']);

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
        $service = $this->userRequestService->setUser($user);

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
        $service = $this->userRequestService->setUser($user);

        return $service->customerSaUnAssignment($user['params']);

    }


    private function checkCoachTime(): bool
    {
        $targetTime = new DateTime('2025-01-01 00:00:00');
        // 当前时间
        $now = new DateTime();

        // 判断当前时间是否超过目标时间
        return $now > $targetTime;
    }

    /**
     * 会员聚合筛选条件
     *
     * @return string[][][]
     */
    public function customerOption(): array
    {

        $returnArr = [
            'type' => [
                ['key' => '0', 'value' => "全部", 'valueStr' => "All"],
                ['key' => '1', 'value' => "Top Client", 'valueStr' => "Top Client"],
                ['key' => '2', 'value' => "New High Potential", 'valueStr' => "New High Potential"],
                ['key' => '3', 'value' => "TC or NHP", 'valueStr' => "Top Client or New High Potential"]],
            'vipLevel' => [
                ['key' => '0', 'value' => '全部', 'valueStr' => 'All'],
                ['key' => '1', 'value' => '银卡', 'valueStr' => 'SILVER'],
                ['key' => '2', 'value' => '金卡', 'valueStr' => 'GOLD'],
                ['key' => '3', 'value' => '钻石卡', 'valueStr' => 'DIAMOND']],
            'priceRange' => [
                ['key' => '0', 'value' => '全部', 'valueStr' => 'All'],
                ['key' => '7', 'value' => '30,000及以上'],
                ['key' => '6', 'value' => '20,000-29,999(包含)'],
                ['key' => '5', 'value' => '10,000-19,999(包含)'],
                ['key' => '4', 'value' => '8,000-9,999(包含)'],
                ['key' => '3', 'value' => '5,000-7,999(包含)'],
                ['key' => '2', 'value' => '3,000-4,999(包含)'],
                ['key' => '1', 'value' => '3,000以下'],
            ],
            'department' => [
                ['key' => '0', 'value' => '全部', 'valueStr' => 'All'],
                ['key' => 'D01', 'value' => '女式包包', 'valueStr' => 'WOMENS BAGS'],
                ['key' => 'D02', 'value' => '女式钱包', 'valueStr' => 'WOMENS WALLETS'],
                ['key' => 'D03', 'value' => '女式小皮件', 'valueStr' => 'WOMENS SLGs'],
                ['key' => 'D04', 'value' => '男式包包', 'valueStr' => 'MENS BAGS'],
                ['key' => 'D05', 'value' => '男式钱包', 'valueStr' => 'MENS WALLETS'],
                ['key' => 'D06', 'value' => '女式旅行用品', 'valueStr' => 'WOMENS TRAVEL'],
                ['key' => 'D07', 'value' => '女式成衣', 'valueStr' => 'WOMENS READY-TO-WEAR'],
                ['key' => 'D08', 'value' => '男式皮带及配件', 'valueStr' => 'MENS SPECIALTY ACCESSORIES'],
                ['key' => 'D09', 'value' => '女式手表', 'valueStr' => 'WOMENS WATCHES'],
                ['key' => 'D10', 'value' => '男式旅行产品', 'valueStr' => 'MENS TRAVEL'],
                ['key' => 'D11', 'value' => '女式鞋履', 'valueStr' => 'WOMENS SHOES'],
                ['key' => 'D12', 'value' => '女式眼镜', 'valueStr' => 'WOMENS EYEWEAR'],
                ['key' => 'D13', 'value' => '珠宝', 'valueStr' => 'JEWELRY'],
                ['key' => 'D14', 'value' => '其他', 'valueStr' => 'OTHER'],
                ['key' => 'D15', 'value' => '小孩产品', 'valueStr' => 'KIDS/BABY'],
                ['key' => 'D16', 'value' => '男式成衣', 'valueStr' => 'MENS READY-TO-WEAR'],
                ['key' => 'D17', 'value' => '女式香水', 'valueStr' => 'WOMENS FRAGRANCE'],
                ['key' => 'D18', 'value' => '男式手表', 'valueStr' => 'MENS WATCHES'],
                ['key' => 'D19', 'value' => '男式鞋履', 'valueStr' => 'MENS SHOES'],
                ['key' => 'D20', 'value' => '男式眼镜', 'valueStr' => 'MENS EYEWEAR'],
                ['key' => 'D21', 'value' => '男式香水', 'valueStr' => 'MENS FRAGRANCE'],
                ['key' => 'D22', 'value' => '女式围巾', 'valueStr' => 'WOMENS SOFT ACCESSORIES'],
                ['key' => 'D23', 'value' => '油膏及配饰', 'valueStr' => 'WOMENS SPECIALTY ACCESSORIES'],
            ],
            'gender' => [
                ['key' => '0', 'value' => '全部', 'valueStr' => 'All'],
                ['key' => '1', 'value' => '女性', 'valueStr' => 'Female'],
                ['key' => '2', 'value' => '男性', 'valueStr' => 'Male'],
                ['key' => '3', 'value' => '不确定', 'valueStr' => 'Unknown'],
            ],
            'birthMonth' => [
                ['key' => '00', 'value' => '全部', 'valueStr' => 'All'],
                ['key' => '01', 'value' => '1', 'valueStr' => 'January'],
                ['key' => '02', 'value' => '2', 'valueStr' => 'February'],
                ['key' => '03', 'value' => '3', 'valueStr' => 'March'],
                ['key' => '04', 'value' => '4', 'valueStr' => 'April'],
                ['key' => '05', 'value' => '5', 'valueStr' => 'May'],
                ['key' => '06', 'value' => '6', 'valueStr' => 'June'],
                ['key' => '07', 'value' => '7', 'valueStr' => 'July'],
                ['key' => '08', 'value' => '8', 'valueStr' => 'August'],
                ['key' => '09', 'value' => '9', 'valueStr' => 'September'],
                ['key' => '10', 'value' => '10', 'valueStr' => 'October'],
                ['key' => '11', 'value' => '11', 'valueStr' => 'November'],
                ['key' => '12', 'value' => '12', 'valueStr' => 'December'],
            ],

        ];
        if ($this->checkCoachTime()) {
            unset($returnArr['type']);
            $returnArr['vipLevel'][0] = ['key' => '0', 'value' => '全部已购', 'valueStr' => 'All'];
        }
        return $returnArr;

    }

    public function customerListWeclienteling(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        return $service->customerListWeclienteling($user['params']);

    }

}