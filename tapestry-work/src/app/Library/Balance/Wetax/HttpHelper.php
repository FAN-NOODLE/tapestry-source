<?php


namespace App\Library\Balance\Wetax;


use App\Exceptions\Error;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;
use Illuminate\Support\Facades\Log;

/**
 * 接口请求工具类
 * Class HttpTool
 * @package App\Library\Balance\Wetax
 */
class HttpHelper
{
    const CONTENT_TYPE = 'application/json';

    /**
     * 接口版本
     */
    const API_VERSION = "2.0";

    /**
     * 签名方法
     */
    const API_SIGN_TYPE = 'sha256';

    /**
     * api请求超时时间
     */
    const API_TIMEOUT = 10;

    /**
     * 接口请求/回调 成功 code
     */
    const SUCCESS_CODE = 0;

    /**
     * @var SettlementConfig
     */
    private $config;

    /**
     * @var Client
     */
    private $client;

    /**
     * HttpHelper constructor.
     * @param SettlementConfig $config
     */
    public function __construct(SettlementConfig $config)
    {
        $this->config = $config;
        $this->client = new Client(['base_uri' => rtrim($this->config->getHost(), '/')]);
    }

    /**
     * 发送接口请求
     * @param string $uri
     * @param array $data
     * @param string $projectCode
     * @param string $callbackUrl
     * @return mixed
     */
    public function request(string $uri, array $data, string $projectCode = '', string $callbackUrl = '')
    {
        Log::info("balance settlement request params", func_get_args());
        $headers = $this->getCommonHeader($projectCode, $callbackUrl);
        $body = json_encode([
            // 传输数据进行加密
            'data' => EncryptUtil::encrypt(json_encode($data), $this->config->getAppSecret())
        ]);
        $options = [
            'headers' => $headers,
            'body' => $body,
            'timeout' => self::API_TIMEOUT
        ];

        Log::info("balance settlement request start", ['headers' => $headers, 'body' => $body]);
        $startTime = microtime(true);
        $responseData = null;
        try {
            $response = $this->client->post($uri, $options);
            $responseData = $response->getBody();
            Log::info("balance settlement request success", [
                'httpCode' => $response->getStatusCode(),
                'responseData' => (string) $responseData
            ]);
        } catch (RequestException $e) {
            Log::error("balance settlement request exception", [
                'code' => $e->getCode(),
                'msg' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw new SettlementException(SettlementException::HTTP_REQUEST_FAILURE, $e->getMessage(), $e);
        } finally {
            $endTime = microtime(true);
            Log::info("balance settlement request time", [
                'startTime' => $startTime,
                'endTime' => $endTime,
                'duration' => sprintf("%.3f", $endTime - $startTime),
            ]);
        }
        return $this->getResponseData($responseData);
    }

    /**
     * 处理http返回结果
     * @param string $responseData
     * @return mixed
     */
    private function getResponseData(string $responseData)
    {
        $result = json_decode($responseData, true);
        if (!isset($result['code'])) {
            // 抛出异常
            throw new SettlementException(SettlementException::CODE_ERROR, "code参数不存在");
        }
        if (self::SUCCESS_CODE != $result['code']) {
            // 抛出异常
            $failReason = $result['msg'] ?? '错误code:' . $result['code'];
            throw new SettlementException(SettlementException::CODE_ERROR, $failReason);
        }
        return $result['data'];
    }

    /**
     * 公共header信息
     * @param string $projectCode
     * @param string $callbackUrl
     * @return array
     */
    private function getCommonHeader(string $projectCode = '', string $callbackUrl = '') : array
    {
        $timestamp = time();
        $hashKey = $this->config->getAppKey() . $timestamp . $this->config->getAppSecret();
        $signature = hash_hmac(self::API_SIGN_TYPE, $hashKey, $this->config->getAppSecret()); // 签名
        $requestId = $_SERVER['x_request_id'] ?? md5(uniqid()) . rand(10000, 99999);
        $header = [
            'Content-Type' => self::CONTENT_TYPE,
            'appkey' => $this->config->getAppKey(),
            'request_id' => $requestId, // 请求标识，每次请求唯一
            'timestamp' => $timestamp, // 时间戳，30分钟过期
            'sign_type' => self::API_SIGN_TYPE, // 签名方法
            'signature' => $signature,
            'version' => self::API_VERSION,
        ];
        // 回调地址
        if (!empty($callbackUrl)) {
            $header['callback_url'] = base64_encode($callbackUrl);
        }
        // 项目编码
        if (!empty($projectCode)) {
            $header['project_code'] = $projectCode;
        }
        return $header;
    }

}