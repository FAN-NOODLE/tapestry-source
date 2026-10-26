<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/12 0012
 * Time: 16:41
 */

namespace App\Library\DataLake\Auth;


use App\Library\DataLake\Request;
use Illuminate\Support\Facades\Redis;

class TokenCache
{
    const ACCESS_TOKEN = "LOREAL:CRM:ACCESS_TOKEN";
    const REFRESH_TOKEN = 'LOREAL:CRM:REFRESH_TOKEN';

    private $forceRefresh = false;
    private $enterpriseId;

    private function getAccessTokenKey()
    {
        return static::ACCESS_TOKEN.':'.$this->enterpriseId;
    }

    private function getRefreshTokenKey()
    {
        return static::REFRESH_TOKEN.':'.$this->enterpriseId;
    }


    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
        $this->request = new Request($enterpriseId);
        if ($this->forceRefresh) {
            Redis::connection()->del([$this->getAccessTokenKey(), $this->getRefreshTokenKey()]);
        }
    }

    public function getAccessToken()
    {
        if (!($token = Redis::connection()->get($this->getAccessTokenKey()))) {
            return $this->getRefreshToken()->getAccessToken();
        }
        //echo "直接读取";
        return $token;

    }

    private function getRefreshToken(): AccessToken
    {

        if (!($refreshToken = Redis::connection()->get($this->getRefreshTokenKey()))) {
            //echo "重新登陆";
            $token = $this->request->requestApi(new GetAccessToken());
        } else {
            //echo "刷新token";
            $token = $this->request->requestApi(new RefreshAccessToken($refreshToken));
        }

        $token = new AccessToken($token->getBody());

        $this->setAccessAndRefreshToken($token);

        return $token;

    }

    private function setAccessAndRefreshToken(AccessToken $token): void
    {

        Redis::connection()->setex(
            $this->getRefreshTokenKey(),
            2 * 3600 - 5 * 60, // refreshToken 两小时过期
            $token->getRefreshToken()
        );

        Redis::connection()->setex(
            $this->getAccessTokenKey(),
            $token->getExpireIn() - 5 * 60, // accessToken提前五分钟刷新
            $token->getAccessToken()
        );

    }


}