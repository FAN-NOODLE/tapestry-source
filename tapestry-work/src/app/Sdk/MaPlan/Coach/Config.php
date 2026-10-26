<?php

namespace App\Sdk\MaPlan\Coach;
final class Config
{
    const  LOG_NAME = 'coach';

    // coach的Ma测试账号
    static $test = [
        'username' => 'cl0168186eea8db5f',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
//        'username' => 'cl03b1188335d24d8',
//        'password' => 'REDACTED_FOR_SOURCE_SHARING',
    ];
    // coach生产账号
    static $prod = [
        'username' => 'cl03b1188335d24d8',
        'password' => 'REDACTED_FOR_SOURCE_SHARING',
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