<?php

namespace App\Library\Cmq\Api;

class Queue
{
    private $queue_name;
    private $cmq_client;
    private $encoding;

    public function __construct($queue_name, $cmq_client, $encoding = false)
    {
        $this->queue_name = $queue_name;
        $this->cmq_client = $cmq_client;
        $this->encoding   = $encoding;
    }

    /**
     * 设置是否对消息体进行base64编码
     *
     * @param bool $encoding 是否对消息体进行base64编码
     */
    public function set_encoding($encoding)
    {
        $this->encoding = $encoding;
    }

    /**
     * 创建队列
     *
     * @param object $queue_meta QueueMeta对象，设置队列的属性
     */
    public function create($queue_meta)
    {
        $params = array(
            'queueName'           => $this->queue_name,
            'pollingWaitSeconds'  => $queue_meta->pollingWaitSeconds,
            'visibilityTimeout'   => $queue_meta->visibilityTimeout,
            'maxMsgSize'          => $queue_meta->maxMsgSize,
            'msgRetentionSeconds' => $queue_meta->msgRetentionSeconds,
            'rewindSeconds'       => $queue_meta->rewindSeconds,
        );
        if ($queue_meta->maxMsgHeapNum > 0) {
            $params['maxMsgHeapNum'] = $queue_meta->maxMsgHeapNum;
        }
        $this->cmq_client->create_queue($params);
    }

    /**
     * 获取队列属性
     *
     * @return object QueueMeta 队列的属性
     */
    public function get_attributes()
    {
        $params                = array(
            'queueName' => $this->queue_name
        );
        $resp                  = $this->cmq_client->get_queue_attributes($params);
        $queue_meta            = new QueueMeta();
        $queue_meta->queueName = $this->queue_name;
        $this->__resp2meta__($queue_meta, $resp);
        return $queue_meta;
    }

    /**
     * 设置队列属性
     *
     * @param object $queue_meta QueueMeta对象，设置队列的属性
     */
    public function set_attributes($queue_meta)
    {
        $params = array(
            'queueName'           => $this->queue_name,
            'pollingWaitSeconds'  => $queue_meta->pollingWaitSeconds,
            'visibilityTimeout'   => $queue_meta->visibilityTimeout,
            'maxMsgSize'          => $queue_meta->maxMsgSize,
            'msgRetentionSeconds' => $queue_meta->msgRetentionSeconds,
            'rewindSeconds'       => $queue_meta->rewindSeconds
        );
        if ($queue_meta->maxMsgHeapNum > 0) {
            $params['maxMsgHeapNum'] = $queue_meta->maxMsgHeapNum;
        }

        $this->cmq_client->set_queue_attributes($params);
    }

    public function rewindQueue($backTrackingTime)
    {
        $params = array(
            'queueName'        => $this->queue_name,
            'startConsumeTime' => $backTrackingTime
        );
        $this->cmq_client->rewindQueue($params);
    }

    /**
     * 删除队列
     */
    public function delete()
    {
        $params = array('queueName' => $this->queue_name);
        $this->cmq_client->delete_queue($params);
    }

    /**
     * 发送消息
     *
     * @param object $message 发送的Message object
     * @param int    $delayTime
     *
     * @return object Message 消息发送成功的返回属性，包含MessageId
     */
    public function send_message($message, $delayTime = 0)
    {
        if ($this->encoding) {
            $msgBody = base64_encode($message->msgBody);
        } else {
            $msgBody = $message->msgBody;
        }
        $params        = array(
            'queueName'    => $this->queue_name,
            'msgBody'      => $msgBody,
            'delaySeconds' => $delayTime
        );
        $msgId         = $this->cmq_client->send_message($params);
        $retmsg        = new Message();
        $retmsg->msgId = $msgId;
        return $retmsg;
    }

    /**
     * 批量发送消息
     *
     * @param object $messages 发送的Message object list
     * @param int    $delayTime
     *
     * @return array 多条消息发送成功的返回属性，包含MessageId
     */
    public function batch_send_message($messages, $delayTime = 0)
    {
        $params = array(
            'queueName'    => $this->queue_name,
            'delaySeconds' => $delayTime
        );
        $n      = 1;
        foreach ($messages as $message) {
            $key = 'msgBody.' . $n;
            if ($this->encoding) {
                $params[$key] = base64_encode($message->msgBody);
            } else {
                $params[$key] = $message->msgBody;
            }
            $n += 1;
        }
        $msgList        = $this->cmq_client->batch_send_message($params);
        $retMessageList = array();
        foreach ($msgList as $msg) {
            $retmsg            = new Message();
            $retmsg->msgId     = $msg['msgId'];
            $retMessageList [] = $retmsg;
        }
        return $retMessageList;
    }

