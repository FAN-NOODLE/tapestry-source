<?php

namespace App\Http\Controllers\App\V1\Common;


use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use Carbon\Carbon;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class CheckLoginController extends Controller
{
    protected $isRequireLogin = false;

    /**
     * 验证token 是否有效
     *
     * @param Request $request
     * @param Account $account
     *
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function __invoke(Request $request, Account $account)
    {
        // 要检验的Token
        $token = $request->get('token');
        if (empty($token)) {
            error(652007);
        }
        // 解析Token
        $tokenData = json_decode(Crypt::decrypt($token), true);
        if ($tokenData['expiredTime'] < time()) {
            error(652011);
        }
        // 获取 Cache里的loginId
        $cacheData = Cache::get($tokenData['login_cache_key']);
        Log::info('cache为空：', ['data' => $cacheData]);

        if ($cacheData && $tokenData['loginId'] !== $cacheData['loginId']) {
            error(652011);
        }
        // 返回此登录管理员的有权限的列表
        $enterprises['enterprise'] = [];
        try {
            if ($tokenData['user_type'] == 2) {

                $this->loginCache($cacheData['expiredTime'], $tokenData['login_cache_key'], $cacheData);
                Log::info('partner_enterprise_list',[$cacheData]);
                $result = [
                    'user_type'               => $tokenData['user_type'],
                    'partner_enterprise_list' => array_values($cacheData['partner_enterprise_list']),
                    'login_cache_key'         => $tokenData['login_cache_key'],
                ];
                Log::info('partner_enterprise_list_18',[array_values($cacheData['partner_enterprise_list'])]);
                return success($result);
            }

            $accountEnterpriseList = $account->accountEnterpriseList([
                'account' => sprintf(Constant::ACCOUNT_FORMAT, $cacheData['mobile']),
                'client'  => Constant::APP_CLIENT_NAME
            ]);
            // 过滤未注册商城小程序企业
            foreach ($accountEnterpriseList['enterprise_admin_list'] as $k => $v) {
                if (!$v['shop_name']) {
                    unset($accountEnterpriseList['enterprise_admin_list'][$k]);
                }
            }

            $enterprises['enterprise']      = $accountEnterpriseList['enterprise_admin_list'];
            $enterprises['login_cache_key'] = $tokenData['login_cache_key'];
            $enterprises['user_type']       = $tokenData['user_type'];
            $this->loginCache($cacheData['expiredTime'], $tokenData['login_cache_key'], $cacheData);
            Log::info("check-login-cache-data", $cacheData ?? []);
            return success($enterprises);
        } catch (\Exception $e) {

            return success($enterprises);
        }

    }

    /**
     * 缓存
     *
     * @param $expiredTime
     * @param $loginCacheKey
     * @param $cacheData
     */
    public function loginCache($expiredTime, $loginCacheKey, $cacheData)
    {
        $expiredTime = Carbon::createFromTimestamp($expiredTime);
        Cache::forget($loginCacheKey);
        Cache::add($loginCacheKey, $cacheData, $expiredTime);

    }
}
