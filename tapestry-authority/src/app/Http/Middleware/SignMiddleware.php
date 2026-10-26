<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;

class SignMiddleware
{
    /**
     * Handle an incoming request.
     *
     * @param Request $request
     * @param Closure $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        // 开发环境不校验签名
        if (config('app.env') === 'dev') {
            return $next($request);
        }

        // 测试环境不校验签名
        if (config('app.env') === 'test') {
            return $next($request);
        }

        // 便于调试
        if ($request->input('sign') === 'CRS.Micro') {
            return $next($request);
        }

        $sign      = $request->header('Auth-Sign');
        $client    = $request->header('Auth-Client');
        $timestamp = $request->header('Auth-Timestamp');

        if (empty($sign)) {
            notice(100001);
        }

        if (empty($timestamp)) {
            notice(100003);
        }

        $key = config('client.key');

        if (empty($key)) {
            notice(100004);
        }

        $verifySign = md5($client . $key . $timestamp);
        $difference = time() - $timestamp;

        if ($difference > 3600 || $difference < -5) {
            notice(100005);
        }

        if ($sign !== $verifySign) {
            notice(100006);
        }

        return $next($request);
    }
}
