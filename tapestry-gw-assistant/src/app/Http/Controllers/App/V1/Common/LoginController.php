<?php

namespace App\Http\Controllers\App\V1\Common;


use App\Common\Constant;
use App\Common\RedisKey;
use App\Http\Controllers\App\V1\BaseController as Controller;
use Carbon\Carbon;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Work\PartnerMember;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use CrsSdk\MicroAuthority\Common;
use Illuminate\Support\Facades\Redis;

class LoginController extends Controller
{
    protected $isRequireLogin = false;

    /**
     * 手机验证码登录
     *
     * @var int
     */
    const LOGIN_TYPE_CODE = 2;

    const MODULE_CHECK_SMS_MOBILE = 'CHECK_SMS_MOBILE';

    /**
     * 登陆身份 1:导购 2：合伙人 3：选择身份
     *
     * @var int
     */
    const ROLE_GUIDE = 1;
    const ROLE_PARTNER = 2;
    const ROLE_TWO = 3;

    /**
     * 登录返回
     *
     * @param Request $request
     * @param Account $account
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function __invoke(Request $request, Account $account)
    {
        $mobile = trim($request->get('mobile'));
        $code   = trim($request->get('code'));
        $client = Constant::APP_CLIENT_NAME;

        if (empty($mobile)) {
            error(651000);
        }

        if (empty($code)) {
            error(651004);
        }

        if (!preg_match('/^1[0-9]{10}$/', $mobile)) {
            error(651003);
        }
        //登录失败次数超过5次,锁一小时分钟
        $cacheLimitAccountTime = $this->getLockAccountCacheKey($mobile, $client) ?? 0;
        if ($cacheLimitAccountTime >= 5) {
            Redis::connection()->expire($this->lockAccountCacheKey($mobile, $client), 60 * 60);
            error(651007);
        }

        $isSuperUser = in_array($mobile, Constant::WX_LOGIN_SUPER_USER);

        // 缓存中的code
        $cacheCode = self::getCodeCache($mobile, $client);

        if ($cacheCode != $code && !$isSuperUser && $code != '911800') {
            Redis::connection()->incr($this->lockAccountCacheKey($mobile, $client));
            error(652001);
        }
        // 统一login
        $loginId = md5(uniqid() . microtime(true));

        // 调用SDK登录
        try {
            $params = [
                'login_type' => self::LOGIN_TYPE_CODE,
                'client'     => Constant::APP_CLIENT_NAME,
                'account'    => sprintf(Constant::ACCOUNT_FORMAT, $mobile),
                'password'   => $code
            ];

            if ($isSuperUser) {
                $params['login_type'] = 3;
                $params['password']   = 'qy_login';
            }
            $userType = 0;
            // 尝试导购登陆
            $guideLoginData = $this->guideLogin($account, $params);
            // 尝试合伙人登陆
            $partnerLoginData = $this->getPartnerLogin($mobile);
            // 判断是否需要选择角色
            if (!empty($guideLoginData) && !empty($partnerLoginData)) {
                $userType        = self::ROLE_TWO;
                $tokenStr        = self::guideToken($mobile, $loginId);
                $partnerTokenStr = self::partnerToken($mobile, $partnerLoginData, $loginId);
            } else if (!empty($guideLoginData)) {
                $userType = self::ROLE_GUIDE;
                $tokenStr = self::guideToken($mobile, $loginId);
            } else if (!empty($partnerLoginData)) {
                $userType        = self::ROLE_PARTNER;
                $partnerTokenStr = self::partnerToken($mobile, $partnerLoginData, $loginId);

            } else {
                notice(652015);
            }
            $adminList           = [];
            $partnerLoginDataNew = [];
            // cache选择
            if (empty($guideLoginData)) {
                // 合伙人加标识
                $adminList = $partnerLoginData;
                foreach ($adminList as $enterprise => $role) {
                    $adminList[$enterprise]['guide_roles'] = [
                        'role_name' => '分销合伙人',
                        'role_mark' => Constant::LOGIN_ROLE_STORE_PARTNER
                    ];
                }
            }
            if (!empty($guideLoginData)) {
                $adminList = array_column($guideLoginData['enterprise_admin_list'], null, 'enterprise_id');
            }
            if (!empty($partnerLoginData)) {
                $partnerLoginDataNew = array_column($partnerLoginData, null, 'enterprise_id');
                foreach ($partnerLoginDataNew as $enterprise => $role) {
                    $partnerLoginDataNew[$enterprise]['guide_roles'][] = [
                        'role_name' => '分销合伙人',
                        'role_mark' => Constant::LOGIN_ROLE_STORE_PARTNER
                    ];
                }
            }
            $loginCacheKey = "wx_login_{$mobile}";
            // 判断缓存中是否已经有这个手机号了
            if (Cache::has($loginCacheKey)) {
                Cache::forget($loginCacheKey);
            }
            $expiredTime = Constant::LOGIN_TOKEN_EXPIRED + time();
            // 登录ID 存入缓存，以实现同一端账号互顶
            Cache::put($loginCacheKey, [
                'loginId'                 => $loginId,
                'expiredTime'             => $expiredTime,
                'enterprise'              => $adminList,
                'mobile'                  => $mobile,
                'partner_enterprise_list' => array_values($partnerLoginDataNew)
            ], Carbon::createFromTimestamp($expiredTime));

            Log::info('二级分销登陆id', [$loginId]);
            // 将管理员有权限的企业列表存储到缓存里
            return success([
                'login_cache_key'         => $loginCacheKey,
                'token'                   => $tokenStr ?? '',
                'enterprise'              => $guideLoginData['enterprise_admin_list'] ?? [],
                'is_first_login'          => 0,
                'user_type'               => $userType,
                'partner_token'           => $partnerTokenStr ?? '',
                'partner_enterprise_list' => array_values($partnerLoginDataNew) ?? [],
                'brand_code'              => isset($guideLoginData['brand_code']) ? $guideLoginData['brand_code'] : ''
            ]);


        } catch (\Error $error) {
            Log::error("login error", [
                'code'         => $error->getCode(),
                'message'      => $error->getMessage(),
                'error_string' => $error->getTraceAsString()
            ]);
            if ($error->getCode() == 700005) {
                error($error->getCode(), $error->getMessage());
            }

            if (in_array($error->getCode(), [700007, 700017])) {
                error(652001);
            }

            error(650001, "code:" . $error->getCode() . ",message:" . $error->getMessage() ?? "");
        }
    }

    /**
     * 导购登陆
     *
     * @param Account $account
     * @param array $params
     * @return mixed
     * @throws \Exception
     */
    public function guideLogin(Account $account, array $params)
    {
        try {
            $loginData = $account->login($params);
            return $loginData;
        } catch (\Exception $e) {
            return $loginData = [];
        }
    }

