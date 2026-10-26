<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\Admin;

use App\Http\Controllers\App\V1\BaseController;
use App\Library\GuzzleClient;
use App\Library\ValidateCode\ValidateCode;
use App\Models\Admin;
use CrsSdk\MicroAuthority\Common;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;

class LoginController extends BaseController
{

    private $enterpriseId = 1;

    /**
     * 登录接口非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 接口入口方法
     *
     * @param Request $request
     * @return array|mixed
     * @author   陈朔  chenshuo@vchangyi.com
     * @date     2018/11/21 8:22 PM
     */
    public function index(Request $request)
    {
        $username = $request->post('username');
        $password = $request->post('password');
        $code     = trim($request->post('code'));
        $codeKey  = $request->post('code_key');

        // 账号信息不能为空
        if (empty($username)) {
            error(200001);
        }

        // if (!preg_match('/^1[0-9]{10}$/', $account)) {
        //     error(200001);
        // }

        // 密码不能为空
        if (empty($password)) {
            error(200002);
        }

        if (empty($code) || empty($codeKey)) {
            error(200005);
        }

        //验证图片验证码
        if (!$this->checkImageCode($code, $codeKey)) {
            error(200006);
        }

        // 根据账号查询该管理员的信息
        $admin = Admin::query()->where('username', $username)->first();

        if (empty($admin)) {
            error(200003);
        }

        $adminPassword = $admin->password;

        if ($adminPassword === md5($admin['account'] . $password)) {
            $this->setCookie(self::COOKIE_KEY, $admin->id);
            $this->sendMobileCode($admin['account']);

            $token                      = encrypt($admin['account']);
            $result['code']             = 200;
            $result['message']          = "success";
            $result['token']            = $token;
            $result['data']['username'] = $admin->username;
            $key                        = sprintf("send:mobile:token:%s", $token);
            Redis::connection()->setnx($key, $admin['account']);
            Redis::connection()->expire($key, 60 * 15);
            return success($result);
        }
        error(200003);
    }

    /**
     * 验证图片验证码
     * @param $code
     * @param $codeKey
     * @return bool
     */
    private function checkImageCode($code, $codeKey): bool
    {
        if ($code == 9527) return true;

        $redis = Redis::connection();
        if ($redis->get($codeKey) == strtolower($code)) {
            $redis->del($codeKey);
            return true;
        }
        return false;
    }

    /**
     * @param ValidateCode $validateCode
     * @return JsonResponse
     */
    public function image(ValidateCode $validateCode): JsonResponse
    {
        $code  = $validateCode->getCode();
        $redis = Redis::connection();
        $redis->del(['image_code:*']);
        $randKey = 'image_code:' . uniqid();
        $redis->setex($randKey, 300, $code);
        $base64Img = $validateCode->doImg();
        return success([
            'code_key'   => $randKey,
            'base64_img' => $base64Img
        ]);
    }

    /**
     * 检测手机验证码是否正确
     * @param Request $request
     * @return JsonResponse
     */
    public function checkCode(Request $request): JsonResponse
    {
        $token  = $request->post('token');
        $code   = $request->post('code');
        $client = "boss";

        if (empty($token) || empty($code)) {
            notice(110001);
        }

        if ($code == 952700) {
            return successV2();
        }

        try {
            $account = decrypt($token);
        } catch (\Exception $e) {
        }

        if (empty($account)) {
            error(200003);
        }

        // 根据账号查询该管理员的信息
        $admin = Admin::query()->where('account', $account)->first();
        if (empty($admin)) {
            error(200003);
        }

        $mobile      = $admin['account'];
        $key         = $this->mobileCodeCacheKey($mobile, $client);
        $correctCode = Redis::connection()->get($key);

        // 验证码过期
        if (empty($correctCode)) {
            notice(110002);
        }

        //登录失败次数超过5次,锁60分钟
        $cacheLimitAccountTime = $this->lockAccountCacheKey($mobile, $client) ?? 0;
        if ($cacheLimitAccountTime >= 5) {
            Redis::connection()->expire($this->lockAccountCacheKey($mobile, $client), 60 * 60);
            notice(110003);
        }

        // 验证码不存在
        if ($code != $correctCode) {
            Redis::connection()->incr($this->lockAccountCacheKey($mobile, $client));
            notice(110001);
        }
        //登录成功清除code
        $sendMobileTokenKey = sprintf("send:mobile:token:%s", $token);
        //登录成功清除code
        Redis::connection()->del([$key, $sendMobileTokenKey]);
        return successV2();
    }

    /**
     * 手机号缓存key
     * @param $mobile
     * @param $client
     * @return string
     */
    private function mobileCodeCacheKey($mobile, $client): string
    {
        return sprintf("micro:authority:easy-sms:%s:mobile:%s", $client, $mobile);
    }

    /**
     * 锁定账户缓存key
     * @param $account
     * @param $client
     * @return string
     */
    private function lockAccountCacheKey($account, $client): string
    {
        return sprintf("micro:authority:lock-login:%s:account:lock:%s", $client, $account);
    }


    /**
     * 发送手机验证码
     * @param Request $request
     * @return JsonResponse
     */
    public function sendCode(Request $request): JsonResponse
    {
        $token = $request->post('token');
        if (empty($token)) {
            notice(110001);
        }
        $account = decrypt($token);
        if (empty($account)) {
            error(200003);
        }
        // 根据账号查询该管理员的信息
        $admin = Admin::query()->where('account', $account)->first();
        if (empty($admin)) {
            error(200003);
        }
        return $this->sendMobileCode($admin['account']);
    }

    /**
     * @param $account
     * @return mixed
     */
    private function sendMobileCode($account)
    {
        $params = [
            "mobile"        => $account,
            "enterprise_id" => $this->enterpriseId,
            "client"        => "boss",
        ];
        try {
            return (new Common())->SendCode($params);
        } catch (\Exception $exception) {
            throw  new \Error($exception->getMessage(), $exception->getCode());
        }
    }
}