<?php


namespace App\Library\Ddc;


use App\Exceptions\Error;

class DdcConfig
{

    private static $instance;

    private $host;

    private $authKey;

    private $authSecret;

    /**
     * DDCConfig constructor.
     */
    private function __construct()
    {
        $this->host = rtrim(env('DDC_HOST', ''), '/');
        $this->authKey = env('DDC_AUTH_KEY', '');
        $this->authSecret = env('DDC_AUTH_SECRET', '');

        if (empty($this->host) || empty($this->authKey) || empty($this->authSecret)) {
            throw new Error(110007,"DDC配置错误");
        }
    }


    public static function getInstance()
    {
        if (is_null(self::$instance)) {
            self::$instance = new self();
        }
        return self::$instance;
    }


    /**
     * @return mixed
     */
    public function getHost()
    {
        return $this->host;
    }

    /**
     * @return mixed
     */
    public function getAuthKey()
    {
        return $this->authKey;
    }

    /**
     * @return mixed
     */
    public function getAuthSecret()
    {
        return $this->authSecret;
    }


    private function __clone()
    {

    }
}