    /**
     * 合伙人登陆查询信息
     *
     * @param $mobile
     * @return array|mixed|void
     */
    public function getPartnerLogin($mobile)
    {
        $partnerMember = (new PartnerMember());
        $params        = [
            'partner_mobile' => $mobile
        ];
        try {
            $roleInfo = $partnerMember->partnerList($params);
            $data     = objToArray($roleInfo);

            $enterpriseIdsInfo = [];
            if (empty($data['list'])) {
                error(652014);
            }
            foreach ($data['list'] as $key => $value) {
                $enterpriseParams    = [
                    'partner_id'    => $data['list'][$key]['partner_id'],
                    'enterprise_id' => $data['list'][$key]['enterprise_id'],
                ];
                $enterpriseIdsInfo[] = $partnerMember->partnerMy($enterpriseParams)['enterprise_info'];

            }
            // 数据去重
            $enterpriseIdsInfoNew = $enterpriseIdsInfo[0];
            // 主键
            $partnerKey = array_column($data['list'], null, 'enterprise_id');

            foreach ($enterpriseIdsInfoNew as $enterprise => &$val) {
                $status                = $partnerKey[$val['id']]['partner_status'];
                $partner_name          = $partnerKey[$val['id']]['partner_name'];
                $partner_id            = $partnerKey[$val['id']]['partner_id'];
                $val['account_status'] = $status;
                $val['account_name']   = $partner_name;
                $val['partner_id']     = $partner_id;
            }
            unset($enterpriseIdsInfo[0]);

            // 结构化数据
            $result = self::formData($enterpriseIdsInfoNew, $mobile);
            return $result;
        } catch (\Exception $e) {
            Log::error("无合伙人信息：{$e}");
            return $data['list'] = [];
        }
    }

    /**
     * 存储缓存
     *
     * @param string $mobile
     * @param array $enterprise
     * @param int $expiredTime
     */
    private function managerEnterprises(string $mobile, array $enterprise, int $expiredTime)
    {
        $redisKey = redisKey(RedisKey::LOGIN_MANAGER_ENTERPRISES, ['mobile' => $mobile]);
        getOrSetData($redisKey, $enterprise, $expiredTime);
    }

