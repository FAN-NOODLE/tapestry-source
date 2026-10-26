<?php

namespace App\Sdk\Active;


use CrsSdk\Boss\EnterpriseConfig;

class CrmHelper
{

    // 企业名称
    const STORE_NAME_KEY = 'store_name';

    // 企业配置的redis-key
    const REDIS_ENTERPRISE_CONFIG_KEY = "CRM:ENTERPRISE_CONFIG_KEY";

    const PrintEcho = true;

    /**
     * 测试用函数,打印请求和返回
     *
     * @param string $log
     * @param string $title
     * @param array  $context
     */
    public static function echoLine(string $log, string $title, array $context = [])
    {
        if (self::PrintEcho && config('app.env') == 'dev') {
            $str     = PHP_EOL;
            $str     .= $title . PHP_EOL;
            $str     .= var_export($context, true);
            $str     .= PHP_EOL;
            $logPath = env('LOGGING_PATH', sprintf('/data/cloud-log/%s', date('Ymd')));
            if (!is_dir($logPath))
                mkdir($logPath, 0777, true);
            file_put_contents($logPath . '/log_crm_' . $log . '.log', $str, FILE_APPEND);
        }
    }

    /**
     * @param string $prefix
     *
     * @return string
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/26
     */
    public static function createUid(string $prefix = ''): string
    {

        $chars = md5(uniqid(mt_rand(), true));
        $uuid  = substr($chars, 0, 8) . '-'
            . substr($chars, 8, 4) . '-'
            . substr($chars, 12, 4) . '-'
            . substr($chars, 16, 4) . '-'
            . substr($chars, 20, 12);
        return $prefix . $uuid;

    }


    /**
     * 获取缓存企业配置
     *
     * @param $enterpriseId
     *
     * @return bool|mixed
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/18
     */
    public static function getEnterpriseConfig($enterpriseId)
    {
        $redisKey         = $enterpriseId . ":" . self::REDIS_ENTERPRISE_CONFIG_KEY;
        $enterpriseConfig = getOrSetData($redisKey);
        if (!$enterpriseConfig) {
            // 根据企业查询企业配置
            $configList = (new EnterpriseConfig())->info([
                'enterprise_id' => $enterpriseId,
                'key'           => self::STORE_NAME_KEY
            ]);

            $enterpriseConfig                  = current($configList['data'] ?? []);
            $enterpriseConfig['enterprise_id'] = $enterpriseId;
            getOrSetData($redisKey, $enterpriseConfig, 86400);

        }

        return $enterpriseConfig;
    }
}
