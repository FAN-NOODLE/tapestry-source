<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Micros\Common;
use CrsSdk\MicroAuthority\Account;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;


class ChatGroupController extends BaseController
{
    /**
     * 获取社群的数据罗盘统计
     * @param Request $request
     * @return JsonResponse
     */
    public function getSummary(Request $request): JsonResponse
    {
        $params                  = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile']        = $this->mobile;
        $params['is_sm']         = 0;
        if ($this->role == Constant::LOGIN_ROLE_ENTERPRISE_ADMIN && empty($params['third_ids'])) {
            $params['is_sm'] = 1;
            $dataRange       = [];
            try {
                $dataRange = (new Account())->dataRangeList([
                    'admin_id' => $this->adminId,
                ]);
            } catch (Exception $e) {
            }
            if (!empty($dataRange)) {
                $params['third_ids'] = array_column($dataRange, "object_id");
            }
        }
        try {
            $result = Common::microScrm()->setApi("micro-scrm/v1/customer/data-compass/chat-group-summary")->post($params);
        } catch (Exception $e) {
            error(10002, $e->getMessage());
        }

        return success($result);
    }
}