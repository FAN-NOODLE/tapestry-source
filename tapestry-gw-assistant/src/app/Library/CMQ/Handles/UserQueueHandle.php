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
 * 使用CMQ将用户信息推送到CMQ队列
 */

use App\Exceptions\Error;
use App\Library\CMQ\Account;
use App\Library\CMQ\ClientException\CMQExceptionBase;
use App\Library\CMQ\CmqQueue;
use App\Library\CMQ\Message;
use App\Library\CMQ\QueueMeta;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class UserQueueHandle
{
    private $secretId;
    private $secretKey;
    private $endpoint;
    private $param; //请求参数
    /**
     * @var CmqQueue
     */
    private $my_queue; //队列对象
    private $queue_name; //队列名称

    public function __construct($secretId, $secretKey, $endpoint, $queue_name, $param)
    {
        $this->secretId  = $secretId;
        $this->secretKey = $secretKey;
        $this->endpoint  = $endpoint;
        $this->param     = $this->validateParam($param);
        $this->queue_name = $queue_name;
    }

    public function run()
    {

        Log::info("CMQ服务用户队列：".$this->queue_name.",消息推送开始");
        $my_account       = new Account($this->endpoint, $this->secretId, $this->secretKey);

        $this->my_queue = $my_account->get_queue($this->queue_name);

        $queue_meta                      = new QueueMeta();
        $queue_meta->queueName           = $this->queue_name;
        $queue_meta->pollingWaitSeconds  = config("cmq.pollingWaitSeconds"); //消息接收长轮询等待时间。取值范围200ms - 30s，默认值200ms。
        $queue_meta->visibilityTimeout   = config("cmq.visibilityTimeout"); //消息可见性超时。取值范围1 - 43200秒（即12小时内），默认值30。
        $queue_meta->maxMsgSize          = config("cmq.maxMsgSize"); //消息最大长度(单位：Byte)
        $queue_meta->msgRetentionSeconds = config("cmq.msgRetentionSeconds");  //消息生命周期(4天)

        try {
            switch ($this->param['action']) {
                //单条发送队列消息
                case 'send_message':
                    $this->send_message($queue_meta);
                    break;
                //单条消费队列消息
                case 'receive_message':
                    $this->receive_message();
                    break;
                //批量发送队列消息
                case 'batch_send_message':
                    break;
                //批量消费队列消息
                case 'batch_receive_message':
                    break;
                //删除队列
                case 'delete':
                    break;
                default:
                    ;
            }

        } catch (CMQExceptionBase $exception) {
            Log::info("CMQ服务用户队列：".$this->queue_name.",消息推送结束");
            return false;
        }
        Log::info("CMQ服务用户队列：".$this->queue_name.",消息推送结束");
    }

    /**
     * 校验参数
     *
     * @param $param
     *
     * @return mixed
     */
    private function validateParam($param)
    {
        $rules    = [
            'action' => 'required',
        ];
        $messages = [
            'action.required' => 'CMQ队列操作名称不能为空',
        ];

        $validator = Validator::make($param, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(codeNum('system_error'), $validator->errors()->first(), 210);
        }

        return $param;
    }

    /**
     * 检查队列是否存在
     */
    private function check_queue()
    {
        //获取队列属性
        try {
            $queue_info = $this->my_queue->get_attributes();
        } catch (CMQExceptionBase $exception) {
            if (in_array($exception->getCode(),[4000,4440])) {
                $user_id = $this->param['userInfo']['id'];
                $action   = $this->param['action'];
                Log::info("CMQ用户服务check_queue报错：", [
                    'queue_name' => $this->queue_name,
                    'user_id'   => $user_id,
                    'action'     => $action,
                    'msg'        => $exception->getMessage(),
                    'code'       => $exception->getCode(),
                    'file'       => $exception->getFile(),
                    'line'       => $exception->getLine(),
                ]);
                $queue_info = [];
            }
        }


        if (!empty($queue_info)) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * 发送队列消息
     *
     * @param $queue_meta
     *
     * @return bool
     */
    public function send_message($queue_meta)
    {
        try {
            $msg_body = json_encode($this->param, JSON_UNESCAPED_UNICODE);
            $msg      = new Message($msg_body);
            $re_msg   = $this->my_queue->send_message($msg);
            Log::info("消息发送到CMQ成功,MessageID:" . $re_msg->msgId . ",MessageBody:" . $msg_body);
        } catch (CMQExceptionBase $exception) {
            if (in_array($exception->getCode(),[4000,4440])) {
                $user_id = $this->param['userInfo']['id'];
                $action   = $this->param['action'];
                Log::info("CMQ用户send_message报错：", [
                    'queue_name' => $this->queue_name,
                    'user_id'   => $user_id,
                    'action'     => $action,
                    'msg'        => $exception->getMessage(),
                    'code'       => $exception->getCode(),
                    'file'       => $exception->getFile(),
                    'line'       => $exception->getLine(),
                ]);
            }
        }
    }

    /**
     * 消费队列消息
     */
    public function receive_message()
    {
        $recv_msg = $this->my_queue->receive_message(3);
        Log::info("消费CMQ队列消息成功" . $recv_msg);
        $this->my_queue->delete_message($recv_msg->respInter);
    }

    /**
     * 删除队列
     */
    public function delete()
    {
        $this->my_queue->delete();
    }
}
