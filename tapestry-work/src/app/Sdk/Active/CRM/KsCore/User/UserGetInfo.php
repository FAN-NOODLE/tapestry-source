<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * User
 * 通过unionid获取企微用户信息
 */
class UserGetInfo extends Base
{
    public $errStr = 'query user info error';
    public $api = '/api/WeCom/GetUserInfo?';


    // UnionId	string	Y	UnionId必须填
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