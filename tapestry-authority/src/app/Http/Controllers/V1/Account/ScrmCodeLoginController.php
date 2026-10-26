<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\BrandAuth;
use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseService;
use App\Traits\CodeTrait;
use App\Common\Constant;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;
use Matrix\Exception;

class ScrmCodeLoginController extends BaseController
{
    use CodeTrait;

    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $client  = $request->post('client', 'console');
        $epIds   = $request->post('enterprise_ids', [$enterpriseId]);
        $account = $request->post('account');
        $code    = trim($request->post('password'));
        $ttl     = trim($request->post('ttl', EnterpriseAdminService::TOKEN_EXPIRE));

        if (empty($account)) {
            notice(700006);
        }
        if (empty($code)) {
            notice(700007);
        }
        if (empty($client)) {
            notice(700032);
        }
        if (empty($enterpriseId)) {
            notice(700020);
        }

        $account = str_replace('+86-', '', $account);
        $mobile  = '';
        if (!is_numeric($account) && strpos($account, Constant::TOKEN_SALT)) {
            if ((strlen($account) - strlen(Constant::TOKEN_SALT)) == 32)
                //如果不是数字且存在token_salt并且长度为32为则为管理后台短信验证码登录
                //兼容双元素认证登录功能，将手机号转换成加密串返回给前端
                //前端提交加密串发送验证码，避免绕过密码登录，该加密串存入redis五分钟，在这里验证
                if (!$trueAccount = $this->redis->get($account)) {
                    notice(700032);
                }
//            $this->redis->del($account);
            $mobile = $trueAccount;
        }

        if (!$mobile) {
            notice(700007);
        }

        //登录失败次数超过5次,锁60分钟
        $cacheLimitAccountTime = $this->getLockAccountCacheKey($mobile, $client) ?? 0;
        if ($cacheLimitAccountTime >= $this->failTimes) {
            Redis::connection()->expire($this->lockAccountCacheKey($mobile, $client), $this->lockTime);
            notice(700027);
        }

        list($mobile, $enterpriseId) = (new EnterpriseAdminService())->checkLoginByMobile($mobile, $enterpriseId, $code, $client);

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

        $token = (new EnterpriseAdminService())->generateToken($enterpriseInfo, $tmp, $ttl);

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
            'token'                 => $token,
            'enterprise_admin_list' => $adminList,
            'enterprise_id'         => $enterpriseId,
            'brand_code'            => $auth && isset($auth[$enterpriseId]) ? $auth[$enterpriseId] : ''
        ];
        $this->redis->del($account);
        return success($data);
    }

}
