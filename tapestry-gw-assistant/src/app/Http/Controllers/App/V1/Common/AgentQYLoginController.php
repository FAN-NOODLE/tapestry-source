<?php

namespace App\Http\Controllers\App\V1\Common;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Controller;
use Carbon\Carbon;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Work\StoreWorker;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;


class AgentQYLoginController extends Controller
{
    // scrm001
//    var $corpSecret = "REDACTED_FOR_SOURCE_SHARING";
    // 云店的secret
    var $corpSecret = "REDACTED_FOR_SOURCE_SHARING";


    var $corpId = "ww29202b3a45270ff3";

    public function getToken()
    {
        //验证请求参数
        $params = $this->request->all();

        if (empty($params['code'])) {
            return error(652012);
        }
        $enterprise_id = $params['enterprise_id'];
        if (empty($enterprise_id)) {
            return error(652012);
        }
        $accessToken = $this->getAccessToken();
        if (empty($accessToken)) {
            error(110002, '获取access_token失败');
        }
        $loginUrl = sprintf("https://qyapi.weixin.qq.com/cgi-bin/miniprogram/jscode2session?access_token=%s&js_code=%s&grant_type=authorization_code", $accessToken, $params['code']);
        $result   = curl_get_https($loginUrl);
        $wxResult = json_decode($result, true);
        $userId   = "";
        if (empty($wxResult)) {
            error(110002, '登录失败');
        } else {
            if ($wxResult['errcode'] > 0) {
                $errmsg  = $wxResult['errmsg'];
                $errcode = $wxResult['errcode'];
                Log::error("企微自建应用登录失败:{$errcode}-------{$errmsg}");
                error(110002, '获取企微登录失败');
            }
            $userId = $wxResult['userid'];
        }
        // uc返回mobile为空，则直接查询导购表获取
        $storeWorker = (new StoreWorker())->detail([
            'userid'        => $userId,
            'enterprise_id' => $enterprise_id,
            'use_scrm_db'   => true
        ]);
        if (empty($storeWorker)) {
            notice(652013);
        }
        Log::info('导购详情', $storeWorker);

//        $epDetail['data']['mobile'] = $storeWorker['mobile'];

        $mobile = $storeWorker['mobile'];
        // 登录
        $loginParams = [
            'client'         => Constant::APP_CLIENT_NAME,
            'account'        => sprintf(Constant::ACCOUNT_FORMAT, $mobile),
            'password'       => 'REDACTED_FOR_SOURCE_SHARING',
            'login_type'     => "third",
            'enterprise_ids' => $params['enterprise_id'],
            'enterprise_id'  => $params['enterprise_id'],
            'name'           => $storeWorker['name']
        ];
        try {
            $loginData = (new Account())->login($loginParams);
        } catch (Error $error) {
            Log::error("Account->login出错，错误信息：{$error->getMessage()}, 错误code：{$error->getCode()}", $loginParams);

            notice(652002);
        }

        if (empty($loginData['enterprise_admin_list']) || empty($loginData['enterprise_admin_list'][0]['guide_roles'])) {
            notice('login_manager_no_authority');
        }

        // 管理员存在的角色不能登录智慧导购
        $guideRoles = array_column($loginData['enterprise_admin_list'][0]['guide_roles'], 'role_mark');
        $hasPower   = array_intersect($guideRoles, Constant::LOGIN_ROLE_LIST);
        if (empty($hasPower)) {
            notice('login_manager_no_authority');
        }
        $loginCacheKey = "qy_agent_login_{$mobile}_" . $enterprise_id;
        if (Cache::has($loginCacheKey)) {
            Cache::forget($loginCacheKey);
        }
        // 生成登录ID
        $loginId = md5(uniqid() . microtime(true));
        // 登录ID 存入缓存，以实现同一端账号互顶
        $expiredTime = Constant::LOGIN_TOKEN_EXPIRED + time();
        Cache::put($loginCacheKey, [
            'loginId'       => $loginId,
            'expiredTime'   => $expiredTime,
            'enterprise_id' => $enterprise_id,
            'enterprise'    => array_column($loginData['enterprise_admin_list'], null, 'enterprise_id'),
            'mobile'        => $mobile
        ], Carbon::createFromTimestamp($expiredTime));

        Log::info("登录缓存内容", [
            'login_cache_info' => Cache::get($loginCacheKey)
        ]);

        // 生成门店端Token
        $tokenStr = Crypt::encrypt(json_encode([
            'login_cache_key' => $loginCacheKey,
            'method'          => 'qy_login',
            'mobile'          => $mobile,
            'loginId'         => $loginId,
            'expiredTime'     => $expiredTime,
            'user_type'       => 1
        ]));

        // 存储身份认证1小时50分钟
        Redis::connection()->setex($tokenStr, 6600, 'true');

        return success([
            'login_cache_key' => $loginCacheKey,
            'token'           => $tokenStr,
            'enterprise'      => $loginData['enterprise_admin_list'],
            'is_first_login'  => 0,
            'userid'          => $storeWorker['userid'],
            'brand_code'      => isset($loginData['brand_code']) ? $loginData['brand_code'] : ''
        ]);

    }

    public function getAccessToken()
    {
//        return "gRcpZDRceHgBG1jvye7NkAQb8DCy4T-i91qlEiK0heTu1_lTN2vrqc0y99sDaa_Jtqe2kMDKNc-0A8-I4EdKzPThfQ0sR_ircn5wQejm2-bw7cMb6dAZ93tkME_bU3SUekf78rBQIEnY3uD24-gO9kxV0uIZ032S4Fy43iY61TlHZ5OofuD3jXhKRDJ3hjJ_dfG0NLCYtHXqiTxNzncrTg";
        $Key         = "scrm001AccessToken";
        $accessToken = getOrSetData($Key);
        if (!$accessToken) {
            $qyWxloginUrl = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s";
            $url          = sprintf($qyWxloginUrl, $this->corpId, $this->corpSecret);
            $result       = curl_get_https($url);
            $wxResult     = json_decode($result, true);
            if (empty($wxResult)) {
                throw new Error(110002, '获取session_key 及 openID异常，企业微信内部错误');
            } else {
                if ($wxResult['errcode'] > 0) {
                    $errmsg  = $wxResult['errmsg'];
                    $errcode = $wxResult['errcode'];
                    throw new Error('110002', $errcode . '-' . $errmsg);
                } else {
                    getOrSetData($Key, $wxResult['access_token'], $wxResult['expires_in'] - 1000);
                    return $wxResult['access_token'];
                }
            }
        }
        return $accessToken;
    }
}