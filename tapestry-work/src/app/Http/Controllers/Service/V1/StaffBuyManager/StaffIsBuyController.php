<?php


namespace App\Http\Controllers\Service\V1\StaffBuyManager;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class StaffIsBuyController extends BaseController
{

    /**
     * 判断用户是否有内购权限
     * @param Request $request
     * @return array
     * @throws \Exception
     * @api /service/v1/staff-manager/is-right
     */
    public function index(Request $request)
    {
        $enterpriseId  = $request->input('enterprise_id', 0);
        $memUid        = $request->input('mem_uid', '');
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', 'contact');

        // 企业id不能为空
        if (empty($enterpriseId)) {
            notice(600000);
        }

        // 员工memUid不能为空
        if (empty($memUid)) {
            notice(720005);
        }

        if (empty($qyDomain)) {
            notice(600003);
        }

        // 获取和判断员工是不是属于允许内购得员工
        $res = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, $qyDomain, $appIdentifier);

        return success(['data' => ['is_right' => $res]]);
    }

}