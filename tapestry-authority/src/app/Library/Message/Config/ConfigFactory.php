<?php
namespace App\Library\Message\Config;


class ConfigFactory
{
    private static $configPool;

    /**
     * @param $enterpriseId
     * @return Config
     */
    static function getInstance($enterpriseId)
    {
        if (!isset(static::$configPool[$enterpriseId]) || !static::$configPool[$enterpriseId]) {
            static::$configPool[$enterpriseId] = new Config($enterpriseId);
        }
        return static::$configPool[$enterpriseId];
    }
}