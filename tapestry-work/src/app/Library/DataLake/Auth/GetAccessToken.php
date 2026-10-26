<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/12 0012
 * Time: 14:38
 */

namespace App\Library\DataLake\Auth;


use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

/**
 * Class GetAccessToken
 * @package App\Library\LorealCRM\Auth
 * 注token 过期时间为半小时
 */
class GetAccessToken extends Base
{

    public $errStr = 'get AccessToken error';
    public $api = '/api/interface/oauth/token';

    protected $type = 'form_data';

    protected $tokenUrl = true;

    protected function setQueryParam()
    {
        $this->queryParam = [
            "grant_type" => "password",
            "username"   => $this->config->getTokenAuthUser(),
            "password"   => $this->config->getTokenAuthPassWord()
        ];
    }


    public function setHeader()
    {
        $this->header = [
            'auth' => [$this->config->getHttpBaseClientId(), $this->config->getHttpBaseClientSecret()]
        ];
    }


}