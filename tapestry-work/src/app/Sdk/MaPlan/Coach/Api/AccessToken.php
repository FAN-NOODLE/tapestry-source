<?php

namespace App\Sdk\MaPlan\Coach\Api;

use App\Library\WechatNotice\Chat;
use App\Sdk\MaPlan\Coach\Request;
use Illuminate\Support\Facades\Redis;

class AccessToken extends Base
{

    protected $uri = "v2/oauth2/token?grant_type=client_credentials";
    protected $access_token;


    protected $prefix = "ma:access_token:";

    public function __construct($enterpriseId,$originType)
    {
        parent::__construct($enterpriseId,$originType);
        $this->uri = $this->apiUri . "/" . $this->uri;
    }

    public function getAccessToken()
    {
        $key = $this->prefix . $this->enterpriseId;
        $accessToken = Redis::connection()->get($key);
        if (!empty($accessToken)) {
            return $accessToken;
        }
        $response = Request::Instance($this->enterpriseId,$this->originType)->getAccessToken($this->uri);
        if (!isset($response['success']) && !empty($response['error'])) {
            Chat::getInstance()->send('请求MA接口返回错误码', [
                'env' => config('app.env'),
                'api' => $this->uri,
                "error_code" => $response['error']['code'],
                'error_info' => $response['error']['message']
            ]);
            return "";
        }
        $this->access_token = $response['access_token'];
        Redis::connection()->setex($key, $response['expires_in'] - 300, $this->access_token);
        return $this->access_token;
    }


}