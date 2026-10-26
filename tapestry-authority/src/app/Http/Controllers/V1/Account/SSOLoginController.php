<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\BrandAuth;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseService;
use App\Traits\CodeTrait;
use Illuminate\Http\Request;

class SSOLoginController extends BaseController
{
    use CodeTrait;

    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $client       = $request->post('client', 'console');
        $mobile      = $request->post('mobile');


        if (empty($mobile)) {
            notice(700006);
        }
        if (empty($client)) {
            notice(700032);
        }
        if (empty($enterpriseId)) {
            notice(700020);
        }

        // 该账号对应的企业列表
        $adminList = (new EnterpriseAdminService())->enterpriseAdminList($mobile, $client, '', [$enterpriseId]);

        if (empty($adminList)) {
            notice(700012);
        }
        $enterpriseTotal = count($adminList);
        $enterpriseInfo  = $adminList[0];

        $tmp = true;
        if ($enterpriseTotal == 1) {
            $tmp     = false;
            $adminId = $enterpriseInfo['admin_id'];
            $version = $enterpriseInfo['version'];

            $adminInfo = $adminList[0];
            if ($version != EnterpriseService::VERSION_SHOP) {
                $adminInfo['enterprise_role_list'] = (new EnterpriseAdminService())->enterpriseAdminRoleList($adminId);
            }

            if ($adminInfo['account_status'] == EnterpriseAdminService::ACCOUNT_STATUS_BAN) {
                notice(700015);
            }

            $adminList[0] = $adminInfo;
        }

        $token = (new EnterpriseAdminService())->generateToken($enterpriseInfo, $tmp, EnterpriseAdminService::TOKEN_EXPIRE);


        $data = [
            'is_first_login'        => false,
            'token'                 => $token,
            'enterprise_admin_list' => $adminList,
            'enterprise_id'         => $enterpriseId,
        ];
        return success($data);
    }
}
