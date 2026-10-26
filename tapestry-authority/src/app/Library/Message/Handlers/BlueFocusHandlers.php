<?php


namespace App\Library\Message\Handlers;


use App\Library\Message\Exception\HRSendMessageError;
use App\Library\Message\Response\BlueFocusResponse;
use App\Library\Message\Traits\HttpRequest;
use App\Library\Validator;

class BlueFocusHandlers extends BaseHandlers
{
    use HttpRequest;
    use Validator;

    const MESSAGE_CONTENT_DATA = '您好，您的本次登录验证码为%s，15分钟内输入有效！请勿随意告知他人。';

    private $content = '';

    const SEND_URL = "https://www.futurersms.com/api/sms/send";

    public function __construct(array $config)
    {
        parent::__construct($config);
        $this->checkConfig();
    }

    /**
     * @throws \Exception
     */
    protected function checkConfig()
    {
        $this->required([
            'user' => [true, '该字段为必填字段'],
            'pwd'    => [true, '该字段为必填字段'],
        ], $this->config);
    }

    /**
     * @inheritDoc
     */
    public function send(array $param): bool
    {
        $params = [
            'user'    => $this->config['user'],
            'pwd'       => md5($this->config['pwd']),
            'mobile'    => $param['subject'],
            'msg'   => $this->makeContent($param['data']),
            'exno'      => $param['extNo'] ?? '',
            'custid'    => $param['msgId'] ?? '',
        ];

        $result   = $this->get(self::SEND_URL, $params);
        $response = new BlueFocusResponse($result);
        $isSendSuccess = false;

        if (0 == $response->getStateCode()) {
            $isSendSuccess = true;
        }

        if (!$isSendSuccess) {
            throw new HRSendMessageError($response, $response->getCodeMsg());
        }
        return true;
    }

    /**
     * 短信内容处理
     */
    private function makeContent($data): string
    {
        $content = "";

        if (is_string($data)) {
            $content = $data;
        }

        if (isset($data['code']) && !empty($data['code'])) {
            $content = sprintf(self::MESSAGE_CONTENT_DATA, $data['code']);
        }

        return preg_replace(['/\//', '/\+/'], ['_', '-'], base64_encode($content));
    }

    /**
     * @inheritDoc
     */
    public function getAccount(): string
    {
        return md5($this->config['user']);
    }
}