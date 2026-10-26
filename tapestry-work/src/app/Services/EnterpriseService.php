<?php

namespace App\Services;

use App\Exceptions\Error;
use App\Models\StoreDepartment;
use CrsSdk\Boss\BaseClient;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\Boss\EnterpriseWorkPlugin;

class EnterpriseService extends BaseService
{


    /**
     * 根据企业id获取企业信息
     * @param $enterpriseId
     * @return mixed
     * @throws \Exception
     */
    public function getDetailByEnterpriseId($enterpriseId)
    {
        $enterpriseData = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
        return $enterpriseData['data'];
    }

    public static function EpId2Domain($enterpriseId)
    {
        $service = new self();
        if($info = $service->getDetailByEnterpriseId($enterpriseId)){
            return $info['domain'];
        }
        throw new Error(10000,"企业id查不到企业hash");
    }

    /**
     * 查询企业domain获取企业详情
     *
     * @param $domain
     *
     * @return mixed
     * @throws \Exception
     */
    public function getDomainInfo($domain)
    {
        $lockKey = redisKey('ENTERPRISE_DOMAIN_INFO', ['domain' => $domain]);

        // 查询缓存
        $enterpriseData = getOrSetData($lockKey);
        if (empty($enterpriseData)) {
            $enterpriseData = (new Enterprise())->detail(['domain' => $domain]);
            // 获取该企业的授权信息
            $enterpriseData = $enterpriseData['data'];
            getOrSetData($lockKey, $enterpriseData, 86400);
        }

        return $enterpriseData;
    }


    /**
     * 查询企业hash获取企业详情
     *
     * @param $enterpriseHash
     *
     * @return mixed
     * @throws \Exception
     */
    public function getHashInfo($enterpriseHash)
    {
        $lockKey = redisKey('ENTERPRISE_HASH_INFO', ['hash' => $enterpriseHash]);

        // 查询缓存
        $enterpriseData = getOrSetData($lockKey);
        if (empty($enterpriseData)) {
            $enterpriseData = (new Enterprise())->detail(['hash' => $enterpriseHash]);
            // 获取该企业的授权信息
            getOrSetData($lockKey, $enterpriseData, 86400);
        }

        return $enterpriseData;
    }

    /**
     * 查询企业hash获取企业详情
     *
     * @param $corpId
     * @return mixed
     * @throws \Exception
     */
    public function getCorpIdInfo($corpId)
    {
        $lockKey = redisKey('ENTERPRISE_CORP_ID_INFO', ['corp_id' => $corpId]);

        // 查询缓存
        $enterpriseData = getOrSetData($lockKey);
        if (empty($enterpriseData)) {
            $data               = (new Enterprise())->detail([
                'corp_id' => $corpId
            ])['data'];
            // 获取该企业的授权信息
            getOrSetData($lockKey, $data, 86400);
        }

        return $enterpriseData;
    }

    /**
     * getEnterpriseIdentifierList
     * @desc 获取企业标识列表
     * @param $enterpriseId
     *
     * @return array
     * @date 2020-07-25
     * @throws \Exception
     */
    public static function getEnterpriseIdentifierList($enterpriseId) {
        $lockKey = redisKey('ENTERPRISE_WORK_PLUGIN', ['id' => $enterpriseId]);
        // 查询缓存
        $usePluginList = getOrSetData($lockKey);
        if (empty($usePluginList)) {
            $usePluginList = (new EnterpriseWorkPlugin())->lists([
                'enterprise_id' => $enterpriseId
            ]);
            getOrSetData($lockKey, $usePluginList, 86400);
        }
        return $usePluginList;
    }


    public static function getEnterpriseConfig($enterpriseId, $key) {
        $lockKey = redisKey('ENTERPRISE_CONFIG_INFO', [$key => $enterpriseId]);
        // 查询缓存
        $configList = getOrSetData($lockKey);
        if (empty($configList)) {
            $configList = (new EnterpriseConfig())->info([
                'enterprise_id' => $enterpriseId,
                'key'           => $key
            ]);
            getOrSetData($lockKey, $configList, 86400);
        }
        return $configList;
    }


    /**
     * 查询企业hash获取企业详情
     *
     * @param $enterpriseId
     * @return mixed
     * @throws \Exception
     */
    public function detail($enterpriseId)
    {
        $lockKey = redisKey('ENTERPRISE_ID_INFO', ['enterprise_id' => $enterpriseId]);
        // 查询缓存
        $enterpriseData = getOrSetData($lockKey);
        if (empty($enterpriseData)) {
            $enterpriseData               = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);
            // 获取该企业的授权信息
            getOrSetData($lockKey, $enterpriseData, 86400);
        }

        return $enterpriseData;
    }


}
