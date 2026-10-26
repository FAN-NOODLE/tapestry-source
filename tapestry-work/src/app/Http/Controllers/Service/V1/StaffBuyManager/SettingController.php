<?php

namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class SettingController extends BaseController
{
    /**
     * 内购员工/部门设置返回
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id', 0);
        $type         = $request->post('type', 1);

        // 企业id不能为空
        if (empty($enterpriseId)) {
            notice(600000);
        }

        // 查询设置
        $data = (new StaffManagerService())->getStaffSetting($enterpriseId,$type);

        return success(['data' => $data]);
    }
}
