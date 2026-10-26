<?php

namespace App\Http\Controllers\App\V1\Common;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController as Controller;
use App\Micros\Common;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Models\MobileWhite;
use App\Models\UserLoginLog;
use App\Services\MobileWhiteService;
use Carbon\Carbon;
use CrsSdk\Boss\BaseClient;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class QYLoginController extends Controller
{
    const LOGIN_TYPE_QY = 3;
    /**
     * 标准标识
     */
    const STANDARD_IDENTIFIER = 1;

    protected $isRequireLogin = false;

    public function __invoke()
    {
        //验证请求参数
        $params = $this->request->all();

        if (empty($params['code'])) {
            error(652012);
        }
        $enterprise = (new Enterprise())->detail([
            'corp_id' => $params['corp_id']
        ]);
        if (empty($enterprise['data'])) {
            Log::info('获取企业相关信息失败', [
                'corp_id' => $params['corp_id']
            ]);
            notice('enterprise_not_fund');
        }
        // 获取mem_uid
        $authInfo = GoMicroScrm::code2session([
            'js_code'       => $params['code'],
            'enterprise_id' => $enterprise['data']['id'],
            'domain'        => $enterprise['data']['domain']
        ]);
        if (empty($authInfo) || empty($authInfo['mobile'])) {
            notice('login_manager_no_authority');
        }
        Log::info('micro-scrm登录接口返回信息', $authInfo);
        // 登录
        $loginParams = [
            'mobile'        => $authInfo['mobile'],
            'name'          => $authInfo['name'],
            'enterprise_id' => $enterprise['data']['id'],
        ];
        $adminInfo   = Common::Authority()->setApi('micro-authority/v1/account/mini-program-login')->post($loginParams);
        if (empty($adminInfo) || empty($adminInfo['info'])) {
            notice(652003, [], [], 500);
        }

        // 登录记录
        $where = [
            'enterprise_id'  => $enterprise['data']['id'],
            'mem_uid'        => '',
            'mem_userid'     => $authInfo['userid'],
            'login_day_time' => strtotime(date('Y-m-d') . '00:00:00'),
        ];
        UserLoginLog::query()->firstOrCreate($where, ['id' => idGenerate()]);

        $res=MobileWhiteService::WhiteList($enterprise['data']['id'],$authInfo['mobile'],MobileWhite::NEW_CUSTOMER_LIST);
        $customerListWhite=$res?1:0;
        $loginCacheKey = "qy_login_{$authInfo['mobile']}_" . $enterprise['data']['id'];
        if (Cache::has($loginCacheKey)) {
            Cache::forget($loginCacheKey);
        }
        // 生成登录ID
        $loginId = md5(uniqid() . microtime(true));
        // 登录ID 存入缓存，以实现同一端账号互顶
        $expiredTime = Constant::LOGIN_TOKEN_EXPIRED + time();
        Cache::put($loginCacheKey, [
            'loginId'          => $loginId,
            'expiredTime'      => $expiredTime,
            'enterprise'       => [$adminInfo['info']],
            'choiceEnterprise' => array_merge($adminInfo['info'], ['admin_id' => $adminInfo['info']['id'], 'domain' => $enterprise['data']['domain']]),
            'data_range_list'  => $adminInfo['data_range_list'],
            'enterprise_id'    => $enterprise['data']['id'],
            'role'             => '',
            'member_auth_info' => $authInfo,
            'mobile'           => $authInfo['mobile'],
            'is_manager'       => $authInfo['is_manager'],
            "is_beta_user"     => $customerListWhite,
        ], Carbon::createFromTimestamp($expiredTime));
        Log::info("登录缓存内容", [
            'login_cache_info' => Cache::get($loginCacheKey)
        ]);

        // 生成门店端Token
        $tokenStr = Crypt::encrypt(json_encode([
            'login_cache_key' => $loginCacheKey,
            'method'          => 'qy_login',
            'mobile'          => $authInfo['mobile'],
            'loginId'         => $loginId,
            'expiredTime'     => $expiredTime,
            'user_type'       => self::USERTYPE_WORKER,
            "is_beta_user"    => $customerListWhite,
        ]));

        // 神策数据上报
        Common::MicroScrm()->setApi('micro-scrm/v1/member/report-data')->post([
            'enterprise_id'    => $enterprise['data']['id'],
            'userid'          => $authInfo['userid'],
        ]);

        return success([
            'login_cache_key'           => $loginCacheKey,
            'token'                     => $tokenStr,
            'enterprise'                => [$adminInfo['info']],
            'userid'                    => $authInfo['userid'],
            'store_code'                => $authInfo['store_code'],
            'store_name'                => $authInfo['store_name'],
            'is_manager'                => $authInfo['is_manager'],
            'store_city'                => $authInfo['store_city'],
            'store_phone'               => $authInfo['store_phone'],
            'store_inventory_name'      => $authInfo['store_inventory_name'],
            "is_beta_user"              => $customerListWhite,
        ]);

    }
}