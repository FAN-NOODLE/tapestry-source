<?php

namespace App\Sdk\Active\CRM\KsCore\User;

use App\Sdk\Active\CRM\KsCore\Base;
use App\Sdk\Active\CRM\KsCore\Config;

/**
 * User
 * 通过unionid查询用户是否是会员接口
 */
class UserJoin extends Base
{
    public $errStr = 'query user info error';
    public $api = '/api/WeCom/Joinin?';


    /**
     * UnionId    string    Y    UnionId必须填
     * CampaignCode    String    Y    活动代码(Wedo提供)
     * ShopCode    String    Y    门店code
     * ClerkCode    String    Y    店员code
     */

    //错误代码errcode
    //0表示成功
    //4001表示已经参与领取过
    //4000活动过期

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