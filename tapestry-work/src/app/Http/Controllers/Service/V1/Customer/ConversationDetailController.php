<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Models\MemberUserRelation;
use App\Services\CustomerService;
use App\Services\EnterpriseService;
use Illuminate\Http\Request;
use App\Library\Uc\External;
use App\Http\Controllers\Service\V1\BaseController;

class ConversationDetailController extends BaseController
{


    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        if (!$enterpriseId) {
            error(600000);
        }
        $qyDomain = $request->input('qy_domain', '');
        if (empty($qyDomain)) {
            error(600003);
        }

        $memUid  = $request->input('mem_uid', '');
        $unionIds = $request->input('union_ids', []);

        //根据导购二维码id统计外部联系人数
        $memberUserList = MemberUserRelation::query()
            ->where('mem_uid', $memUid)
            ->whereIn('union_id', $unionIds)
            ->where('delete_type', '=', 0)
            ->orderBy('created_at', 'desc')
            ->get(['union_id','user_id', 'ec_userid'])
            ->toArray();

        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);
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

            $ecUserids                = array_column($memberUserList, 'ec_userid');

            $customerExchangeInfoList = CustomerService::getCustomerExchangeList($qyDomain,
                $identifier,
                $ecUserids);

            // 修改列表
            foreach ($memberUserList as &$value) {
                $value['ec_userid'] = $customerExchangeInfoList[$value['ec_userid']]['pltUserid'] ?? $value['ec_userid'];
            }
        }

        return success($memberUserList);
    }
}