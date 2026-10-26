<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 11:40
 */

namespace App\Library\DataLake\Config;


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