<?php

namespace App\Library\Cmq;

use App\Library\Cmq\Api\Account;
use App\Library\Cmq\Api\CMQExceptionBase;
use App\Library\Cmq\Api\SubscriptionMeta;
use App\Library\Cmq\Api\TopicMeta;

/**
 * CMQ PHP Demo
 *
 * 1、Account类对象不是线程安全的，如果多线程使用，需要每个线程单独初始化Account类对象；
 * 2、Topic与Queue使用不同的endpoint, 因此需要需要分别初始化Account；
 * 3、创建订阅的时候，需要设置订阅的属性，订阅属性参见SubscriptionMeta的定义。
 *
 * Class TopicDemo
 * @package App\Library\Cmq
 */
class TopicDemo
{
    private $secretId;
    private $secretKey;
    private $endPoint;

    public function __construct($secretId, $secretKey, $endPoint)
    {
        $this->secretId  = $secretId;
        $this->secretKey = $secretKey;
        $this->endPoint  = $endPoint;
    }

    public function run()
    {
        try {
            // create account and topic
            $topicName = "Topic-Test";
            $myAccount = new Account($this->endPoint, $this->secretId, $this->secretKey);
            $myTopic   = $myAccount->get_topic($topicName);
            // $myTopicMeta = new TopicMeta();
            // $myTopic->create($myTopicMeta);

            $myTopicMeta             = $myTopic->get_attributes();
            $myTopicMeta->maxMsgSize = 1024;
            $myTopic->set_attributes($myTopicMeta);
            // list topic
            $topicList = $myAccount->list_topic();

            // publish message and batch publish message without tags
            $msg   = "666666";
            $msgId = $myTopic->publish_message($msg);
            echo "publish message without tag \n";

            // create subscription
            $subscriptionName = "subscription-test";
            $mySubscription   = $myAccount->get_subscription($topicName, $subscriptionName);
            $subscriptionMeta = new SubscriptionMeta();

            // get and set subscription meta
            // please input your endpoint and protocol
            $subscriptionMeta->Endpoint = "http://distribution-api-dev-crs.vchangyi.com/callback/test";
            $subscriptionMeta->Protocol = "HTTP";
            $mySubscription->create($subscriptionMeta);
            echo "create sub \n";
            $subscriptionMeta = $mySubscription->get_attributes();
            echo "get attributes\n";
            echo $subscriptionMeta;
            $mySubscription->set_attributes($subscriptionMeta);
            echo "set attributes\n";

            // list subscription
            $subscriptionList = $myTopic->list_subscription($topicName);
            echo $subscriptionList;
            echo "list subscription \n";

            // delete subscription and topic
            $mySubscription->delete();
            echo "delete subscription \n";
            $myTopic->delete();
            echo "delete topic \n";
        } catch (CMQExceptionBase $e) {
            echo $e;
        }
    }
}

// 从腾讯云官网查看云api的密钥信息
$secretId  = "REDACTED_FOR_SOURCE_SHARING";
$secretKey = "REDACTED_FOR_SOURCE_SHARING";
$endPoint  = "https://cmq-topic-sh.api.qcloud.com";

$instance = new TopicDemo($secretId, $secretKey, $endPoint);
$instance->run();
