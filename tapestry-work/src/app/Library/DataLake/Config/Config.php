<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/26 0026
 * Time: 11:41
 */

namespace App\Library\DataLake\Config;


use App\Exceptions\Error;
use App\Models\DataLakeConfig;
use Illuminate\Support\Facades\Redis;
use App\Common\Constant;

class Config
{
    private $config;
    private $enterpriseId;
    const CHANNEL = "6P0";

    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;

    }

    private function getEnterpriseDLConfig(int $enterpriseID): array
    {
        $config = DataLakeConfig::query()
            ->where('enterprise_id', $enterpriseID)
            ->first();
//        if (!$config) {
//            throw new Error(100000, '企业dataLake配置找不到');
//        }
//
//        return $config->toArray();
        return $this->config = $config?$config->toArray():[];
    }

    public function getBrandCode()
    {
        return $this->config['brand_code'];
    }

    public function getProgramCode()
    {
        return $this->config['program_code'];
    }

    public function getStoreCode()
    {
        return $this->config['store_code'];
    }

    public function getSellerName()
    {
        return $this->config['seller_name'];
    }

    public function getHttpBaseClientId()
    {
        return $this->config['client_id'];
    }

    public function getHttpBaseClientSecret()
    {
        return $this->config['client_secret'];
    }

    public function getTokenAuthUser()
    {
        return $this->config['auth_user'];
    }

    public function getTokenAuthPassWord()
    {
        return $this->config['auth_password'];
    }

    public function getEnterpriseId()
    {
        return $this->enterpriseId;
    }

    public function getChannelCode()
    {
        return '6P0';
    }

}
