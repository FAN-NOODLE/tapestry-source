<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * 通过unionid获取用户订单列表信息
 */
class UserGetOrderList extends Base
{
    public $errStr = 'query user info error';
    public $api = '/api/WeCom/GetUserOrderList?';

    /**
     * UnionId    string    Y    UnionId必须填
     * Index    int    Y    页数 默认1
     * Size    Int    Y    取的数量 默认10
     */
    private $unionId;
    private $index = 1;
    private $size = 10;

    public function __construct($unionId, $index, $size)
    {
        parent::__construct();
        $this->unionId = $unionId;
        $this->index   = $index ?? $this->index;
        $this->size    = $size ?? $this->size;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'timestamp' => time(),
            'nonce'     => Config::NONCE,
            'unionId'   => $this->unionId,
            'index'     => $this->index,
            'size'      => $this->size,
        ];

    }
}