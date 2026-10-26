<?php

namespace App\Sdk\Active;

use App\Services\LogService;
use App\Services\SystemService;

class CrmFactory
{
    /**
     * 获取对应SDK
     *
     * @param $config array 配置信息
     *
     * @return mixed
     */
    public static function getInstance(array $config)
    {
        if (!isset($config['name'])){
            error("配置错误");
        }
        switch ($config['name']) {
            case "coach":
                $sdkClass = __NAMESPACE__ . "\CRM\\CoachCrm";
                break;
            case "sw":
                $sdkClass = __NAMESPACE__ . "\CRM\\SwCrm";
                break;
            case "ks":
                $sdkClass = __NAMESPACE__ . "\CRM\\KsCrm";
                break;
            default:
                $sdkClass = __NAMESPACE__ . "\CRM\\{$config['name']}";
        }

        // 类不存在走默认的ERP
        if (!class_exists($sdkClass)) {
            $sdkClass = __NAMESPACE__ . "\CRM\\DefaultCrm";
        }

        return new $sdkClass($config);
    }
}
