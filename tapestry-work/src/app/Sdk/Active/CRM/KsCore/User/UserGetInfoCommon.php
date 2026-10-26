<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * User
 * 获取企微用户信息
 */
class UserGetInfoCommon extends Base
{
    public $errStr = 'query user info error';
    public $api    = '/api/WeCom/GetUserInfo?';

    private $stringCode;
    private $stringValue;

    public function __construct(string $stringCode, string $stringValue)
    {
        parent::__construct();
        $this->stringCode  = $stringCode;
        $this->stringValue = $stringValue;
        $this->setQueryParam();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp'       => time(),
            'nonce'           => Config::NONCE,
            $this->stringCode => $this->stringValue,
        ];
    }
}