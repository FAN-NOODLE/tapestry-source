<?php
/**
 * Created by PhpStorm.
 * User: BaiZhongYang
 * Date: 2020/7/17
 * Time: 15:40
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Models\CrsOrder;
use App\Services\EnterpriseService;
use App\Services\GuideAchievementSettingService;
use CrsSdk\Shop\Store;
use CrsSdk\Work\Member;
use App\Exceptions\Error;
use CrsSdk\Shop\CouponUser;
use App\Library\Uc\External;
use Illuminate\Http\Request;
use App\Services\CustomerService;
use App\Models\MemberUserRelation;
use Illuminate\Support\Facades\Log;
use App\Models\MemberOrderStatistics;
use App\Http\Controllers\Service\V1\BaseController;

class CustomerWisdomDetailController extends BaseController
{
    /**
     * index
     *
     * @desc 顾客详情接口
     *
     * @param Request $request
     *
     * @return array|mixed
     * @date 2020-07-17
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        if (!$enterpriseId) {
            return error(600000);
        }
        $domain = $request->input('qy_domain', '');
        if (empty($domain)) {
            return error(600003);
        }
        $hash = $request->input('qy_hash', '');
        if (empty($hash)) {
            return error(110002);
        }
        $memUid = $request->input('mem_uid', '');
        if (!$memUid) {
            return error(720005);
        }
        $unionId = $request->input('union_id', '');
        // 列表返回的外部联系人主键id
        $relationId = $request->input('relation_id', 0);
        if (empty($relationId)) {
            return error(600002);
        }
        $userId = $request->input('user_id', 0);

        $params = ['enterprise_id' => $enterpriseId];
        switchAppHost($params);

        $guideAchievementSetting = (new GuideAchievementSettingService())->detail($params);
        if (empty($guideAchievementSetting)) {
            $guideAchievementSetting =
                [
                    'enterprise_id'  => $enterpriseId,
                    'binding_type'   => 1,
                    'duration_time'  => 0,
                    'change_owner'   => 1,
                    'protected_time' => 0
                ];
        }

        // 查询顾客基本信息
        $customerInfo = $this->customerInfo($enterpriseId, $userId, $memUid, $relationId, $unionId, $domain, $hash);
        if (empty($customerInfo)) {
            return error(600012);
        }

        return success([
            'customer_info'             => $customerInfo,
            'guide_achievement_setting' => $guideAchievementSetting,
        ]);
    }


    /**
     * 获取客户信息
     *
     * @param $enterpriseId
     * @param $userId
     * @param $memUid
     * @param $relationId
     * @param $unionId
     * @param $qyDomain
     * @param $hash
     *
     * @return array
     * @throws \Exception
     */
    protected function customerInfo($enterpriseId, $userId, $memUid, $relationId, $unionId, $qyDomain, $hash)
    {
        Log::info('获取客户信息参数', [$enterpriseId, $memUid, $relationId, $unionId, $qyDomain, $hash]);
        $customerRelation = MemberUserRelation::getCustomerById($relationId);
        if (!$customerRelation){
            return  [
                'ec_name'         => "用户关系不存在",
                'ec_avatar'       => '',
                'ec_gender'       => 0,
                'ec_userid'       => '',
                'real_name'       => '',
                'third_member_no' => '',
                'tags'            => [],
                'used_coupon_num' => 0,
                'relation_at'     => 0,
                'user_id'         => '',
                'mem_uid'         => '',
                'latest_buy_time' => 0,
                'source_type'     => 0,
            ];
        }
        // 如果是好物集则不查询企业信息
    /*    if(Constant::HAO_WU_JI_ENTERPRISE==$enterpriseId){
            $customerInfo     = CustomerService::getCustomerListFromShop(0, [$unionId]);
        }else{
    */
            $customerInfo     = CustomerService::getCustomerListFromShop($enterpriseId, [$unionId]);
    //    }

        if (empty($customerInfo)) {
            $customerInfo = [
                'id' => 0,
            ];
        } else {
            $customerInfo = $customerInfo[0];
        }

        // 获取系统标签
        $tagInfo = CustomerService::getCustomerTagsFromUC($qyDomain, [$customerRelation['ec_id']]);

        // 获取优惠券核销总数
        if (!empty($userId)) {
            $countCondition = [
                'user_id'       => $userId,
                'enterprise_id' => $enterpriseId,
            ];
            $couponData     = (new CouponUser())->couponCountUseful($countCondition);
        }

        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);
        $identifierList = array_column($enterpriseIdentifierList, 'identifier');

        $exchangeInfo = [];
        if (in_array(Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT, $identifierList)) {
            $identifier = Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT;
            if (env('APP_ENV') == 'dev') {
                $identifier = 'crs_store_assistant_dev';
            } else {
                $list = array_column($enterpriseIdentifierList, null, 'identifier');
                if ($list[$identifier]['created_at'] > env('ICLK_APP_CREATED_AT')) {
                    $identifier = 'iclk_crs_store_assistant';
                }
            }
            $exchangeInfo = (new External())->exchangeInfo($qyDomain, $identifier,
                ['ecUserid' => $customerRelation['ec_userid']]);
        }

        // 查询最后支付时间
        $order = MemberOrderStatistics::query()
            ->where('union_id', $unionId)
            ->where('enterprise_id', $enterpriseId)
            ->first(['latest_buy_time']);

        $customerRelation['ec_userid'] = $exchangeInfo['pltUserid'] ?? $customerRelation['ec_userid'];

        $customerData = [
            'ec_name'         => $customerRelation['ec_name'] ?? '',
            'ec_avatar'       => $customerRelation['ec_avatar'] ?? '',
            'ec_gender'       => $customerRelation['ec_gender'] ?? 0,
            'ec_userid'       => $customerRelation['ec_userid'],
            'real_name'       => $customerInfo['real_name'] ?? '',
            'third_member_no' => $customerInfo['third_member_no'] ?? '',
            'tags'            => array_values($tagInfo[$customerRelation['ec_id']] ?? []),
            'used_coupon_num' => $couponData['used_count'] ?? 0,
            'relation_at'     => $customerRelation['relation_at'] ?? 0,
            'user_id'         => $customerInfo['id'],
            'mem_uid'         => $customerRelation['mem_uid'] ?? '',
            'latest_buy_time' => $order->latest_buy_time ?? 0,
            'source_type'     => $customerRelation['source_type'] ?? 0,
        ];

        return $customerData;
    }

}