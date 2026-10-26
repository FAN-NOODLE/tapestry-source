<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\Enterprise;
use App\Services\EnterpriseService;
use App\Http\Controllers\V1\BaseController;

class AccountSystemController extends BaseController
{
    /**
     * 获取创建编辑账号时企业可使用的系统
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $enterpriseId = $this->request->input('enterprise_id', 0);

        $detail = (new EnterpriseService())->enterpriseDetail($enterpriseId);

        if (empty($detail)) {
            notice(701005);
        }

        $data = [];
        if (!empty($detail['available_system'])) {
            $availableSystem = json_decode($detail['available_system']);
            array_multisort($availableSystem);

            foreach ($availableSystem as $v) {
                $sysList = Enterprise::SYSTEM[$v] ?? 0;
                if (!$sysList) continue;
                $data[] = [
                    'id'   => $v,
                    'name' => Enterprise::SYSTEM[$v]
                ];
            }
        }

        return success([
            'list' => $data
        ]);
    }
}