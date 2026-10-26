<?php


namespace App\Library\Message\Handlers;

use App\Library\Message\Exception\HRSendMessageError;
use App\Library\Message\Response\HRResponse;
use App\Library\Message\Response\MWResponse;
use App\Library\Message\Traits\HttpRequest;
use App\Library\Validator;
use Exception;

/**
 * Class MWHandler
 * @package App\Library\Message\Handlers
 *
 * 梦网短信服务平台
 */
class MWHandler extends BaseHandlers
{
    use HttpRequest;
    use Validator;


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

    // 梦网短信平台接口集合
    const SEND_URLS = [
        // 单条发送接口
        'send_single' => 'http://%s:%s/sms/v2/std/send_single',
    ];
    // 密码加密固定字符串
    const ENCRYPTION_STR = '00000000';
    // 梦网配置
    //protected $config = [];

    /**
     * MWHandler constructor.
     * @param array $config 发送短信需要的配置
     * @throws Exception
     */
    public function __construct(array $config)
    {
        $this->config = $config;

        $this->checkConfig();
    }

    /**
     * 检测企业的账号配置是否
     * @throws Exception
     */
    protected function checkConfig()
    {
        $this->required([
            'userid' => [true, '该字段为必填字段'],
            'pwd'    => [true, '该字段为必填字段'],
        ], $this->config);
    }

    public function send(array $param): bool
    {
        return $this->sendSingle($param);
    }

    /**
     * 单条发送
     */
    private function sendSingle(array $param): bool
    {
        $timestamp = (string)time();

        $params = [
            'userid'    => $this->config['userid'],
            'pwd'       => $this->encryptionPassword($timestamp),
            'mobile'    => $param['subject'],
            'content'   => $this->makeContent($param['data']),
            'timestamp' => $timestamp,
            'exno'      => $param['extNo'] ?? '',
            'custid'    => $param['msgId'] ?? '',
        ];

        $url      = sprintf(self::SEND_URLS['send_single'], $this->config['ip'], $this->config['port']);
        $result   = $this->post($url, $params);
        $response = new MWResponse($result);

        $this->isSendSuccess = false;

        if ($response->getStateCode() == 0 && $response->getMsgId() != '') {
            $this->isSendSuccess = true;
        }

        if (!$this->isSendSuccess) {
            throw new HRSendMessageError($response, $response->getCodeMsg());
        }
        return $this->isSendSuccess;
    }

    /**
     * 密码加密
     * 加密规则：strtolower(md5(账号+固定字符串+密码+时间戳))
     */
    private function encryptionPassword(string $timestamp): string
    {
        $str = sprintf("%s%s%s%s", $this->config['userid'], self::ENCRYPTION_STR,
            $this->config['pwd'], $timestamp);

        return strtolower(md5($str));
    }

    /**
     * 短信内容处理
     */
    private function makeContent($data): string
    {
        $content = "";

        if (is_string($data)) {
            $content = urlencode($data);
        }

        if (isset($data['code']) && !empty($data['code'])) {
            $content = sprintf(self::MESSAGE_CONTENT_DATA, $data['code']);
        }

        return $content;
    }

    /**
     * 获取调用接口的账号
     * @return string
     */
    public function getAccount(): string
    {
        return md5($this->config['userid']);
    }
}