<?php

namespace App\Http\Controllers\App\V1\Member;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\MicroAuthority\Account;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class MemberStoreController extends BaseController
{
    /**
     * 获取导购的部门列表
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request): JsonResponse
    {
        $admin_id = $this->loginAdminEnterprise['admin_id'];

        try {
            $dataRange = (new Account())->dataRangeList(['admin_id' => $admin_id]);
            $third_id  = array_column($dataRange, 'object_id');
            $params    = ['enterprise_id' => $this->enterpriseId, "third_ids" => $third_id, "mobile" => $this->mobile, "is_sm" => 1];
            if (!$this->isManager) {// 如果角色是导购的话，导购查询自己的门店列表
                $params['is_sm'] = 0;//is_sm  1为店长，即企业管理，0为店员，即企业导购
            }

            $result = Common::microScrm()->setApi("micro-scrm/v1/member/store-lists")->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }
        return success($result);
    }

    /**
     * coach 店长(SM)获取管理的导购(SA)
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function staffLists(Request $request): JsonResponse
    {
        $params   = $request->all();
        $store_id = $params['store_id'] ?? 0;
        $unionid  = $params['unionid'] ?? "";
        if ($this->isManager){
            $dataRange           = $this->dataRangeList;
            $third_id = array_column($dataRange, "object_id");
            $store_ids=$third_id;
        } else {
            $store_ids[] = $store_id;
        }
        $input  = ['store_ids' => $store_ids, 'enterprise_id' => $this->enterpriseId, 'mobile' => $this->mobile, "unionid" => $unionid];
        $result = Common::MicroScrm()->setApi("micro-scrm/v1/member/staff-lists")->post($input);
        return success($result);
    }
}