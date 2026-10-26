<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/17
 * Time: 10:26 上午
 */

namespace App\Http\Controllers\Service\V1\DistributionPartner;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsOrder;
use App\Models\DistributionOrderSettlement;
use App\Models\MemberUserRelation;
use App\Models\PartnerEnterpriseMap;
use App\Models\ShoppingGuideCashWithdrawal;
use App\Services\CustomerService;
use App\Services\EnterpriseService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class DisableDetailController extends BaseController
{
    /**
     * 订单结算状态
     *
     */
    const ORDER_STATUS_FINISH = 100;


    public function index(Request $request)
    {
        $params = $request->only([
            'enterprise_id',
            'partner_id',
            'mobile',
        ]);

        // 验证规则
        $rules = [
            'enterprise_id' => 'integer',
            'partner_id'    => 'integer',
            'mobile'        => 'integer',
        ];

        // 提示语
        $messages = [
            'partner_id.integer' => '合伙人信息需为int类型',
            'mobile.integer'     => '手机号需为int类型',
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        // 字段筛选
        $selectFields = [
            'partner_member.id', 'partner_member.partner_name', 'partner_member.avatar', 'partner_enterprise_map.user_id', 'partner_enterprise_map.created_at',
            'partner_enterprise_map.partner_status', 'partner_enterprise_map.forbidden_reason',
            'partner_enterprise_map.worker_uid', 'partner_enterprise_map.worker_name', 'partner_enterprise_map.mobile_show',
            'partner_enterprise_map.mobile_sensitive'
        ];
        // 详情
        $sql = PartnerEnterpriseMap::query()
            ->leftJoin('partner_member', 'partner_enterprise_map.partner_id', '=', 'partner_member.id')
            ->when(!empty($params['partner_mobile']), function ($query) use ($params) {
                return $query->where('partner_member.partner_mobile', $params['partner_mobile']);
            })
            ->where('partner_enterprise_map.enterprise_id', '=', $params['enterprise_id'])
            ->where('partner_member.id', '=', $params['partner_id']);

        $detail = $sql->first($selectFields);
        if (!empty($detail)) {
            $detail = $detail->toArray();
        } else {
            throw new Error(700000, '合伙人不存在');
        }
        Log::info('合伙人详情', $detail);

        // 查询合伙人对应的外部联系人userid
        $detail['ec_userid'] = '';
        if ($detail['user_id']) {
            switchAppHost(['enterprise_id' => $params['enterprise_id']]);
            $userInfo = (new User())->detail([
                'enterprise_id' => $params['enterprise_id'],
                'id' => $detail['user_id']
            ])['data'];
            Log::info('userinfo', $userInfo);
            if (!empty($userInfo)) {
                $ecUserId = MemberUserRelation::query()
                    ->where('enterprise_id', '=', $params['enterprise_id'])
                    ->where('mem_uid', '=', $detail['worker_uid'])
                    ->where('union_id', '=', $userInfo['union_id'])
                    ->value('ec_userid');
                if (!empty($ecUserId)) {
                    $detail['ec_userid'] = $this->convertEcUserId($ecUserId, $params);
                }

            }
        }

        // 分销订单数
        $distributionOrders = DistributionOrderSettlement::query()
            ->where('enterprise_id', '=', $params['enterprise_id'])
            ->where('partner_id', '=', $params['partner_id'])
            ->where('order_status', '>=', 30) // 已支付
            ->count();

        // 佣金结算笔数
        $numberOfCommissionSettlement = DistributionOrderSettlement::query()
            ->where('enterprise_id', '=', $params['enterprise_id'])
            ->where('partner_id', '=', $params['partner_id'])
            ->where('settlement_status', '=', DistributionOrderSettlement::ADOPT_SETTLE)
            ->count();

        // 佣金结算金额
        $commissionSettlementAmount = ShoppingGuideCashWithdrawal::query()
            ->where('enterprise_id', '=', $params['enterprise_id'])
            ->where('worker_id', '=', $params['partner_id'])
            ->select(['cash_withdrawal','frozen_amount', 'available_cash_amount'])
            ->first();
        if(!empty($commissionSettlementAmount)){
            $commissionSettlementAmount = ($commissionSettlementAmount->cash_withdrawal
                    + $commissionSettlementAmount->frozen_amount
                    + $commissionSettlementAmount->available_cash_amount)/100;
        }

        $detail['orders']                       = $distributionOrders;
        $detail['commission_settlement']        = $numberOfCommissionSettlement;
        $detail['commission_settlement_amount'] = $commissionSettlementAmount ?? 0 ;

        unset($distributionOrders);
        unset($numberOfCommissionSettlement);
        unset($commissionSettlementAmount);

        return success($detail);
    }

    /**
     * 转换ecUserId
     * @param $ecUserId
     * @param $params
     * @return mixed|string
     * @throws \Exception
     */
    private function convertEcUserId($ecUserId, $params)
    {
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']])['data'];
        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($params['enterprise_id']);
        $identifierList = array_column($enterpriseIdentifierList, 'identifier');

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
            $customerExchangeInfoList = CustomerService::getCustomerExchangeList($enterpriseInfo['domain'],
                $identifier,
                [$ecUserId]);
            Log::info('转换后ec_userid', $customerExchangeInfoList);
            $ecUserId = $customerExchangeInfoList[$ecUserId]['pltUserid'] ?? $ecUserId;
        }

        return $ecUserId;
    }
}