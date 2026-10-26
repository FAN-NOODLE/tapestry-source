<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/16
 * @Time: 12:18
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: BaseHandlers.php
 */

namespace App\Library\Message\Handlers;


use App\Library\Message\Contracts\SendHandlerInterface;
use App\Library\Message\Config\Config;

abstract class BaseHandlers implements SendHandlerInterface
{
    /**
     * 短信配置
     * @var array|null
     */
    protected $config = null;

    public function __construct(array $config)
    {
        $this->config = $config;
    }

    public function getConfig()
    {
        return $this->config;
    }

    public function setConfig(Config $config)
    {
        $this->config = $config;

        return $this;
    }

    public function checkMobile($mobile)
    {

    }
}