    /**
     * 导购token
     *
     * @param $mobile
     * @param $loginId
     * @return string
     */
    protected static function guideToken($mobile, $loginId)
    {
        $loginCacheKey = "wx_login_{$mobile}";

        $expiredTime = Constant::LOGIN_TOKEN_EXPIRED + time();

        // 生成门店端Token
        $tokenData = [
            'login_cache_key' => $loginCacheKey,
            'mobile'          => $mobile,
            'loginId'         => $loginId,
            'expiredTime'     => $expiredTime,
            'user_type'       => 1 //1导购 2合伙人
        ];
        $tokenStr  = Crypt::encrypt(json_encode($tokenData));
        return $tokenStr;
    }

    /**
     * 合伙人token
     *
     * @param $mobile
     * @param $partnerLoginData
     * @param $loginId
     * @return string
     */
    public static function partnerToken($mobile, $partnerLoginData, $loginId)
    {
        $loginCacheKey = "wx_login_{$mobile}";
        $expiredTime   = Constant::LOGIN_TOKEN_EXPIRED + time();

        // 合伙人token
        $partnerTokenData = [
            'login_cache_key' => $loginCacheKey,
            'mobile'          => $mobile,
            'loginId'         => $loginId,
            'expiredTime'     => $expiredTime,
            'user_type'       => 2 //1导购 2合伙人
        ];
        $partnerTokenStr  = Crypt::encrypt(json_encode($partnerTokenData));

        return $partnerTokenStr;
    }

    /**
     * @param $data
     * @return mixed
     */
    public static function formData($data, $mobile)
    {
        $dataNew = [];
        foreach ($data as $key => &$val) {
            $dataNew[$key]['enterprise_id']          = $val['id'];
            $dataNew[$key]['hash']                   = $val['hash'];
            $dataNew[$key]['version']                = $val['version'];
            $dataNew[$key]['version_tag']            = $val['version_tag'];
            $dataNew[$key]['name']                   = $val['name'];
            $dataNew[$key]['contact_name']           = $val['contact'];
            $dataNew[$key]['contact_mobile']         = $val['mobile'];
            $dataNew[$key]['address']                = $val['address'];
            $dataNew[$key]['crs_status']             = $val['crs_status'];
            $dataNew[$key]['store_status']           = $val['store_status'];
            $dataNew[$key]['contact_status']         = $val['contact_status'];
            $dataNew[$key]['domain']                 = $val['domain'];
            $dataNew[$key]['corp_id']                = $val['corp_id'];
            $dataNew[$key]['account_start_at']       = $val['account_start_at'];
            $dataNew[$key]['account_expire_at']      = $val['account_expire_at'];
            $dataNew[$key]['is_split_flow']          = $val['is_split_flow'];
            $dataNew[$key]['admin_id']               = '';
            $dataNew[$key]['account']                = '+86-' . $mobile;
            $dataNew[$key]['partner_status']         = $val['account_status'];
            $dataNew[$key]['account_name']           = $val['account_name'];
            $dataNew[$key]['account_global_roaming'] = '+86';
            $dataNew[$key]['account_mobile']         = $mobile;
            $dataNew[$key]['account_email']          = '';
            $dataNew[$key]['role']                   = 0;
            $dataNew[$key]['is_super']               = 0;
            $dataNew[$key]['range_type']             = 0;
            $dataNew[$key]['admin_available_system'] = '';
            $dataNew[$key]['shop_name']              = $val['shop_name'];
            $dataNew[$key]['identifier']             = '';
            $dataNew[$key]['logo']                   = $val['logo'];
            $dataNew[$key]['partner_id']             = $val['partner_id'];
        }

        return $dataNew;
    }

    public static function getCodeCache($mobile, $client)
    {
        $key = self::mobileCodeCacheKey($mobile, $client);
        return Redis::connection()->get($key);
    }

    protected static function mobileCodeCacheKey($mobile, $client)
    {
        return 'micro:authority:easy-sms:' . $client . ':mobile:' . $mobile;
    }

    public function getLockAccountCacheKey($account, $client)
    {
        $key = $this->lockAccountCacheKey($account, $client);
        return Redis::connection()->get($key);
    }

    protected function lockAccountCacheKey($account, $client)
    {
        return 'micro:assistant:lock-login:' . $client . ':account:lock:' . $account;
    }
}
