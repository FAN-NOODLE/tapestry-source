<?php
/**
 * KsCrm class
 */

namespace App\Sdk\Active\CRM;

use App\Sdk\Active\CRM\KsCore\Config;
use App\Sdk\Active\CRM\KsCore\Request;
use App\Sdk\Active\CRM\KsCore\UserRequestService;
use Exception;

class KsCrm extends BaseCrm
{

    private $enterpriseId;
    private $userRequestService;

    public function __construct($enterpriseConfig)
    {

        $this->enterpriseId       = $enterpriseConfig['enterprise_id'];
        $this->userRequestService = new UserRequestService();

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

        $info = $service->getUserInfoCommon();
        if (empty($info)) {
            return parent::reqMemberInfoNo($user);
        }

        return $this->handleCustomerInfo($info, $user);

    }

    public function handleCustomerInfo($info, $user): array
    {
        /**
         * 参数名    类型    说明
         * UserCode    string    会员编号
         * Address    string    用户详细地址
         * Year    Int    年
         * Month    Int    月
         * Day    Int    日
         * Mobile    String    手机号
         * AreaCode    string    区号
         * LevelName    String    会员等级描述
         * Level    String    会员等级
         * LevelExpireTime    DateTime    会员到期时间
         * LatelyUpgradeTime    DateTime    会员最近升级时间
         * ShoppingMoney    decimal    3年内总消费金额（从当前天往前计算
         * UpgradePoint    int    升级所需积分
         * Name    String    会员姓名
         * Gender    String    男 女
         * AvailablePoint    Int    可用积分
         */
        $codeDescription = [
            ''  => '非会员',
            '0' => '注册粉丝',
            '1' => '粉蜜会员',
            '2' => '银卡',
            '3' => '金卡',
            '4' => '钻卡',
        ];
        return [
            'status'                   => true,
            'customer_code'            => $info['UserCode'],    //string	會員卡号。可含有數字和字母
            'customer_type'            => $info['LevelName'],   //String	会员等级信息
            'customer_level'           => (string)$info['Level'],       //String	会员等级
            'customer_card'            => (string)$info['Level'],       //String	会员等级类型
            'customer_card_str'        => $info['LevelName'] ?? '',       //String	会员等级类型名称
            'union_id'                 => $user['union_id'] ?? '',    //String	微信unionId
            'mobile'                   => $info['Mobile'],      //String	手机号
            'name'                     => $info['Name'],        //String	姓名
            'birth_year'               => $info['Year'],        //Int       YYYY 出生年份
            'birth_month'              => $info['Month'],       //Int		MM: 出生月份，1~12
            'birth_day'                => $info['Day'],         //Int		DD: 出生日，1~31
            'birthday'                 => ($info['Year'] . '-' . $info['Month'] . '-' . $info['Day']),
            'expiry_date'              => handleDateTimeUTCToDay($info['LevelExpireTime']), // 会员到期时间
            'upgrade_date_time'        => handleDateTimeUTCToDay($info['LatelyUpgradeTime']), // 会员最近升级时间
            'upgrade_need_points'      => $info['UpgradePoint'],    //Float	升级需要积分
            'last_12m_purchase_amount' => 0.00,   //Float 过去12个月消费总额
            'last_36m_purchase_amount' => $info['ShoppingMoney'] ?? 0.00,   //Float 过去36个月消费总额
            'total_valid_points'       => $info['AvailablePoint'] ?? 0.00,  //Float 积分
            'address'                  => $info['Address'] ?? '',
            'membership'               => $info['UserCode'],      //String	主卡卡号
            'total_inta_valid_points'  => (int)isset($info['total_inta_valid_points']) ? $info['total_inta_valid_points'] : 0, // 可兑换积分
            'pts_need_to_upgrade'      => (float)isset($info['pts_need_to_upgrade']) ? $info['pts_need_to_upgrade'] : 0.00, // 升级所需消费金额
            'pts_need_to_renew'        => (float)isset($info['pts_need_to_renew']) ? $info['pts_need_to_renew'] : 0.00, // 保级所需消费金额
            'msg'                      => '会员信息返回成功'
        ];

    }

    /**
     * 获取用户信息接口
     *
     * @param array $user
     *
     * @return array
     */
    public function getMemberInfo(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        $info = $service->getUserInfo();
        if (empty($info)) {
            return parent::reqMemberInfoNo($user);
        }

        return $this->handleCustomerInfo($info, $user);

    }

