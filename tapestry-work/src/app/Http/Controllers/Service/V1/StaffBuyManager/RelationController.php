<?php

namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class RelationController extends BaseController
{
    /**
     * 查询内购员工对应部门
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId  = $request->post('enterprise_id', 0);
        $memUid        = $request->post('mem_uid', '');
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', 'contact');

        // 企业id不能为空
        if (empty($enterpriseId)) {
            notice(600000);
        }

        if (empty($qyDomain)) {
            notice(600003);
        }

        if (empty($memUid)) {
            notice(600006);
        }

        // 查询对应部门
        list($dpId, $dpName) = (new StaffManagerService())->getStaffRightDepartmentRelation($enterpriseId, $memUid, $qyDomain, $appIdentifier);

        return success(['data' => ['dp_id' => $dpId, 'dp_name' => $dpName]]);
    }
}
