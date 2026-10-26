<?php

namespace App\Sdk\Active\CRM\KsCore\Customer;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * 根据UnionId查询会员信息
 */
class CustomerQueryInfo extends Base
{
    public $errStr = 'query customer info error';
    public $api = '/api/customer/query?';


    private $mobile;
    private $openId;
    private $unionId;

    public function __construct(string $mobile, string $openId, string $unionId)
    {
        parent::__construct();
        $this->mobile  = $mobile;
        $this->openId  = $openId;
        $this->unionId = $unionId;
        $this->setQueryParam();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'Mobile'    => $this->mobile,
            'OpenId'    => $this->openId,
            'UnionId'   => $this->unionId,
        ];
    }
}