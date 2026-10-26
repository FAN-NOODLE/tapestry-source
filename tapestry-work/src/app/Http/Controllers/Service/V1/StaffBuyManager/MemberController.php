<?php


namespace App\Http\Controllers\Service\V1\StaffBuyManager;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Resources\MemberResource;
use App\Library\Uc\Member;
use App\Services\StaffManagerService;
use Illuminate\Http\Request;

class MemberController extends BaseController
{
    /**
     * @api /service/v1/staff-manager/member
     * @param Request $request
     * @param string qy_domain 企业唯一标识domain
     * @param string app_identifier 应用标识，默认contact
     * @param array  params
     * @return array
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', 'contact');
        $params        = $request->post('params', []);

        $result = (new StaffManagerService())->getStaffListFromUc($qyDomain, $appIdentifier, $params);
        // 格式化数据
        $data = (new StaffManagerService())->formatUcMember($result);

        return success(['data' => $data]);
    }
}
