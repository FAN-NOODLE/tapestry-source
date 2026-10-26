<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/19 11:30
 */

namespace App\Services;


use App\Common\Constant;
use App\Exceptions\Error;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class Decrypt
{
    public function index($cache)
    {
        try {
            $cache = Crypt::decrypt($cache);
            $cache = json_decode($cache, true);
        } catch (\Exception $exception) {
            Log::error('token:' . var_export($cache, true));
            Log::error('ERROR:110003:token解密错误');
            throw new Error(110003, 'token解密错误');
        }
        return $cache;
    }

    /**
     * 获得小程序配置
     * @param $identifier
     * @return array
     * @throws \Exception
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2020/11/25 16:46
     */
    public static function getEnvMiniProgram($identifier)
    {
        switch ($identifier) {
            case Constant::STANDARD_IDENTIFIER :
                //正式appId
                $appId      = env('PROGRAM_APPID');
                $secret     = env('PROGRAM_SECRET');
                break;
            case Constant::RELEASE_IDENTIFIER :
                //灰度appId
                $appId      = env('RELEASE_PROGRAM_APPID');
                $secret     = env('RELEASE_PROGRAM_SECRET');
                break;
            case Constant::STANDARD_ICLK_IDENTIFIER :
                //灰度appId
                $appId      = env('ICLK_PROGRAM_APPID');
                $secret     = env('ICLK_PROGRAM_SECRET');
                break;
            default :
                //正式appId
                $appId      = env('PROGRAM_APPID');
                $secret     = env('PROGRAM_SECRET');
        }
        return [$appId,$secret];
    }
}