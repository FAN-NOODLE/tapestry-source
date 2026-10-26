<?php

namespace App\Http\Middleware;

use Closure;

class SignMiddleware
{
    /**
     * Handle an incoming request.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        // 开发环境不校验签名
        $env = config('app.env');
        if (in_array($env, ['dev', 'test', 'uat'])) {
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
            error(100001, [], 401);
        }

        if (empty($timestamp)) {
            error(100003, [], 401);
        }

        $key = config('client.key');

        if (empty($key)) {
            error(100004, [], 401);
        }

        $verifySign = md5($client . $key . $timestamp);
        $difference = time() - $timestamp;

        if ($difference > 3600 || $difference < -3600) {
            error(100005, [], 401);
        }

        if ($sign !== $verifySign) {
            error(100006, [], 401);
        }

        return $next($request);
    }
}
