<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class Access
{
    /**
     * 记录访问日志
     *
     * @param Request $request
     * @param \Closure $next
     *
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        // 全局request_id
        $headers = $request->header();
        $xRequestId = '';
        foreach ($headers as $key => $val) {
            if (strtolower($key) == 'x-request-id') {
                $xRequestId = $val[0];
                break;
            }
            if ($key === 'enterprise-hash') {
                $_SERVER['enterprise_hash'] = $val;
            }
        }
        if (empty($xRequestId)) {
            $xRequestId = LOG_ID;
        }
        $_SERVER['x_request_id'] = $xRequestId;

        $response = $next($request);

        try {
            $responseData = json_decode($response->getContent(), true);

            if (is_null($responseData) || !is_array($responseData)) {
                $responseData = [$responseData];
            }

            Log::info('请求全部数据', [
                'url'         => $request->url(),
                'method'      => $request->method(),
                'params'      => $request->all(),
                'return_code' => $response->getStatusCode(),
                'return_body' => $responseData,
            ]);

        } catch (\Exception $exception) {

            Log::info('请求全部数据', [
                'url'    => $request->url(),
                'method' => $request->method(),
                'params' => $request->all(),
            ]);
        }

        return $response;
    }
}
