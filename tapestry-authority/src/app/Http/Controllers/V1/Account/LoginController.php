<?php

namespace App\Http\Controllers\V1\Account;

use App\Common\Constant;
use App\Models\BrandAuth;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Modules\BuiltInRole\BuiltInRoleActive\CommonRole;
use App\Modules\BuiltInRole\RoleStrategy;
use App\Services\BaseService;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseService;
use App\Traits\CodeTrait;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;
use Matrix\Exception;

class LoginController extends BaseController
{
    use CodeTrait;

    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $client       = $request->post('client', 'console');
        $epIds        = $request->post('enterprise_ids', []);
        $mobile       = $request->post('account');
        $name         = $request->post('name','');
        $code         = trim($request->post('password'));
        $loginType    = trim($request->post('login_type', EnterpriseAdminService::LOGIN_BY_CODE));
        $ttl          = trim($request->post('ttl', EnterpriseAdminService::TOKEN_EXPIRE));

        if (empty($mobile)) {
            notice(700006);
        }
        if (empty($code)) {
            notice(700006);
        }
        if (empty($client)) {
            notice(700032);
        }

        $mobile = str_replace('+86-', '', $mobile);
        //登录失败次数超过5次,锁24小时
        $cacheLimitAccountTime = $this->getLockAccountCacheKey($mobile, $client) ?? 0;
        if ($cacheLimitAccountTime >= $this->failTimes) {
            Redis::connection()->expire($this->lockAccountCacheKey($mobile, $client), $this->lockTime);
            notice(700027);
        }

        if ($loginType == EnterpriseAdminService::LOGIN_BY_CODE) {
            list($mobile, $enterpriseId) = (new EnterpriseAdminService())->checkLoginByMobile($mobile, $enterpriseId, $code, $client);
        } else {
            notice(700032);
        }

        // 该账号对应的企业列表
        $adminList = (new EnterpriseAdminService())->enterpriseAdminList($mobile, $client, '', $epIds);

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

        // 获取定制化权限标识
        $auth = (new BrandAuth())->getBrandAuth(
            $enterpriseId,
            '',
            [],
            [
                'enterprise_id',
                'brand_code'
            ]
        );
        $auth = array_column($auth, 'brand_code', 'enterprise_id');

        $token = (new EnterpriseAdminService())->generateToken($enterpriseInfo, $tmp, $ttl);

        $data = [
            'token'                 => $token,
            'enterprise_admin_list' => $adminList,
            'enterprise_id'         => $enterpriseId,
            'brand_code'            => $auth && isset($auth[$enterpriseId]) ? $auth[$enterpriseId] : ''
        ];

        return success($data);
    }

    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function miniProgramLogin(Request $request): \Illuminate\Http\JsonResponse
    {
        $enterpriseId = $request->post('enterprise_id');
        $mobile       = $request->post('mobile');
        $name       = $request->post('name');
        $adminInfo = (new EnterpriseAdminService())->checkLoginByWeChatWorkLogin($mobile, $enterpriseId, $name);
        if (empty($adminInfo)) {
            notice(700012);
        }
        $dataRangeList = (new EnterpriseDataRange())->accountDataRange($adminInfo['id']);
        return success(['info' => $adminInfo,'data_range_list' => $dataRangeList]);
    }

}
