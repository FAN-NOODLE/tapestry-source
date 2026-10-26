<?php

namespace App\Library\Message\Config;


use App\Exceptions\Error;
use App\Models\ShortMessageConfig;

class Config
{
    private $config;
    private $enterpriseId;

    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
        $config             = ShortMessageConfig::query()
            ->where('enterprise_id', $enterpriseId)
            ->first();
        if (!$config) {
            throw new Error(100000, '企业的短信配置未找到');
        }
        $this->config = $config->toArray();
    }

    public function getConfig()
    {
        return json_decode(unserialize($this->config['config']), true);
    }

    public function getType()
    {
        return $this->config['sms_type'];
    }

}