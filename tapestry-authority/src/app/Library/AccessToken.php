<?php

namespace App\Library;

use Illuminate\Support\Facades\Log;
use Overtrue\LaravelWeChat\Facade as LaravelWechat;
use Psr\SimpleCache\InvalidArgumentException;

class AccessToken
{
    /**
     * 获取access_token
     *
     * @param $appId
     * @param $appSecret
     * @param false $refresh
     * @return string
     * @throws InvalidArgumentException
     */
    public function getAccessToken($appId, $appSecret, $refresh = false)
    {
        try {
            $miniProgram = LaravelWechat::miniProgram();
            $miniProgram->config['app_id'] = $appId;
            $miniProgram->config['secret'] = $appSecret;
            $accessToken = $miniProgram->access_token->getToken($refresh);
        } catch (\Exception $e) {
            Log::error('请求微信获取Token异常', [
                'message' => $e->getMessage(),
                'code'    => $e->getCode(),
            ]);
            $accessToken['access_token'] = '';
        }

        return $accessToken['access_token'];
    }
}
