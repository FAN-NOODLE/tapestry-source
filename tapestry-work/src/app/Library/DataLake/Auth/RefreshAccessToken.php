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

class RefreshAccessToken extends Base
{
    public $errStr = 'refresh AccessToken error';
    public $api = '/api/interface/oauth/token';

    protected $type = 'form_data';

    protected $tokenUrl = true;
    private $refreshToken;

    function __construct($refreshToken)
    {
        $this->refreshToken = $refreshToken;
        parent::__construct();
    }

    protected function  setQueryParam()
    {
        $this->queryParam =  [
            "grant_type"=>"refresh_token",
            "refresh_token"=> $this->refreshToken
        ];
    }

    public function setHeader()
    {
        $this->header = [
            'auth'=>[$this->config->getHttpBaseClientId(), $this->config->getHttpBaseClientSecret()]
        ];
    }


}