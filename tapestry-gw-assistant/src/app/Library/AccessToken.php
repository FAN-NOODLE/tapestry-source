<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019-07-11
 * Time: 10:20
 */

namespace App\Library;

use CrsSdk\Boss\Auth;
use EasyWeChat\Factory;

class AccessToken
{
    /**
     * 获取access_token
     *
     * @param int    $enterpriseId 企业id
     * @param string $appId 小程序app_id
     * @param string $appSecret 小程序app_secret
     * @param string $overdueToken 过期的token
     *
     * @return mixed
     * @throws \EasyWeChat\Kernel\Exceptions\HttpException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidArgumentException
     * @throws \EasyWeChat\Kernel\Exceptions\InvalidConfigException
     * @throws \EasyWeChat\Kernel\Exceptions\RuntimeException
     * @throws \Psr\SimpleCache\InvalidArgumentException
     */
    public function getAccessToken($enterpriseId, $appId = '', $appSecret = '', $overdueToken = '')
    {
        $redisKey = redisKey('ACCESS_TOKEN_BY_ENTERPRISE_ID', ['enterprise_id' => $enterpriseId]);

        $accessToken = getOrSetData($redisKey);
        if (!empty($accessToken) && $overdueToken != $accessToken) {
            return $accessToken;
        }else{
            delRedis([$redisKey]);
        }

        $config = $this->getConfig($enterpriseId, $appId, $appSecret);

        // 获取access_token
        $accessToken = Factory::miniProgram($config)->access_token->getToken(true);

        $accessToken = $accessToken['access_token'];

        getOrSetData($redisKey, $accessToken, 3600);

        return $accessToken;
    }

    public function getConfig($enterpriseId, $appId = '', $appSecret = '')
    {
        if (empty($appId) && empty($appSecret)) {
            // 获取商城小程序的id和secret
            $authAccount = (new Auth())->getConfig(['enterprise_id' => $enterpriseId]);

            $appId     = $authAccount['data']['app_id'];
            $appSecret = $authAccount['data']['app_secret'];
        }

        // 参数
        return [
            'app_id' => $appId,
            'secret' => $appSecret
        ];
    }

}
