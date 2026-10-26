<?php
/**
 * Created By PhpStorm
 * @Date: 2021/6/24
 * @Time: 19:13
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: CheckAuth.php
 */

namespace App\Http\Middleware;

use App\Exceptions\Error;
use Closure;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;
use CrsSdk\MicroAuthority\BrandAuth;
use Illuminate\Support\Facades\Redis;

class checkAuth
{

    /**
     * 校验权限
     * @param \Illuminate\Http\Request $request
     * @param Closure $next
     * @return mixed|void
     * @throws \Exception
     * @Date: 2021/6/30
     * @Time: 18:42
     */
    public function handle($request, Closure $next)
    {
        $param = [
            'enterprise_id' => $this->getLoginEnterpriseId($request),
            'business'      => array_reverse(explode('/', $request->path()))[1],
        ];

        $auth            = (new BrandAuth())->getAuth($param);
        $isHasPermission = $auth ? true : false;

        if (!$isHasPermission) {
            return error(2300001, [], [], 401);
        }

        Log::info('Request', [
            'IP'     => $request->getClientIp(),
            'Port'   => $request->getPort(),
            'Method' => $request->method(),
            'Url'    => $request->url(),
            'Params' => $request->all(),
        ]);

        return $next($request);
    }

    /**
     * @param \Illuminate\Http\Request $request
     * @return mixed
     * @Date: 2021/6/30
     * @Time: 19:03
     */
    private function getLoginEnterpriseId($request)
    {
        $token = $request->header('Auth-Token');

        if ($token === null) {
            throw new Error(450000, '请登录后操作');
        }

        try {
            $tokenInfo = json_decode(Crypt::decrypt($token), true);
        } catch (\Exception $exception) {
            throw new Error(450000, 'Auth-Token error: ' . $exception->getMessage());
        }

        //使用redis判断过期时间
        $key = sprintf("gw-scrm:token:%s", $token);

        //15分钟未操作则过期
        if (!Redis::connection()->expire($key, 3600 * 24)) {
            throw new Error(451004, "身份过期", [], 401);
        }

        if (!isset($tokenInfo['is_tmp']) || $tokenInfo['is_tmp'] == 1) {
            throw new Error(451005, '请先换取正式登录token');
        }

        return $tokenInfo['enterprise_id'];
    }
}