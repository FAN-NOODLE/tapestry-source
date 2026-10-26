<?php

namespace App\Library\Cmq;

use App\Library\Cmq\Api\Account;
use App\Library\Cmq\Api\Message;
use Illuminate\Support\Facades\Log;

/**
 * CMQ QUEUE
 *
 * Class CmqQueue
 * @package App\Library\Cmq
 */
class CmqQueue
{
    private $queue;

    public function __construct($queueName)
    {
        $secretId  = env('SECRET_ID');
        $secretKey = env('SECRET_KEY');
        $endPoint  = env('QUEUE_END_POINT');

        $account = new Account($endPoint, $secretId, $secretKey);
        if (!empty($queueName)) {
            $this->queue = $account->get_queue($queueName);
        }
    }

    /**
     * 发送消息通知
     *
     * @param $message
     *
     * @return string
     * @author: zhaojie
     * @create: 2019/10/25 21:19
     */
    public function sendMessage($message)
    {

        try {
            // 消息发送
            $msg    = new Message($message);
            $result = $this->queue->send_message($msg);
            return $result;
        } catch (\Exception $e) {
            Log::error('cmq send message error:' . $e->getMessage());
        }
        return null;
    }

    /**
     * 消费消息
     * @return string
     * @author: zhaojie
     * @create: 2019/10/25 21:19
     */
    public function receiveMessage()
    {
        try {
            $result = $this->queue->receive_message(3);
            return $result;
        } catch (\Exception $e) {
            Log::error('cmq receive message error:' . $e->getMessage());
        }
        return null;
    }


    /**
     * 删除消息
     *
     * @param $receiptHandle
     *
     * @author: zhaojie
     * @create: 2019/10/26 12:00
     */
    public function deleteMessage($receiptHandle)
    {
        try {
            $this->queue->delete_message($receiptHandle);
        } catch (\Exception $e) {
            Log::error('cmq delete message error:' . $e->getMessage());
        }
    }
}
