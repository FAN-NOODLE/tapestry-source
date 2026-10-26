<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:01
 */

namespace App\Library\WechatNotice;

use App\Jobs\AlertMsgJob;
use App\Library\WechatNotice\Format\FormatInterface;
use App\Library\WechatNotice\Format\MarkDownFormat;


class Chat
{
    private static $instance;

    /**
     * @var string
     */
    const API_URL = "https://qyapi.weixin.qq.com";

    /**
     * @var string
     */
    private $web_hook_url;

    /**
     * @var string
     */
    private $webHookKey = '';

    /**
     * @var FormatInterface
     */
    private $format;

    private function __construct()
    {
        $this->setDefaultWebHookKey();
        $this->format = new MarkDownFormat();
    }

    private function setDefaultWebHookKey()
    {
        switch (config('app.env')) {
            case 'test':
            case 'uat':
                $hookKey = "58dac20b-979d-46f0-b10d-96a8bcbd1645";
                break;
            case 'prod' :
                $hookKey = '4cbe8c17-8100-48ab-a911-8f2bdf773edf';
                break;
            default :
                $hookKey = "17c80e16-5719-4312-9471-1d68fc15c47c";
                break;
        }

        $this->setHook($hookKey);
    }

    /**
     * 单例
     *
     * @return static
     */
    public static function getInstance()
    {
        if (!static::$instance) {
            static::$instance = new static;
        }
        return static::$instance;
    }

    public function setHook($hookKey)
    {
        $this->webHookKey   = $hookKey;
        $this->web_hook_url = sprintf("/cgi-bin/webhook/send?key=%s", $hookKey);
        return $this;
    }


    /**
     * @param string $msg
     * @param $data
     * @param bool|null $sync
     * @return bool
     * @throws \Exception
     */
    public function send(string $msg, ?array $data = [], ?bool $sync = false)
    {

        if (!$sync) {
            //异步发送
            dispatch(new AlertMsgJob($msg, $data, $this->webHookKey));
            return true;
        }
        $content = "### " . $msg . "\n";
        if ($data) {
            $data = var_export((array)$data, true);
        } else {
            $data = '';
        }

        $content = $this->format->format($content . $data);

        $response = Request::requestPostJson(self::API_URL, $this->web_hook_url, $content);
        $response = json_decode($response, true);
        if ($response['errcode'] !== 0) {
            throw new SendWorkWechatMsgException((int)$response['errcode'], $response['errmsg']);
        }
        return true;
    }

}