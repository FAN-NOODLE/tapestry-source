<?php

namespace App\Sdk\Active\CRM\KsCore\Customer;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * 根据Mobile查询会员详细信息
 */
class CustomerQueryInfoByMobile extends Base
{
    public $errStr = 'query customer info error';
    public $api = '/api/customer/query?';

    private $mobile;

    public function __construct(string $mobile)
    {
        parent::__construct();
        $this->mobile = $mobile;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'Mobile'    => $this->mobile,
        ];
    }
}