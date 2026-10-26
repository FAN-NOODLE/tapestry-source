<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/15
 * Time: 3:19 PM
 */

namespace App\Http\Middleware;

use Closure;
use Illuminate\Support\Facades\Crypt;

class ReserveAccess
{
    /**
     * 验证预约单微信小程序端的token
     *
     * @param \Illuminate\Http\Request $request
     * @param \Closure $next
     * @return mixed
     */
    public function handle($request, Closure $next)
    {
        $token = $request->header('CustomerToken');
        if (empty($token)) {
            error(652006);
        }
        // 解析Token
        $tokenData = json_decode(Crypt::decrypt($token), true);

        if ($tokenData['expire_time'] < time()) {
            error(652004);
        }
        $request->attributes->add(['tokenData'=>$tokenData]);
        return $next($request);
    }
}
