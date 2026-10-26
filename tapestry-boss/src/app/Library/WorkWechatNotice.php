<?php

namespace App\Library;

use Exception;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;

/**
 * Class 企业微信提醒
 *
 * @package App\Library
 */
class WorkWechatNotice
{
    /**
     * web hook地址
     */
    const WEB_HOOK_URL = 'https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s';

    /**
     * 执行环境
     */
    const HANDLE_ENVS = ['prod', 'test'];

    /**
     * 环境
     *
     * @var string
     */
    private $env;

    public function __construct(string $env)
    {
        $this->env = $env;
    }

    /**
     * 发送消息
     *
     * @param string $keywords
     * @param array $context
     */
    public function send(string $keywords, array $context = [])
    {
        if (in_array($this->env, self::HANDLE_ENVS)) {

            $contentStr = '注意：' . $keywords . PHP_EOL;
            foreach ($context as $key => $value) {
                $contentStr .= '> <font color="#3399FF">' . $key . '：</font><font color="#666666">' . (is_array($value) ? Json::encode($value) : $value) . '</font>' . PHP_EOL;
            }
            try {
                (new Client())->request('POST', sprintf(self::WEB_HOOK_URL, env("NOTICE_KEY", "58dac20b-979d-46f0-b10d-96a8bcbd1645")), [
                    'http_errors' => false,
                    'headers'     => [
                        'Content-Type' => 'application/json'
                    ],
                    'body'        => Json::encode([
                        'msgtype'  => 'markdown',
                        'markdown' => [
                            'content' => $contentStr
                        ]
                    ])
                ]);
            } catch (GuzzleException $exception) {
            } catch (Exception $exception) {
            }
        }
    }
}
