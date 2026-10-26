<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class Access
{
    /**
     * 记录访问日志
     *
     * @param Request $request
     * @param Closure $next
     * @return mixed
     */
    public function handle(Request $request, Closure $next)
    {
        // 获取短链接参数
        $traceId = $request->header('Trace-Id', 0);
        $traceState = $request->header('Trace-State', 0);

        // 是否需要解密
        if ($traceId && $traceState) {
            $decryptData = decryptShortLink($traceId);
            if (!empty($decryptData['args'])) {
                $request->offsetSet('traceId', $traceId);
                foreach ($decryptData['args'] as $v) {
                    $request->offsetSet($v['name'], $v['value']);
                }
            }
        }

        $response = $next($request);

        Log::debug('Request Url: ' . $request->url());
        Log::debug('Request Method: ' . $request->method());
        Log::debug('Request Params: ' . json_encode($request->all()));

        // 解析Token
        if ($token = $request->header('token')) {
            $tokenData = json_decode(Crypt::decrypt($token), true);
            $tokenTtl = $tokenData['expiredTime'] - time();
            if ($tokenTtl < 900) {
                $response->header('Auth-Status', 'false');
            }
        } else {
            $response->header('Auth-Status', 'false');
        }
        return $response;
    }
}
