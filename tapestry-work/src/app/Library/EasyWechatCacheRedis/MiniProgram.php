<?php

namespace App\Library\EasyWechatCacheRedis;

use EasyWeChat\Factory;
use Symfony\Component\Cache\Adapter\RedisAdapter;

class MiniProgram
{

    public static function init($config): \EasyWeChat\MiniProgram\Application
    {
        $app = Factory::miniProgram($config);
        $app->rebind('cache', new RedisAdapter(app('redis')->connection()->client()));
        return $app;
    }
}
