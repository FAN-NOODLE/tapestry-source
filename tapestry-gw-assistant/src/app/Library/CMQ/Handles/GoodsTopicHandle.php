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
 * 使用CMQ将商品信息推送到CMQ队列
 */

use App\Exceptions\Error;
use App\Library\CMQ\Account;
use App\Library\CMQ\ClientException\CMQExceptionBase;
use App\Library\CMQ\Topic;
use App\Library\CMQ\TopicMeta;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;

class GoodsTopicHandle
{
    private $secretId;
    private $secretKey;
    private $endpoint;
    private $param; //请求参数
    /**
     * @var Topic
     */
    private $my_topic; //主题对象
    private $topic_name; //主题名称

    public function __construct($secretId, $secretKey, $endpoint, $param)
    {
        $this->secretId  = $secretId;
        $this->secretKey = $secretKey;
        $this->endpoint  = $endpoint;
        $this->param     = $this->validateParam($param);
    }

    public function run()
    {
        $env              = config('app.env');

        $this->topic_name = $env . "-crs-goods-topic";
        Log::info("CMQ服务商品消息发布：" . $this->topic_name . ",发布消息开始");
        $my_account     = new Account($this->endpoint, $this->secretId, $this->secretKey);
        $this->my_topic = $my_account->get_topic($this->topic_name);
        $my_topic_meta   = new TopicMeta();

        try {
            switch ($this->param['action']) {
                //单条发布消息
                case 'send_message':
                    $this->send_message($my_topic_meta);
                    break;
                default:
                    ;
            }

        } catch (CMQExceptionBase $exception) {
            Log::info("CMQ服务商品消息发布：" . $this->topic_name . ",发布消息结束");
            return false;
        }
        Log::info("CMQ服务商品消息发布：" . $this->topic_name . ",发布消息结束");
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
            'action.required' => 'CMQ主题操作名称不能为空',
        ];

        $validator = Validator::make($param, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(codeNum('system_error'), $validator->errors()->first(), 210);
        }

        return $param;
    }

    /**
     * 检查主题是否存在
     */
    private function check_topic()
    {
        //获取队列属性
        try {
            $my_topic_meta  = $this->my_topic->get_attributes();
        } catch (CMQExceptionBase $exception) {
            if (in_array($exception->getCode(),[4000,4440])) {
                $good_id = $this->param['goodInfo']['id'];
                $action  = $this->param['action'];
                Log::info("CMQ商品GoodsTopicHandle服务check_topic报错：", [
                    'queue_name' => $this->topic_name,
                    'good_id'    => $good_id,
                    'action'     => $action,
                    'msg'        => $exception->getMessage(),
                    'code'       => $exception->getCode(),
                    'file'       => $exception->getFile(),
                    'line'       => $exception->getLine(),
                ]);
                $my_topic_meta = [];
            }
        }


        if (!empty($my_topic_meta)) {
            return true;
        } else {
            return false;
        }
    }

    /**
     * 发布主题消息
     *
     * @param $my_topic_meta
     */
    public function send_message($my_topic_meta)
    {
        $check_result = false;
        $res          = $this->check_topic();
        try {
            if ($res) {
                $check_result = true;
            } else {
                $this->my_topic->create($my_topic_meta);
                sleep(2);
                $my_topic_meta->maxMsgSize=65530;
                $this->my_topic->set_attributes($my_topic_meta);
                $res = $this->check_topic();
                if ($res) {
                    $check_result = true;
                }
            }
            if ($check_result) {
                //$msg   = json_encode($this->param, JSON_UNESCAPED_UNICODE);
                $msg = "我的测试消息";
                $this->my_topic->publish_message($msg);
            }
        } catch (CMQExceptionBase $exception) {
            if (in_array($exception->getCode(),[4000,4440])) {
                $good_id = $this->param['goodInfo']['id'];
                $action  = $this->param['action'];
                Log::info("CMQ商品GoodsTopicHandle服务send_message报错：", [
                    'queue_name' => $this->my_topic,
                    'good_id'    => $good_id,
                    'action'     => $action,
                    'msg'        => $exception->getMessage(),
                    'code'       => $exception->getCode(),
                    'file'       => $exception->getFile(),
                    'line'       => $exception->getLine(),
                ]);
            }
        }
    }
}
