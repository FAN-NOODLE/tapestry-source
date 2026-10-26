<?php

namespace App\Sdk\Active\CRM\KsCore\Customer;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * Customer Query Coupons
 * 根据手机号或openid查询KS会员已领取的优惠券信息。手机号和openid同时提供时，优先使用手机号进行查询。
 */
class CustomerQueryCoupons extends Base
{
    public $errStr = 'query customer info error';
    public $api = '/api/customer/querycoupons?';


    private $mobile;
    private $openId;

    public function __construct(string $mobile, string $openId)
    {
        parent::__construct();
        $this->mobile = $mobile;
        $this->openId = $openId;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'Mobile'    => $this->mobile,
            'OpenId'    => $this->openId,
        ];

    }
}