    /**
     * 通过unionid获取用户订单列表信息
     *
     * @param array $user
     *
     * @return array
     */
    public function getCustomerPurchase(array $user): array
    {

        $service = $this->userRequestService->setUser($user);

        $data = $service->getUserOrderListByUnionId($user['page'], $user['size']);

        /**
         * 参数名    类型    说明
         * SumSize    int    总数量
         * OrderItems    List
         * OrderCode    String    订单编号
         * TransactionTime    DateTime    订单创建时间
         * ProductName    string    商品名称
         * ProductSize    Int    商品数量
         * Money    decimal    下单金额
         * Source    String    来源
         */
        $records = [];
        if (!empty($data) && is_array($data['OrderItems']) && !empty($data['OrderItems'])) {
            foreach ($data['OrderItems'] as $v) {
                $val                  = [];
                $val['amount']        = '0';
                $skus                 = [];
                $val['store_name']    = $v['Source'] ?? '';
                $val['order_no']      = $v['OrderCode'] ?? '';
                $val['purchase_date'] = strtotime(handleDateTimeUTC($v['TransactionTime'])) * 1000;
                $val['amount']        += $v['Money'];
                $val['amount']        = (string)$val['amount'];

                $skus['store_name'] = $v['Source'] ?? '';
                $skus['order_no']   = $v['OrderCode'] ?? '';
                $skus['sku_name']   = $v['ProductName'] ?? '';
                $skus['quantity']   = (string)$v['ProductSize'];
                $skus['amount']     = (string)$v['Money'];

                $val['purchase_skus'][] = $skus;

                $records[] = $val;

            }
            $res['total'] = $data['SumSize'];

        }

        $res['pageNum']  = $user['page'];
        $res['pageSize'] = $user['size'];
        $res['records']  = $records;

        return $res;

    }

    /**
     * 根据customerInfo查询会员优惠卷记录
     *
     * @param array $user
     *
     * @return array
     */
    public function getCustomerCoupon(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        $data = $service->getUserCouponList();

        /**
         * 参数名    类型    说明
         * CouponType    int    优惠券类别1优惠券
         * CouponName    String    优惠券名称
         * CouponStartTime    datetime    有效开始时间
         * CouponEndTime    datetime    有效结束时间
         * Rule    String    规则描述
         */
        // 优惠券类别
        $couponTypeArr = ['1' => '优惠券'];
        $records       = [];
        if (!empty($data)) {
            foreach ($data as $v) {
                $val                = [];
                $val['type']        = $v['CouponType'];
                $val['type_str']    = $couponTypeArr[$v['CouponType']] ?? '';
                $val['name']        = $v['CouponName'];
                $val['startDate']   = handleDateTimeUTC($v['CouponStartTime']);
                $val['endDate']     = handleDateTimeUTC($v['CouponEndTime']);
                $val['description'] = $v['Rule'];

                $records[] = $val;

            }
        }
        $res['records'] = $records;

        return $res;

    }

    /**
     * 用户是否是会员接口
     *
     * @param array $user
     *
     * @return array
     */
    public function checkMember(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

        $data = $service->checkMember();

        if (isset($data['data']) && $data['data']) {
            $config = Config::getPathConfig(true);
        } else {
            $config = Config::getPathConfig(false);
        }
        $res['status'] = $data['data'];
        $res['config'] = $config;

        return $res;
    }

    /**
     * 用户参加相关活动，获取券
     *
     * @param array $user
     *
     * @return array
     */
    public function userJoin(array $user): array
    {
        $service = $this->userRequestService->setUser($user);

//        $campaignCode = $user['storeType'] == 1 ? Config::CampaignCodeAL : Config::CampaignCode;
        $campaignCode = Config::CampaignCode;
        $havedInfo    = $service->userJoin($campaignCode, $user['shopCode'], $user['clerkCode']);

        //0表示成功
        //4001表示已经参与领取过
        //4000活动过期、会员信息不存在
        // 新用户或未注册默认未领劵
        $res['haved']        = false;
        $res['welfare_info'] = [];
        if (!empty($havedInfo) && $havedInfo['errcode'] == 4001) {
            $res['haved'] = true;
        }
        return $res;
    }

}