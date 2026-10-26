<?php

namespace App\Sdk\Active\CRM\CoachCore;


final class Config
{

    const  LOG_NAME = 'coach';
    const  NONCESTR = 'qaD9ehzzqzpijxdDn8p1ppC3acvng76h';
    const  SecretValue = 'REDACTED_FOR_SOURCE_SHARING';
    const  campaignSID = 'FC749E29-292A-40D8-B906-247943D1DF6D';

    // DES加密解密
    const key = 'e233digi';
    const iv = 'E233D1G1';


    // coach测试账号
    static $test = [
        'username' => 'COACH-CHANGYI',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
        'client_id' => 'cdp',
        'client_secret' => 'cdp',

    ];
    // coach生产账号
    static $prod = [
        'username' => 'CHANGYI',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
        'client_id' => 'cdp',
        'client_secret' => 'cdp',

    ];

    static function getTokenUrl(): string
    {
        if (config('app.env') == 'prod') {
            $tokenUrlReq = "https://api.cdp.coach.com.cn";
        } else {
            $tokenUrlReq = "https://apiqa.cdp.coach.com.cn";
        }
        return $tokenUrlReq;
    }

    static function getUrl(): string
    {
        if (config('app.env') == 'prod') {
            $baseUrl = "https://api.cdp.coach.com.cn";
        } else {
            $baseUrl = "https://apiqa.cdp.coach.com.cn";
        }
        return $baseUrl;
    }


    /**
     * 发放优惠卷环境url
     *
     * @return string
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/5
     */
    static function getUrlCore(): string
    {
        if (config('app.env') == 'prod') {
            $baseUrl = "https://weclt.coach.com/weworkapi";
        } else {
            $baseUrl = "https://weclt.coach.com/weworkapitest";
        }
        return $baseUrl;
    }

    /**
     * 消费订单中商品sku——image url
     *
     * @return string
     * @author: wangle@vchangyi.com
     * @Time: 2023/2/28
     */
    static function getSkuImageUrl(): string
    {
        if (config('app.env') == 'prod') {
            $skuImageUrl = "https://s7d2.scene7.com/is/image/Coach";
        } else {
            $skuImageUrl = "https://s7d2.scene7.com/is/image/Coach";
        }
        return $skuImageUrl;
    }

    static function getClientId(): string
    {
        return Config::isProdEnv() ? Config::$prod['client_id'] : Config::$test['client_id'];
    }

    static function isProdEnv(): bool
    {
        return config("app.env") == 'prod';
    }

    static function getClientSecret(): string
    {
        return Config::isProdEnv() ? Config::$prod['client_secret'] : Config::$test['client_secret'];
    }

    static function getUsername(): string
    {
        return Config::isProdEnv() ? Config::$prod['username'] : Config::$test['username'];
    }

    static function getPassword(): string
    {
        return Config::isProdEnv() ? Config::$prod['password'] : Config::$test['password'];
    }

}