    /**
     * 消费消息
     *
     * @param int $polling_wait_seconds 本次请求的长轮询时间，单位：秒
     *
     * @return object Message object中包含基本属性、临时句柄
     */
    public function receive_message($polling_wait_seconds = NULL)
    {
        $params = array('queueName' => $this->queue_name);
        if ($polling_wait_seconds != NULL) {
            $params['UserpollingWaitSeconds'] = $polling_wait_seconds;
            $params['pollingWaitSeconds']     = $polling_wait_seconds;
        } else {
            $params['UserpollingWaitSeconds'] = 30;
        }
        $resp = $this->cmq_client->receive_message($params);
        $msg  = new Message();
        if ($this->encoding) {
            $msg->msgBody = base64_decode($resp['msgBody']);
        } else {
            $msg->msgBody = $resp['msgBody'];
        }
        $msg->msgId            = $resp['msgId'];
        $msg->receiptHandle    = $resp['receiptHandle'];
        $msg->enqueueTime      = $resp['enqueueTime'];
        $msg->nextVisibleTime  = $resp['nextVisibleTime'];
        $msg->dequeueCount     = $resp['dequeueCount'];
        $msg->firstDequeueTime = $resp['firstDequeueTime'];
        return $msg;
    }

    /**
     * 批量消费消息
     *
     * @param int $num_of_msg           本次请求最多获取的消息条数
     * @param int $polling_wait_seconds 本次请求的长轮询时间，单位：秒
     *
     * @return array 多条消息的属性，包含消息的基本属性、临时句柄
     */
    public function batch_receive_message($num_of_msg, $polling_wait_seconds = NULL)
    {
        $params = array('queueName' => $this->queue_name, 'numOfMsg' => $num_of_msg);
        if ($polling_wait_seconds != NULL) {
            $params['UserpollingWaitSeconds'] = $polling_wait_seconds;
            $params['pollingWaitSeconds']     = $polling_wait_seconds;
        } else {
            $params['UserpollingWaitSeconds'] = 30;
        }
        $msgInfoList    = $this->cmq_client->batch_receive_message($params);
        $retMessageList = array();
        foreach ($msgInfoList as $msg) {
            $retmsg = new Message();
            if ($this->encoding) {
                $retmsg->msgBody = base64_decode($msg['msgBody']);
            } else {
                $retmsg->msgBody = $msg['msgBody'];
            }
            $retmsg->msgId            = $msg['msgId'];
            $retmsg->receiptHandle    = $msg['receiptHandle'];
            $retmsg->enqueueTime      = $msg['enqueueTime'];
            $retmsg->nextVisibleTime  = $msg['nextVisibleTime'];
            $retmsg->dequeueCount     = $msg['dequeueCount'];
            $retmsg->firstDequeueTime = $msg['firstDequeueTime'];
            $retMessageList []        = $retmsg;
        }
        return $retMessageList;
    }

    /**
     * 删除消息
     *
     * @param string $receipt_handle 最近一次操作该消息返回的临时句柄
     */
    public function delete_message($receipt_handle)
    {
        $params = array('queueName' => $this->queue_name, 'receiptHandle' => $receipt_handle);
        $this->cmq_client->delete_message($params);
    }

    /**
     * 批量删除消息
     *
     * @param array $receipt_handle_list batch_receive_message返回的多条消息的临时句柄
     */
    public function batch_delete_message($receipt_handle_list)
    {
        $params = array('queueName' => $this->queue_name);
        $n      = 1;
        foreach ($receipt_handle_list as $receipt_handle) {
            $key          = 'receiptHandle.' . $n;
            $params[$key] = $receipt_handle;
            $n            += 1;
        }
        $this->cmq_client->batch_delete_message($params);
    }

    protected function __resp2meta__($queue_meta, $resp)
    {
        if (isset($resp['queueName'])) {
            $queue_meta->queueName = $resp['queueName'];
        }
        if (isset($resp['maxMsgHeapNum'])) {
            $queue_meta->maxMsgHeapNum = $resp['maxMsgHeapNum'];
        }
        if (isset($resp['pollingWaitSeconds'])) {
            $queue_meta->pollingWaitSeconds = $resp['pollingWaitSeconds'];
        }
        if (isset($resp['visibilityTimeout'])) {
            $queue_meta->visibilityTimeout = $resp['visibilityTimeout'];
        }
        if (isset($resp['maxMsgSize'])) {
            $queue_meta->maxMsgSize = $resp['maxMsgSize'];
        }
        if (isset($resp['msgRetentionSeconds'])) {
            $queue_meta->msgRetentionSeconds = $resp['msgRetentionSeconds'];
        }
        if (isset($resp['createTime'])) {
            $queue_meta->createTime = $resp['createTime'];
        }
        if (isset($resp['lastModifyTime'])) {
            $queue_meta->lastModifyTime = $resp['lastModifyTime'];
        }
        if (isset($resp['activeMsgNum'])) {
            $queue_meta->activeMsgNum = $resp['activeMsgNum'];
        }
        if (isset($resp['rewindSeconds'])) {
            $queue_meta->rewindSeconds = $resp['rewindSeconds'];
        }
        if (isset($resp['inactiveMsgNum'])) {
            $queue_meta->inactiveMsgNum = $resp['inactiveMsgNum'];
        }
        if (isset($resp['rewindmsgNum'])) {
            $queue_meta->rewindmsgNum = $resp['rewindmsgNum'];
        }
        if (isset($resp['minMsgTime'])) {
            $queue_meta->minMsgTime = $resp['minMsgTime'];
        }
        if (isset($resp['delayMsgNum'])) {
            $queue_meta->delayMsgNum = $resp['delayMsgNum'];
        }
    }
}

