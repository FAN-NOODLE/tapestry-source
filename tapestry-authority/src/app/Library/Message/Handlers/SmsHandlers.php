<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-04-14
 * Time: 17:05
 */

namespace App\Library\Message\Handlers;

use App\Library\Message\Contracts\SendHandlerInterface;
use Overtrue\EasySms\EasySms;

class SmsHandlers implements SendHandlerInterface
{
    /**
     * 配置信息
     *
     * @var array
     */
    private $config = [];

    /**
     * SmsHandlers constructor.
     * @param array $config
     */
    public function __construct(array $config)
    {
        $this->config = $config;
    }

    /**
     * 检测参数
     *
     * @param array $messageData
     * @throws \Exception
     */
    protected function checkRequestParams(array $messageData): void
    {
        if (!isset($messageData['subject']) || empty($messageData['subject'])) {
            throw new \Exception('subject is empty');
        }

        if (!isset($messageData['templateCode']) || empty($messageData['templateCode'])) {
            throw new \Exception('templateCode is empty or not fund');
        }

        if (isset($messageData['options'])) {
            $messageData['gateways'] = $messageData['options'];
            unset($messageData['options']);
        }
    }

    /**
     * 发送消息
     *
     * @param array $params
     * @return bool
     * @throws \Exception
     */
    public function send(array $params): bool
    {
        try {
            $this->checkRequestParams($params);
            (new EasySms($this->config))->send($params['subject'], [
                'template' => $params['templateCode'],
                'data'     => $params['data']
            ], $params['gateways'] ?? []);
        } catch (\Exception $exception) {
            throw $exception;
        }
        return true;
    }

    public function getAccount(): string
    {
        return md5($this->config['gateways']['sendcloud']['sms_user']);
    }
}
