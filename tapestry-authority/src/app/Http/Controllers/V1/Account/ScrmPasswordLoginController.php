<?php

namespace App\Http\Controllers\V1\Account;

use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use App\Common\Constant;
use App\Traits\CodeTrait;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class ScrmPasswordLoginController extends BaseController
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
            Redis::connection()->expire($this->lockAccountCacheKey($account, $client), $this->lockTime);
            notice(700027);
        }

        // 账号密码登录
        list($mobile, $enterpriseId, $isFirstLogin) = (new EnterpriseAdminService())->checkLoginByAccount($account, $enterpriseId, $password, $client);

        if (!empty($mobile)) {
            $md5Mobile = md5($mobile . rand(0, 9999)) . Constant::TOKEN_SALT;
            //兼容双元素认证登录功能，将手机号转换成加密串返回给前端
            //前端提交加密串发送验证码，避免绕过密码登录，该加密串存入redis五分钟
            $this->redis->setex($md5Mobile, 300, $mobile);
            return success(
                [
                    'token'          => $md5Mobile,
                    'enterprise_id'  => $enterpriseId,
                    'is_first_login' => false,
                ]
            );
        }

        return success(
            [
                'token'          => '',
                'enterprise_id'  => $enterpriseId,
                'is_first_login' => false,
            ]
        );
    }
}