class QueueMeta
{
    public $queueName;
    public $maxMsgHeapNum;
    public $pollingWaitSeconds;
    public $visibilityTimeout;
    public $maxMsgSize;
    public $msgRetentionSeconds;
    public $createTime;
    public $lastModifyTime;
    public $activeMsgNum;
    public $inactiveMsgNum;
    public $rewindSeconds;
    public $rewindmsgNum;
    public $minMsgTime;
    public $delayMsgNum;

    /* 队列属性
        @note: 设置属性
        :: maxMsgHeapNum: 最大堆积消息数
        :: pollingWaitSeconds: receive message时，长轮询时间，单位：秒
        :: visibilityTimeout: 消息可见性超时, 单位：秒
        :: maxMsgSize: 消息最大长度, 单位：Byte
		:: msgRetentionSeconds: 消息保留周期，单位：秒
		:: rewindSeconds ： 最大回溯时间， 单位：秒

        @note: 非设置属性
        :: activeMsgNum: 可消费消息数，近似值
        :: inactiveMsgNum: 正在被消费的消息数，近似值
        :: createTime: queue创建时间，单位：秒
        :: lastModifyTime: 修改queue属性的最近时间，单位：秒
		:: queue_name: 队列名称
		:: rewindmsgNum:已删除，但是任然在回溯保留时间内的消息数量
		:: minMsgTime: 消息最小未消费时间，单位为秒
		:: delayMsgNum:延时消息数量
    */
    public function __construct()
    {
        $this->queueName           = "";
        $this->maxMsgHeapNum       = -1;
        $this->pollingWaitSeconds  = 0;
        $this->visibilityTimeout   = 30;
        $this->maxMsgSize          = 65536;
        $this->msgRetentionSeconds = 345600;
        $this->createTime          = -1;
        $this->lastModifyTime      = -1;
        $this->activeMsgNum        = -1;
        $this->inactiveMsgNum      = -1;
        $this->rewindSeconds       = 0;
        $this->rewindmsgNum        = 0;
        $this->minMsgTime          = 0;
        $this->delayMsgNum         = 0;
    }

    public function __toString()
    {
        $info = array("visibilityTimeout"   => $this->visibilityTimeout,
                      "maxMsgHeapNum"       => $this->maxMsgHeapNum,
                      "maxMsgSize"          => $this->maxMsgSize,
                      "msgRetentionSeconds" => $this->msgRetentionSeconds,
                      "pollingWaitSeconds"  => $this->pollingWaitSeconds,
                      "activeMsgNum"        => $this->activeMsgNum,
                      "inactiveMsgNum"      => $this->inactiveMsgNum,
                      "createTime"          => date("Y-m-d H:i:s", $this->createTime),
                      "lastModifyTime"      => date("Y-m-d H:i:s", $this->lastModifyTime),
                      "QueueName"           => $this->queueName,
                      "rewindSeconds"       => $this->rewindSeconds,
                      "rewindmsgNum"        => $this->rewindmsgNum,
                      "minMsgTime"          => $this->minMsgTime,
                      "delayMsgNum"         => $this->delayMsgNum);
        return json_encode($info);
    }
}

class Message
{
    public $msgBody;
    public $msgId;
    public $enqueueTime;
    public $receiptHandle;
    public $nextVisibleTime;
    public $dequeueCount;
    public $firstDequeueTime;

    /* 消息属性
        @note: send_message 指定属性
        :: msgBody         消息体

        @note: send_message 返回属性
        :: msgId           消息编号

        @note: receive_message 返回属性，除基本属性外
        :: receiptHandle       下次删除或修改消息的临时句柄
        :: enqueueTime         消息入队时间
        :: nextVisibleTime     下次可被再次消费的时间
        :: dequeueCount        总共被消费的次数
        :: firstDequeueTime    第一次被消费的时间
    */
    public function __construct($message_body = "")
    {
        $this->msgBody          = $message_body;
        $this->msgId            = "";
        $this->enqueueTime      = -1;
        $this->receiptHandle    = "";
        $this->nextVisibleTime  = -1;
        $this->dequeueCount     = -1;
        $this->firstDequeueTime = -1;
    }

    public function __toString()
    {
        $info = array("msgBody"          => $this->msgBody,
                      "msgId"            => $this->msgId,
                      "enqueueTime"      => date("Y-m-d H:i:s", $this->enqueueTime),
                      "nextVisibleTime"  => date("Y-m-d H:i:s", $this->nextVisibleTime),
                      "firstDequeueTime" => date("Y-m-d H:i:s", $this->firstDequeueTime),
                      "dequeueCount"     => $this->dequeueCount,
                      "receiptHandle"    => $this->receiptHandle);
        return json_encode($info);
    }
}
