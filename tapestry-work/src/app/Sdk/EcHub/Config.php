<?php

namespace App\Sdk\EcHub;

class Config
{
    // coach的Ma测试账号
    static $test = [
        'appkey' => 'REDACTED_FOR_SOURCE_SHARING',
        'appsecret' => 'REDACTED_FOR_SOURCE_SHARING',
        'baseUrl'=>"https://ec-hub-stage.coach.com.cn"
    ];
    // coach生产账号
    static $prod = [
        'appkey' => 'REDACTED_FOR_SOURCE_SHARING',
        'appsecret' => 'REDACTED_FOR_SOURCE_SHARING',
        'baseUrl'=>"https://ec-hub.coach.com.cn"
    ];


    public static function getInfo(): array
    {
        if (config("app.env") == "prod") {
            return self::$prod;
        } else {
            return self::$test;
        }
    }

}