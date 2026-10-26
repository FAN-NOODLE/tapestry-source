<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/13
 * @Time: 18:47
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: HRMessage.php
 */

namespace App\Library\Message\Handlers;

use App\Library\Message\Contracts\SendHandlerInterface;
use App\Library\Message\Traits\HttpRequest;
use App\Library\Message\Exception\HRSendMessageError;
use App\Library\Message\Response\HRResponse;

class HRMessageHandlers extends BaseHandlers
{
    use HttpRequest;

    /**
     * 请求地址
     */
    const ENDPOINT_URL = 'http://api.shgmnets.com/sms/send';

    /**
     * 内容模板
     */
    const MESSAGE_CONTENT_DATA = '（欧家智慧导购）您好，您的本次登录验证码为%s，15分钟内输入有效！请勿随意告知他人。';

    /**
     * 是否发送成功
     */
    private $isSendSuccess = false;

    /**
     * 发送内容
     * @var string
     */
    private $content = '';

    public function __construct(array $config)
    {
        parent::__construct($config);
    }

    /**
     * 发送短信
     * @param array $param
     * @return bool
     * @throws \Exception
     */
    public function send(array $param): bool
    {
        try{
            $this->checkRequestParams($param);

            $params = [
                'user' => $this->config['user'],
                'pwd' => $this->md5Password($this->config['pwd']),
                'mobile' => $param['subject'],
                'msg' => $this->getContent($param['data']),
                'extNo' => isset($param['extNo']) ? $param['extNo'] : '',
                'msgId' => isset($param['msgId']) ? $param['msgId'] : '',
            ];

            $result = $this->get(self::ENDPOINT_URL, $params);

            $response = new HRResponse($result);

            $this->isSendSuccess = false;
            if($response->getStateCode() == 0 && $response->getMsgId() != ''){
                $this->isSendSuccess = true;
            }

            if (!$this->isSendSuccess) {
                throw new HRSendMessageError($response, $response->getCodeMsg());
            }

        }catch (\Exception $exception) {
            throw $exception;
        }

        return $this->isSendSuccess;
    }

    /**
     * 加密密码
     * @param $password
     * @return string
     */
    public function md5Password($password)
    {
        return strtolower(md5($password));
    }

    /**
     * 校验参数
     * @param $messageData
     * @throws \Exception
     */
    public function checkRequestParams($messageData)
    {
        if (!isset($messageData['subject']) || empty($messageData['subject'])) {
            throw new \Exception('subject is empty');
        }

        if (!isset($messageData['data']) || empty($messageData['data'])) {
            throw new \Exception('data is empty or not found');
        }
    }

    /**
     * 获取发送短信的内容
     * @param $data
     * @return string
     */
    public function getContent($data)
    {
        if(isset($data['code']) && !empty($data['code'])){
           $this->content = sprintf(self::MESSAGE_CONTENT_DATA, $data['code']);
        }

        if(is_string($data)){
            $this->content = $data;
        }

        $this->content = base64_encode($this->content);

        $this->content = str_replace('+', '-', $this->content);
        $this->content = str_replace('/', '_', $this->content);
        $this->content = str_replace('=', '', $this->content);

        return $this->content;
    }

    /**
     * 获取调用接口的账号
     * @return string
     */
    public function getAccount(): string
    {
        return md5($this->config['user']);
    }
}