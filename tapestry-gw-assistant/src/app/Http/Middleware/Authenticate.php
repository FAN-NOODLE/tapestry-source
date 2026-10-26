<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Contracts\Auth\Factory as Auth;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class Authenticate
{
    /**
     * The authentication guard factory instance.
     *
     * @var Auth
     */
    protected $auth;

    /**
     * Create a new middleware instance.
     *
     * @param Auth $auth
     * @return void
     */
    public function __construct(Auth $auth)
    {
        $this->auth = $auth;
    }

    /**
     * Handle an incoming request.
     *
     * @param Request $request
     * @param Closure $next
     * @param string|null $guard
     * @return mixed
     */
    public function handle(Request $request, Closure $next, string $guard = null)
    {
        if (config('app.env') == 'dev') {
            return $next($request);
        }
        if ($this->auth->guard($guard)->guest()) {
            error(652004, '登录状态失效', [], [], 401);
        }
        $token = $request->header('token');
        if (empty($token)) {
            error(652006);
        }
        // 解析Token
        $tokenData = json_decode(Crypt::decrypt($token), true);

        if ($tokenData['expiredTime'] < time()) {
            error(652004, '登录状态失效', [], [], 401);
        }

        if (!isset($tokenData['login_cache_key'])) {
            error(652004, '登录状态失效', [], [], 401);
        }

        // 获取 Cache里的loginId
        $cacheData = Cache::get($tokenData['login_cache_key']);
        Log::info('校验loginId：', [$cacheData['loginId']]);
        if (empty($cacheData)) {
            error(652004, '登录状态失效', [], [], 401);
        }

        // 添加企业登录护踢
        Log::info("cache_data", $cacheData);

        // 选择的企业信息丢失
        if (empty($cacheData['enterprise_id'])) {
            error(652008);
        }

        /** @var Response $response */
        return $next($request);
    }
}
