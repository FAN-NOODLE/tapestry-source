<?php

namespace App\Sdk\CDP\Coach;
final class Config
{
    // coach的Ma测试账号
    static $test = [
        'domain' => 'https://uat-recommend-newcdp.coach.com.cn',
        'token' => 'REDACTED_FOR_SOURCE_SHARING',
    ];

    // coach生产账号
    static $prod = [
        'domain' => 'https://recommend-newcdp.coach.com.cn',
        'token' => 'REDACTED_FOR_SOURCE_SHARING',
    ];


    public static function getConfig(): array
    {
        if (config("app.env") == "prod") {
            return self::$prod;
        } else {
            return self::$test;
        }
    }

    public static function getEasConfig(): array
    {
        return [
            'url' => env('EAS_RECOMMEND_URL', ''),
            'authorization' => env('EAS_RECOMMEND_AUTHORIZATION', ''),
        ];
    }
}
