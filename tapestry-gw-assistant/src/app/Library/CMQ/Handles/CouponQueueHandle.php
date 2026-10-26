<?php

namespace App\Library\CMQ\Handles;

/*
 * CMQ_V1.0.2
 *
 *  1 Account类对象不是线程安全的，如果多线程使用，需要每个线程单独初始化Account类对象
 *  2 Topic与Queue使用不同的endpoint, 因此需要需要分别初始化Account
 *  3 创建订阅的时候，需要设置订阅的属性，订阅属性参见SubscriptionMeta的定义
 */

/**
 * 使用CMQ将用户的优惠券信息推送到CMQ队列
 */

use App\Library\CMQ\Account;
use App\Library\CMQ\Message;
use App\Library\CMQ\CmqQueue;
use Illuminate\Support\Facades\Log;
use App\Library\CMQ\ClientException\CMQExceptionBase;

class CouponQueueHandle
{
    /**
     * 接收的参数
     * @var
     */
    private $params;

    /**
     * 服务
     * @var CmqQueue
     */
    private $my_queue;

    private $endpoint;

    private $secretId;

    private $secretKey;

    private $queue_name;

    public function __construct($params)
    {
        $this->params = $params;
        $this->secretId = config("cmq.secretId");
        $this->endpoint = config("cmq.endPoint");
        $this->secretKey = config("cmq.secretKey");
        $this->queue_name = config('app.env') . "-shopping-guide-coupon";
    }

    /**
     * 对外接口
     */
    public function run()
    {
        Log::info("优惠券发送信息:开始:" . $this->queue_name, $this->params);

        $my_account = new Account($this->endpoint, $this->secretId, $this->secretKey);
        $this->my_queue = $my_account->get_queue($this->queue_name);

        $isSuccess = $this->send_message();

        if ($isSuccess) {
            Log::info("优惠券发送信息:结束:" . $this->queue_name);
        }
    }

    /**
     * 发送信息
     */
    public function send_message()
    {
        try {
            $msg_body = json_encode($this->params, JSON_UNESCAPED_UNICODE);
            $msg_content = new Message($msg_body);
            $this->my_queue->send_message($msg_content);
        } catch (CMQExceptionBase $exception) {

            if (in_array($exception->getCode(), [4000, 4440])) {
                Log::info("优惠券发送信息:异常:" . $this->queue_name, [
                    'code'       => $exception->getCode(),
                    'file'       => $exception->getFile(),
                    'line'       => $exception->getLine(),
                    'params'     => $this->params,
                    'message'    => $exception->getMessage(),
                    'queue_name' => $this->queue_name,
                ]);
            }

            return false;
        }

        return true;
    }
}
