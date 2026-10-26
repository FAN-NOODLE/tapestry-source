<?php

namespace App\Traits;

use CrsSdk\Boss\Enterprise;

/**
 * 企业信息
 * Trait EnterpriseTools
 *
 * @package App\Traits
 */
trait EnterpriseTools
{
    /**
     * 获取切流企业的hash
     *
     * @param $enterpriseId
     *
     * @return string
     */
    public function getEnterpriseHash($enterpriseId)
    {
        $hash = '';

        $enterpriseSdk = new Enterprise();

        try {

            $enterprise = $enterpriseSdk->detail(['enterprise_id' => $enterpriseId])['data'];
            $hash       = $enterprise['hash'];

        } catch (\Exception $e) {

        }

        return $hash;
    }

    /**
     * 获取切流企业
     *
     * @param $enterpriseId
     *
     * @return array
     */
    public function getEnterpriseDetail($enterpriseId)
    {
        $enterprise = '';

        $enterpriseSdk = new Enterprise();

        try {

            $enterprise = $enterpriseSdk->detail(['enterprise_id' => $enterpriseId])['data'];

        } catch (\Exception $e) {

        }

        return $enterprise;
    }
}
