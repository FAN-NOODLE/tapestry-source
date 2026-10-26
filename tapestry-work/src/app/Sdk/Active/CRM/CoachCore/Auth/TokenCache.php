<?php

namespace App\Sdk\Active\CRM\CoachCore\Auth;

use App\Sdk\Active\CRM\CoachCore\Request;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Throwable;

class TokenCache
{
    const ACCESS_TOKEN = "CRM:CDP:ACCESS_TOKEN";
    const REFRESH_TOKEN = 'CRM:CDP:REFRESH_TOKEN';
    const LOCK_TTL = 30;
    const LOCK_WAIT_TIMEOUT = 10;
    const LOCK_WAIT_INTERVAL = 100000;

    protected $redisAccessTokenLock = "AccessToken:Lock:";
    private $forceRefresh;
    /**
     * @var Request
     */
    private $request;
    private $enterpriseId;

    public function __construct($enterpriseId, $force = false)
    {
        $this->enterpriseId = $enterpriseId;
        $this->request = new Request($enterpriseId);
        $this->forceRefresh = $force;
        if ($this->forceRefresh) {
            Redis::connection()->del([$this->getAccessTokenKey(), $this->getRefreshTokenKey()]);
        }
    }

    public function getAccessTokenKey(): string
    {
        return static::ACCESS_TOKEN . ":" . $this->enterpriseId;
    }

    public function getRefreshTokenKey(): string
    {
        return static::REFRESH_TOKEN . ":" . $this->enterpriseId;
    }

    /**
     * If token cache is invalid, delete it under lock to avoid duplicate cleanup.
     *
     * @return int
     */
    public function delTokenKey(): int
    {
        if ($this->setRedisTokenLock()) {
            try {
                return Redis::connection()->del([$this->getAccessTokenKey(), $this->getRefreshTokenKey()]);
            } finally {
                $this->removeRedisTokenLock();
            }
        }
        return 0;
    }

    protected function setRedisTokenLock(): int
    {
        $lock = $this->redisAccessTokenLock . $this->enterpriseId;
        $locked = Redis::connection()->setnx($lock, 1);
        if ($locked) {
            Redis::connection()->expire($lock, self::LOCK_TTL);
        }
        return $locked;
    }

    public function getAccessToken(): string
    {
        $token = Redis::connection()->get($this->getAccessTokenKey());
        Log::info("AccessToken", [$token]);
        if (!$token) {
            return $this->getRefreshToken();
        }
        return $token;

    }

    private function getRefreshToken(): string
    {
        $refreshToken = Redis::connection()->get($this->getRefreshTokenKey());
        Log::info("refreshToken", ['enterprise_id' => $this->enterpriseId, 'has_refresh_token' => !empty($refreshToken)]);
        $accessToken = Redis::connection()->get($this->getAccessTokenKey());
        if ($accessToken) return $accessToken;

        if ($this->setRedisTokenLock()) {
            try {
                $accessToken = Redis::connection()->get($this->getAccessTokenKey());
                if ($accessToken) return $accessToken;

                $refreshToken = Redis::connection()->get($this->getRefreshTokenKey());
                if ($refreshToken) {
                    try {
                        $response = $this->request->requestTokenApi(new RefreshAccessToken($this->enterpriseId, $refreshToken));
                    } catch (Throwable $e) {
                        Log::warning("refreshToken failed, retry getAccessToken", [
                            'enterprise_id' => $this->enterpriseId,
                            'message' => $e->getMessage(),
                            'code' => $e->getCode()
                        ]);
                        Redis::connection()->del($this->getRefreshTokenKey());
                        $response = $this->request->requestTokenApi(new GetAccessToken($this->enterpriseId));
                    }
                } else {
                    $response = $this->request->requestTokenApi(new GetAccessToken($this->enterpriseId));
                }

                $token = new AccessToken($response->getBody());
                $this->setAccessAndRefreshToken($token);
                return $token->getAccessToken();
            } finally {
                $this->removeRedisTokenLock();
            }
        }

        $startTime = microtime(true);
        while ((microtime(true) - $startTime) < self::LOCK_WAIT_TIMEOUT) {
            usleep(self::LOCK_WAIT_INTERVAL);
            $accessToken = Redis::connection()->get($this->getAccessTokenKey());
            if ($accessToken) return $accessToken;
        }

        $lock = $this->redisAccessTokenLock . $this->enterpriseId;
        $ttl = Redis::connection()->ttl($lock);
        Log::warning("wait access token lock timeout", ['enterprise_id' => $this->enterpriseId, 'lock_ttl' => $ttl]);
        if ($ttl < 0) {
            $this->removeRedisTokenLock();
        }
        return $this->getRefreshToken();

    }

    /**
     * Set RefreshToken and AccessToken.
     *
     * @param AccessToken $token
     *
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/26
     */
    private function setAccessAndRefreshToken(AccessToken $token): void
    {
        // RefreshToken
        Redis::connection()->setex(
            $this->getRefreshTokenKey(),
            $token->getExpireIn() - 3600,     // refreshToken 3 hours

            $token->getRefreshToken()
        );

//        // Token
        Redis::connection()->setex(
            $this->getAccessTokenKey(),
            7200, // accessToken 2 hours
            $token->getAccessToken()
        );

    }

    protected function removeRedisTokenLock()
    {
        $lock = $this->redisAccessTokenLock . $this->enterpriseId;
        Redis::connection()->del($lock);
    }


}
