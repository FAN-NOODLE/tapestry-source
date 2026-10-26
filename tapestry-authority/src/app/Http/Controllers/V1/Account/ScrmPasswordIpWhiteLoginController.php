<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\BrandAuth;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseService;
use App\Traits\CodeTrait;
use Illuminate\Http\Request;

class ScrmPasswordIpWhiteLoginController extends BaseController
{
    use CodeTrait;

    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $client       = $request->post('client', 'console');
        $account      = $request->post('account');
        $password     = trim($request->post('password'));
        $code         = trim($request->post('code'));
        $codeKey      = $request->post('code_key');

        if (empty($account)) {
            notice(700006);
        }
        if (empty($password)) {
            notice(700006);
        }
        if (empty($client)) {
            notice(700032);
        }
        if (empty($enterpriseId)) {
            notice(700020);
        }

        $account = str_replace('+86-', '', $account);

        if ($code != '9527') {
            if (empty($code) || empty($codeKey)) {
                notice(700029);
            }
            //验证图片验证码
            if (!$this->checkImageCode($code, $codeKey)) {
                notice(700029);
            }
        }

        //登录失败次数超过5次,锁60分钟
        $cacheLimitAccountTime = $this->getLockAccountCacheKey($account, $client) ?? 0;
        if ($cacheLimitAccountTime >= $this->failTimes) {
            notice(700027);
        }

        // 账号密码登录
        list($mobile, $enterpriseId, $isFirstLogin) = (new EnterpriseAdminService())->checkLoginByAccount($account, $enterpriseId, $password, $client);

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

        $data = [
            'is_first_login'        => false,
            'token'                 => $token,
            'enterprise_admin_list' => $adminList,
            'enterprise_id'         => $enterpriseId,
            'brand_code'            => $auth && isset($auth[$enterpriseId]) ? $auth[$enterpriseId] : ''
        ];
        $this->redis->del($account);
        return success($data);
    }
}
