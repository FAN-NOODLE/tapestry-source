<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/12 0012
 * Time: 14:39
 */

namespace App\Library\DataLake;


use App\Library\DataLake\Auth\TokenCache;
use App\Library\DataLake\Config\Config;
use App\Library\DataLake\Config\ConfigFactory;

abstract class Base
{
    protected $baseUrl = 'https://dl-api.lorealchina.com';

    protected $errStr = '';

    protected $queryParam = [];

    protected $header = [];

    protected $api = '';

    protected $throttle = false;

    protected $type = 'json';

    protected $tokenUrl = false;

    /**
     * @var Config
     */
    protected $config;

    protected $token;

    function __construct()
    {
    }

    public function notifyEnterpriseId($enterpriseId)
    {

    }

    protected abstract function setQueryParam();

    protected function getApi()
    {
        if (config('app.env') == 'prod' && !$this->tokenUrl) {
            $prefix = "/api/interface/third/";
            return $prefix . 'va/' . substr($this->api, strlen($prefix));
        }

        if (config('app.env') != 'prod' && !$this->tokenUrl) {
            $prefix = "/api/interface/third/";
            return $prefix . 'vb/' . substr($this->api, strlen($prefix));
        }
        return $this->api;

    }

    public function getFullUrl()
    {
        return $this->baseUrl . $this->getApi();
    }

    public function getQueryParam()
    {
        if ($this->type == 'json') {
            return ['body' => json_encode($this->queryParam)];
        } else if ($this->type == 'form_data') {
            return ['form_params' => $this->queryParam];
        }
    }

    public function getOptions()
    {
        $body   = $this->getQueryParam();
        $header = $this->getHeader();
        return array_merge($header, $body);
    }

    public function getHeader()
    {
        return $this->header;
    }

    public function getErrorStr()
    {
        return $this->errStr;
    }

    public function getThrottle()
    {
        return $this->throttle;
    }


    public function setHeader()
    {
        $enterpriseId = $this->config->getEnterpriseId();
        $headers      = [
            'Authorization' => "Bearer " . (new TokenCache($enterpriseId))->getAccessToken(),
            'Content-Type'  => 'application/json',
        ];
        if ($this->throttle) {
            $headers['ruid'] = json_encode(['uuid' => $this->createUid()]);

        }
        $this->header = ['headers' => $headers, 'http_errors' => false];

    }

    public function createUid($prefix = '')
    {

        $chars = md5(uniqid(mt_rand(), true));
        $uuid  = substr($chars, 0, 8) . '-'
            . substr($chars, 8, 4) . '-'
            . substr($chars, 12, 4) . '-'
            . substr($chars, 16, 4) . '-'
            . substr($chars, 20, 12);
        return $prefix . $uuid;

    }


}