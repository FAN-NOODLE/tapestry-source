<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-14
 * Time: 17:08
 */

namespace App\Library\Message;

use App\Library\Message\Contracts\SendHandlerInterface;
use App\Library\Message\Handlers\MiniSubscriptionHandlers;
use App\Library\Message\Handlers\SmsHandlers;
use App\Library\Message\Handlers\WorkMessageHandlers;

class FactoryMessage extends MessageBase
{
    /**
     * 默认可以发送的类型
     *
     * @var array
     */
    private $handlers = [
        'easySms'          => SmsHandlers::class,
        'workMessage'      => WorkMessageHandlers::class,
        'miniSubscription' => MiniSubscriptionHandlers::class
    ];

    /**
     * 获取发送对象
     *
     * @param string $sendType
     * @return SendHandlerInterface
     */
    public function getMessageHandlers(string $sendType): SendHandlerInterface
    {
        if (!isset($this->handlers[$sendType])) {
            throw new \InvalidArgumentException("sendType {$sendType} 不存在");
        }
        $configName = strtolower(preg_replace('/(?<=[a-z])([A-Z])/', '-$1', $sendType));
        $config = config($configName);
        return new $this->handlers[$sendType]($config);
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
}
