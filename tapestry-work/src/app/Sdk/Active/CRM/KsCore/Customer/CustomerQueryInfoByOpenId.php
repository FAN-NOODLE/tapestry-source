<?php

namespace App\Sdk\Active\CRM\KsCore\Customer;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * 根据OpenId查询会员信息
 */
class CustomerQueryInfoByOpenId extends Base
{
    public $errStr = 'query customer info error';
    public $api = '/api/customer/query?';


    private $openId;

    public function __construct(string $openId)
    {
        parent::__construct();
        $this->openId = $openId;
        $this->setQueryParam();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'OpenId'    => $this->openId,
        ];
    }
}