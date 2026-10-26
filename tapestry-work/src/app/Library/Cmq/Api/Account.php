<?php

namespace App\Library\Cmq\Api;

/**
 * Account 对象非线程安全，如果多线程使用，需要每个线程单独初始化Account对象类
 *
 * Class Account
 * @package App\Library\Cmq\Api
 */
class Account
{
    private $host;
    private $secretId;
    private $secretKey;
    private $cmq_client;

    /**
     * Account constructor.
     *
     * @param string $host      访问的url，例如：https://cmq-queue-gz.api.qcloud.com
     * @param string $secretId  用户的secretId, 腾讯云官网获取
     * @param string $secretKey 用户的secretKey，腾讯云官网获取
     */
    public function __construct($host, $secretId, $secretKey)
    {
        $this->host       = $host;
        $this->secretId   = $secretId;
        $this->secretKey  = $secretKey;
        $this->cmq_client = new CMQClient($host, $secretId, $secretKey);
    }

    /**
     * @param string $sign_method only support sha1 and sha256
     *
     * @throws CMQClientParameterException
     */
    public function set_sign_method($sign_method = 'sha1')
    {
        $this->cmq_client->set_sign_method($sign_method);
    }

    /**
     * 设置访问的url
     *
     * @param string $host      访问的url，例如：http://cmq-queue-gz.api.tencentyun.com
     * @param string $secretId  用户的secretId，腾讯云官网获取
     * @param string $secretKey 用户的secretKey，腾讯云官网获取
     */
    public function set_client($host, $secretId = NULL, $secretKey = NULL)
    {
        if ($secretId == NULL) {
            $secretId = $this->secretId;
        }
        if ($secretKey == NULL) {
            $secretKey = $this->secretKey;
        }
        $this->cmq_client = new CMQClient($host, $secretId, $secretKey);
    }

    /**
     * 获取queue client
     *
     * @return object CMQClient
     */
    public function get_client()
    {
        return $this->cmq_client;
    }

    /**
     * 获取Account的一个Queue对象
     *
     * @param string $queue_name 队列名
     *
     * @return object Queue 返回该Account的一个Queue对象
     */
    public function get_queue($queue_name)
    {
        return new Queue($queue_name, $this->cmq_client);
    }

    /**
     * 列出Account的队列
     *
     * @param string $searchWord 队列名的前缀
     * @param int    $limit      list_queue最多返回的队列数
     * @param string $offset     list_queue的起始位置，上次list_queue返回的next_offset
     *
     * @return array QueueURL的列表和下次list queue的起始位置；如果所有queue都list出来，next_offset为""
     */
    public function list_queue($searchWord = "", $limit = -1, $offset = "")
    {
        $params = array();
        if ($searchWord != "") {
            $params['searchWord'] = $searchWord;
        }
        if ($limit != -1) {
            $params['limit'] = $limit;
        }
        if ($offset != "") {
            $params['offset'] = $offset;
        }

        $ret_pkg = $this->cmq_client->list_queue($params);

        if ($offset == "") {
            $next_offset = count($ret_pkg['queueList']);
        } else {
            $next_offset = $offset + count($ret_pkg['queueList']);
        }
        if ($next_offset >= $ret_pkg['totalCount']) {
            $next_offset = "";
        }

        return array("totalCount"  => $ret_pkg['totalCount'],
                     "queueList"   => $ret_pkg['queueList'],
                     "next_offset" => $next_offset);
    }

    /**
     * 列出Account的主题
     *
     * @param string $searchWord 主题关键字
     * @param int    $limit      最多返回的主题数目
     * @param string $offset     list_topic的起始位置，上次list_topic返回的next_offset
     *
     * @return array TopicURL的列表和下次list topic的起始位置；如果所有topic都list出来，next_offset为""
     */
    public function list_topic($searchWord = "", $limit = -1, $offset = "")
    {
        $params = array();
        if ($searchWord != "") {
            $params['searchWord'] = $searchWord;
        }
        if ($limit != -1) {
            $params['limit'] = $limit;
        }
        if ($offset != "") {
            $params['offset'] = $offset;
        }

        $resp = $this->cmq_client->list_topic($params);

        if ($offset == "") {
            $next_offset = count($resp['topicList']);
        } else {
            $next_offset = $offset + count($resp['topicList']);
        }
        if ($next_offset >= $resp['totalCount']) {
            $next_offset = "";
        }

        return array("totalCount"  => $resp['totalCount'],
                     "topicList"   => $resp['topicList'],
                     "next_offset" => $next_offset);
    }

    /**
     * 获取Account的一个Topic对象
     *
     * @param string $topic_name
     *
     * @return object Topic 返回该Account的一个Topic对象
     */
    public function get_topic($topic_name)
    {
        return new Topic($topic_name, $this->cmq_client);
    }

    /**
     * 获取Account的一个Subscription对象
     *
     * @param string $topic_name
     * @param string $subscription_name
     *
     * @return object Subscription 返回该Account的一个Subscription对象
     */
    public function get_subscription($topic_name, $subscription_name)
    {
        return new Subscription($topic_name, $subscription_name, $this->cmq_client);
    }
}
