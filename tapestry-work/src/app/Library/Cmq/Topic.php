<?php

namespace App\Library\Cmq;

use App\Library\Cmq\Api\Account;
use App\Library\Cmq\Api\SubscriptionMeta;
use App\Library\Cmq\Api\TopicMeta;

/**
 * CMQ PHP Demo
 *
 * 1、Account类对象不是线程安全的，如果多线程使用，需要每个线程单独初始化Account类对象；
 * 2、Topic与Queue使用不同的endpoint, 因此需要需要分别初始化Account；
 * 3、创建订阅的时候，需要设置订阅的属性，订阅属性参见SubscriptionMeta的定义。
 *
 * Class Topic
 * @package App\Library\Cmq
 */
class Topic
{
    private $secretId;
    private $secretKey;
    private $endpoint;
    private $connect;

    public function __construct($secretId = null, $secretKey = null, $endpoint = null)
    {
        $this->secretId  = empty($secretId) ? env('SECRET_ID') : $secretId;
        $this->secretKey = empty($secretKey) ? env('SECRET_KEY') : $secretKey;
        $this->endpoint  = empty($endpoint) ? env('TOPIC_END_POINT') : $endpoint;
        $this->connect   = new Account($this->endpoint, $this->secretId, $this->secretKey);
    }

    /**
     * 创建 topic
     *
     * @param       $topicName topic 名称
     * @param array $topicMeta 配置参数
     *
     * @return bool|void
     */
    public function create($topicName, $topicMeta = [])
    {
        // 已存在 topic
        if ($this->get_topic_attributes($topicName)) {
            return true;
        }

        $topic = $this->connect->get_topic($topicName);

        $mTopicMeta = new TopicMeta();
        $topic->create($mTopicMeta);
        if (!empty($topicMeta)) {
            foreach ($topicMeta as $k => $v) {
                $mTopicMeta->$k = $v;
            }
        }

        return $topic->set_attributes($mTopicMeta);
    }

    /**
     * 创建订阅
     *
     * @param string $topicName        topic 名称
     * @param string $subscriptionName 订阅名称
     * @param string $EndPoint         订阅的协议，目前支持两种协议：HTTP、queue。使用 HTTP 协议，用户需自己搭建接受消息的 Web Server。使用 queue，消息会自动推送到 CMQ queue，用户可以并发地拉取消息。
     * @param string $Protocol         接收通知的 endpoint，根据协议 protocol 区分：对于 HTTP，endpoint 必须以http://开头，host 可以是域名或 IP；对于 queue，则填 queueName。
     *
     * @return $this|bool
     */
    public function create_subscription($topicName, $subscriptionName, $EndPoint, $Protocol = 'HTTP')
    {
        if (empty($EndPoint)) {
            return false;
        }

        // 已存在订阅
        if ($this->get_subscription_attributes($topicName, $subscriptionName)) {
            return true;
        }

        $subscription     = $this->connect->get_subscription($topicName, $subscriptionName);
        $subscriptionMeta = new SubscriptionMeta();
        // get and set subscription meta
        // please input your endpoint and protocol
        $subscriptionMeta->Endpoint = $EndPoint;
        $subscriptionMeta->Protocol = $Protocol;
        return $subscription->create($subscriptionMeta);
    }

    /**
     * 发送消息
     *
     * @param $msg
     * @param $topicName
     *
     * @return mixed
     */
    public function publish_message($msg, $topicName)
    {
        $topic = $this->connect->get_topic($topicName);
        return $topic->publish_message($msg);
    }

    /**
     * 列出Account主题 topic 列表
     *
     * @param string $searchWord
     * @param int    $limit
     * @param string $offset
     *
     * @return array
     */
    public function list_topic($searchWord = '', $limit = -1, $offset = "")
    {
        return $this->connect->list_topic($searchWord, $limit, $offset);
    }

    /**
     * 删除订阅
     *
     * @param $topicName
     * @param $subscriptionName
     *
     * @return bool|void
     */
    public function subscription_delete($topicName, $subscriptionName)
    {
        // 是否存在订阅
        if (!$this->get_subscription_attributes($topicName, $subscriptionName)) {
            return false;
        }

        $mySubscription = $this->connect->get_subscription($topicName, $subscriptionName);
        return $mySubscription->delete();
    }

    /**
     * 删除topic，需要先删除订阅
     *
     * @param $topicName
     *
     * @return bool
     */
    public function topic_delete($topicName)
    {
        // 是否存在 topic
        if (!$this->get_topic_attributes($topicName)) {
            return false;
        }

        $topic = $this->connect->get_topic($topicName);
        $topic->delete();
    }

    /**
     * 订阅列表
     *
     * @param $topicName
     *
     * @return array
     */
    public function list_subscription($topicName)
    {
        $topic = $this->connect->get_topic($topicName);
        return $topic->list_subscription($topicName);
    }

    /**
     * topic 属性
     *
     * @param $topicName
     *
     * @return TopicMeta|bool
     */
    private function get_topic_attributes($topicName)
    {
        try {
            $topic = $this->connect->get_topic($topicName);
            return $topic->get_attributes();
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * 订阅属性
     *
     * @param $topicName
     * @param $subscriptionName
     *
     * @return SubscriptionMeta|bool
     */
    private function get_subscription_attributes($topicName, $subscriptionName)
    {
        try {
            $mySubscription = $this->connect->get_subscription($topicName, $subscriptionName);
            return $mySubscription->get_attributes();
        } catch (\Exception $e) {
            return false;
        }
    }
}
