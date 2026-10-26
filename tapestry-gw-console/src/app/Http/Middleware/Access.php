<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/15
 * Time: 3:19 PM
 */

namespace App\Http\Middleware;

use Closure;
use Illuminate\Support\Facades\Log;

class Access
{
    /**
     * 记录访问日志
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        // 全局request_id
        $headers    = $request->header();
        $xRequestId = '';
        foreach ($headers as $key => $val) {
            if (strtolower($key) == 'x-request-id') {
                $xRequestId = $val[0];
                break;
            }
        }
        if (empty($xRequestId)) {
            $xRequestId = LOG_ID;
        }
        $_SERVER['x_request_id'] = $xRequestId;

        Log::info('Request', [
            'IP'     => $request->getClientIp(),
            'Port'   => $request->getPort(),
            'Method' => $request->method(),
            'Url'    => $request->url(),
            'Params' => $request->all(),
        ]);

        return $next($request);
    }
}
