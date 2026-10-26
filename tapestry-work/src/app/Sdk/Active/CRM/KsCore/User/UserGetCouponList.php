<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * 通过unionid获取用户优惠券列表信息
 */
class UserGetCouponList extends Base
{
    public $errStr = 'query user info error';
    public $api = '/api/WeCom/GetUserCouponList?';


    private $unionId;

    public function __construct(string $unionId)
    {
        parent::__construct();
        $this->unionId = $unionId;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'unionId'   => $this->unionId,
        ];

    }
}