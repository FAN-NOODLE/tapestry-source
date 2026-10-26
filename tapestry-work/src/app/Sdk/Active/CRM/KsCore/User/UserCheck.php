<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * User
 * 通过unionid查询用户是否是会员接口
 */
class UserCheck extends Base
{
    public $errStr = 'query user info error';
    public $api = '/api/WeCom/Check?';


    // UnionId	string	Y	UnionId必须填
    // 参数名	类型	说明
    // data	bool	True代表存在 false代表不存在
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