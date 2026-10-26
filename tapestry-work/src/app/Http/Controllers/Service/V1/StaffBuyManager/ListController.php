<?php


namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 参与内购的员工或部门列表返回
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', 'contact');

        // 企业id不能为空
        if (empty($enterpriseId)) {
            notice(600000);
        }

        if (empty($qyDomain)) {
            notice(600003);
        }

        // 查询列表
        // $list = (new StaffManagerService())->getStaffRight($enterpriseId, $qyDomain, $appIdentifier);
        $list = (new StaffManagerService())->getStaffRightList($enterpriseId);

        return success(['data' => $list]);
    }

}