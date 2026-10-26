<?php

namespace App\Http\Middleware;

use Closure;

class SignMiddleware
{
    /**
     * 调试签名
     */
    const DEBUG_SIGN = 'CRS';

    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     *
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        // 开发环境不校验签名
        if (in_array(config('app.env'), ['dev', 'test'])) {
            return $next($request);
        }

        // 便于调试
        if ($request->input('sign') === 'debug') {
            return $next($request);
        }

        $sign      = $request->header('Auth-Sign');
        $client    = $request->header('Auth-Client');
        $timestamp = $request->header('Auth-Timestamp');

        if (empty($sign)) {
            return error(100001, [], 401);
        }

        if (empty($timestamp)) {
            return error(100003, [], 401);
        }

        $key = config('client.key');

        if (empty($key)) {
            return error(100004, [], 401);
        }

        $verifySign = md5($client . $key . $timestamp);
        $difference = time() - $timestamp;

        if (abs($difference) > 3600) {
            return error(100005, [], 401);
        }

        if ($sign !== $verifySign) {
            return error(100006, [], 401);
        }

        return $next($request);
    }
}
