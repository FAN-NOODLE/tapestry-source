<?php

namespace App\Sdk\Active\CRM\CoachCore;


final class ConfigSW
{

    const  LOG_NAME = 'sw';

    // sw测试账号
    static $test = [
        'username' => 'SW_CHANGYI',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
        'client_id' => 'cdp',
        'client_secret' => 'cdp',

    ];
    // sw生产账号
    static $prod = [
        'username' => 'SW_CHANGYI',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
        'client_id' => 'cdp',
        'client_secret' => 'cdp',

    ];


    static function getTokenUrl(): string
    {
        if (config('app.env') == 'prod') {
            $tokenUrlReq = "https://api.cdp.stuartweitzman.cn";
        } else {
            $tokenUrlReq = "https://apiqa.cdp.stuartweitzman.cn";
        }
        return $tokenUrlReq;
    }

    static function getUrl(): string
    {
        if (config('app.env') == 'prod') {
            $baseUrl = "https://prod.cdp.stuartweitzman.cn";
        } else {
            $baseUrl = "https://apiqa.cdp.stuartweitzman.cn";
        }
        return $baseUrl;
    }


    static function getClientId(): string
    {
        return ConfigSW::isProdEnv() ? ConfigSW::$prod['client_id'] : ConfigSW::$test['client_id'];
    }

    static function isProdEnv(): bool
    {
        return config("app.env") == 'prod';
    }

    static function getClientSecret(): string
    {
        return ConfigSW::isProdEnv() ? ConfigSW::$prod['client_secret'] : ConfigSW::$test['client_secret'];
    }

    static function getUsername(): string
    {
        return ConfigSW::isProdEnv() ? ConfigSW::$prod['username'] : ConfigSW::$test['username'];
    }

    static function getPassword(): string
    {
        return ConfigSW::isProdEnv() ? ConfigSW::$prod['password'] : ConfigSW::$test['password'];
    }

}