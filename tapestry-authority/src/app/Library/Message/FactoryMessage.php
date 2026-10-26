<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-14
 * Time: 17:08
 */

namespace App\Library\Message;

use App\Library\Message\Contracts\SendHandlerInterface;
use App\Library\Message\Handlers\BlueFocusHandlers;
use App\Library\Message\Handlers\MWHandler;
use App\Library\Message\Handlers\SmsHandlers;
use App\Library\Message\Handlers\HRMessageHandlers;
use App\Library\Message\Config\ConfigFactory;
use App\Library\Message\Handlers\WemediacnHandlers;

class FactoryMessage extends MessageBase
{
    /**
     * 默认可以发送的类型
     *
     * @var array
     */
    private $handlers = [
        'easySms'   => SmsHandlers::class,
        'HRMessage' => HRMessageHandlers::class,
        'MWMessage' => MWHandler::class,
        'BlueFocus' => BlueFocusHandlers::class,
        'Wemediacn' => WemediacnHandlers::class,
    ];

    private $enterpriseId = '';

    private $configInstance = '';

    private $instanceHandlers = null;

    public function __construct(int $enterpriseId)
    {
        $this->enterpriseId   = $enterpriseId;
        $this->configInstance = ConfigFactory::getInstance($enterpriseId);
    }

    /**
     * 获取发送对象
     *
     * @param string $sendType
     * @return SendHandlerInterface
     */
    public function getMessageHandlers(string $sendType = ''): SendHandlerInterface
    {
        if (!$sendType) {
            $sendType = $this->configInstance->getType();
        }

        if (!isset($this->handlers[$sendType])) {
            throw new \InvalidArgumentException("sendType {$sendType} 不存在");
        }

        $config = $this->getConfig();

        if (!$config) {
            throw new \Exception("config {$config} 未设置");
        }

        $this->setInstanceHandlers(new $this->handlers[$sendType]($config));

        return $this->getInstanceHandlers();
    }

    /**
     * 增加新的消息类型
     *
     * @param array $handlers
     */
    public function registerHandlers(array $handlers)
    {
        $this->handlers = array_merge($this->handlers, $handlers);
    }

    /**
     * 获取配置信息
     * @return mixed
     */
    public function getConfig()
    {
        return $this->configInstance->getConfig();
    }

    /**
     * 获取当前短信的类型
     * @return mixed
     */
    public function getType()
    {
        return $this->configInstance->getType();
    }

    /**
     * 获取短信handler的实例
     */
    public function getInstanceHandlers()
    {
        return $this->instanceHandlers;
    }

    /**
     * 设置短信handler的实例
     * @param $instance
     */
    public function setInstanceHandlers($instance)
    {
        $this->instanceHandlers = $instance;
    }
}
