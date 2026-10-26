<?php

namespace App\Sdk\Active\CRM\SwCoupon;

use App\Library\WechatNotice\Chat;
use Exception;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;

class SwCouponClient
{
    private $appId;
    private $secret;
    private $aesKey;
    private $baseUrl;
    private $version;
    private $client;

    public function __construct()
    {
        // sw测试账号
        $testConfig = [
            'app_id' => 'SW-CY-Acxiom-Openapi-20230001',
            'secret' => 'REDACTED_FOR_SOURCE_SHARING',
            'aes_key' => 'REDACTED_FOR_SOURCE_SHARING',
            'base_url' => 'https://cieapiuat.cdp.stuartweitzman.cn/openapi/api',

        ];
        // sw生产账号
        $prodConfig = [
            'app_id' => 'SW-CY-Acxiom-Openapi-20250001',
            'secret' => 'REDACTED_FOR_SOURCE_SHARING',
            'aes_key' => 'REDACTED_FOR_SOURCE_SHARING',
            'base_url' => 'https://cieapi.cdp.stuartweitzman.cn/openapi/api',
        ];
        $config = $testConfig;
        if (config("app.env") == 'prod') {
            $config = $prodConfig;
        }
        $this->appId = $config['app_id'];
        $this->secret = $config['secret'];
        $this->aesKey = $config['aes_key'];
        $this->baseUrl = $config['base_url'];
        $this->version = '2.0';

        $this->client = new Client([
            'timeout' => 30,
            'connect_timeout' => 10,
            'http_errors' => false,
        ]);
    }

    /**
     * 获取用户优惠券列表
     * couponStatus 4 已生效
     */
    public function getUserCouponList($unionId, $page = 1, $size = 10,$couponStatus = 4)
    {
        $customerInfo = [
            'union_id' => $unionId,
        ];

        $businessInfo = [
            'page' => [
                'current' => $page,
                'size' => $size
            ]
        ];

        if ($couponStatus) {
            $businessInfo['couponStatus'] = $couponStatus;
        }

        $requestParam = [
            'customerInfo' => $customerInfo,
            'businessInfo' => json_encode($businessInfo)
        ];

        return $this->request('coupon.getUserCouponList', $requestParam);
    }

    /**
     * 发送请求到开放接口
     */
    public function request($method, $requestParam)
    {
        $timestamp = round(microtime(true) * 1000);
        $url = $this->baseUrl . "?appId={$this->appId}&method={$method}&version={$this->version}&timestamp={$timestamp}";

        // 加密请求体
        $body = json_encode($requestParam);
        $encryptedBody = $this->aesEncode($body, $this->aesKey);
        // 生成签名
        $signString = $this->generateSignString($this->appId, $this->secret, $method, $this->version, $timestamp, $encryptedBody);
        $sign = $this->sign($signString);
        // 准备请求头和请求体
        $headers = [
            'sign' => $sign,
            'Content-Type' => 'application/json',
            'User-Agent' => 'OpenApiClient/1.0',
        ];

        $requestBody = json_encode(['body' => $encryptedBody]);
        try {
            $response = $this->client->post($url, [
                'headers' => $headers,
                'body' => $requestBody,
                'verify' => false, // 如果需要忽略SSL验证，设置为false
            ]);

            $statusCode = $response->getStatusCode();
            $responseBody = $response->getBody();
            if ($statusCode >= 200 && $statusCode < 300) {
                $result = json_decode($responseBody->getContents(), true);
                if (isset($result['code'])&& $result['code'] == '200') {
                    $result['data'] = json_decode($this->aesDecode($result['data'], $this->aesKey), true);
                } else {
                    Chat::getInstance()->send('请求SW-CRM优惠券接口出错1', [
                        'env' => config('app.env'),
                        'api' => $url,
                        'param' => $body,
                        'response' => $result
                    ]);
                }
                Log::info('SW优惠券接口请求成功', [
                    'api' => $url,
                    'param' => $body,
                    'response' => $result
                ]);
                return $result;
            } else {
                Log::error('SW优惠券接口请求失败1', [
                    'url' => $url,
                    'status' => $statusCode,
                    'response' => $responseBody->getContents()
                ]);
                Chat::getInstance()->send('请求SW-CRM优惠券接口出错2', [
                    'env' => config('app.env'),
                    'api' => $url,
                    'param' => $body,
                    'response' => $responseBody->getContents(),
                ]);
                return [
                    'error' => true,
                    'status_code' => $statusCode,
                    'response' => $responseBody->getContents(),
                    'message' => 'SW优惠券接口请求失败2'
                ];
            }
        } catch (Exception $e) {
            Log::error('OpenAPI其他异常', [
                'url' => $url,
                'error' => $e->getMessage(),
                'trace' => $e->getTraceAsString()
            ]);

            return [
                'error' => true,
                'message' => '系统异常',
                'exception' => get_class($e)
            ];
        }
    }


    /**
     * AES加密
     *
     * @param $content
     * @param $key
     * @return string|null
     */
    private function aesEncode($content, $key): ?string
    {
        if (empty($content) || empty($key)) {
            return null;
        }
        $iv = 'j5A8Gq0YQK2HUlMc';
        $encrypted = openssl_encrypt(
            $content,
            'AES-128-CBC',
            $key,
            OPENSSL_RAW_DATA,
            $iv
        );

        return base64_encode($encrypted);
    }

    private function aesDecode($content, $key)
    {
        if (empty($content) || empty($key)) {
            return null;
        }

        $iv = 'j5A8Gq0YQK2HUlMc'; // 与Java代码中的KEY_VI一致
        return openssl_decrypt(
            base64_decode($content),
            'AES-128-CBC',
            $key,
            OPENSSL_RAW_DATA,
            $iv
        );
    }

    /**
     * 生成签名字符串
     */
    private function generateSignString($appId, $secret, $method, $version, $timestamp, $body): string
    {
        return "appId={$appId}&secret={$secret}&method={$method}&version={$version}&timestamp={$timestamp}&body={$body}";
    }

    /**
     * MD5签名
     */
    private function sign($signString): string
    {
        return md5($signString);
    }

    /**
     * 获取Guzzle客户端实例
     */
    public function getClient(): Client
    {
        return $this->client;
    }
}