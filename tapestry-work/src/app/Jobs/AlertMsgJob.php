<?php

namespace App\Jobs;


use App\Library\WechatNotice\Chat;

/**
 * Class SendMessageJob
 * @package App\Jobs
 */
class  AlertMsgJob extends Job
{
    /**
     * @var
     */
    private $msg;
    private $data;
    /**
     * @var string
     * 通知的企业微信机器人id
     */
    private $webHookKey;
    /**
     * 队列名称
     *
     * @var string
     */
    public $queue = 'alert_msg';

    /**
     * SendAlertMessageJob constructor.
     * @param $msg
     * @param $data
     * @param string $webHookKey
     */
    public function __construct(string $msg, array $data, string $webHookKey)
    {
        $this->data       = $data;
        $this->msg        = $msg;
        $this->webHookKey = $webHookKey;
    }


    public function handle()
    {
        return Chat::getInstance()
            ->setHook($this->webHookKey)
            ->send($this->msg, $this->data, true);
    }
}
