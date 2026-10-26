<?php


namespace App\Micros;

use App\Exceptions\Error;
use CrsSdk\Base;
use Exception;
use Illuminate\Support\Str;

/**
 * @method  static self MicroScrm()
 * @method  static self work()
 * @method  static self Authority()
 */
class Common extends Base
{
    private static $aliasMap = [
        'microscrm' => 'go-micro-scrm',
        'work'      => 'work',
        'authority' => 'micro-authority-v1'
    ];
    private $method = 'post';
    private $params;
    private $api;

    public function __construct($clientAlias)
    {
        parent::__construct($clientAlias);
    }

    static function __callStatic($method, $params): self
    {

        $clientAlias = static::$aliasMap[Str::lower($method)] ?? null;
        if ($clientAlias) {
            return new self($clientAlias);
        }
        throw new Exception($method . '对应端未找到', 999);
    }

    public function setParams($params = [])
    {
        $this->params = $params;
        return $this;
    }

    public function setApi($api)
    {
        $this->api = $api;
        return $this;
    }

    public function get($param = [])
    {
        $this->method = 'get';

        return $this->query($param);
    }

    public function query($param = [])
    {
        if (!$this->api) {
            throw new Error(10000, '请求地址为空');
        }
        return $this->request($this->method, $this->api, $this->params ?? $param);
    }

    public function post($param = [])
    {
        $this->method = 'post';
        return $this->query($param);
    }

}