<?php
/**
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/5/24
 * Time: 17:16
 */

namespace App\Handlers\WxPayApi;

use App\Exceptions\Error;
use App\Handlers\PaymentHandler;
use App\Models\Test;
use App\Services\WxPayAuthService;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;

class WxPayClient
{
    /**
     * 微信支付 v3 host
     */
    const WX_PAYMENT_HOST = 'https://api.mch.weixin.qq.com/v3';

    /**
     * http请求实例
     *
     * @var Client|null
     */
    protected $httpClient = null;

    public function __construct()
    {
        $this->httpClient = new Client();
    }

    /**
     * post请求
     *
     * @param       $url
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    protected function wxHttpPost($url, $data = [])
    {
        $method = 'POST';
        $apiUrl = $this->getApiUrl($url, $data['query']);

        try {
            $authorization = $this->getAuthorization($apiUrl, $method, $data['body'], $data['paymentConfig']);

            $response = $this->httpClient->request($method, $apiUrl, [
                'headers'     => [
                    'Accept'        => 'application/json',
                    'Content-Type'  => 'application/json',
                    'Authorization' => $authorization
                ],
                'body'        => json_encode($data['body']),
                'http_errors' => false
            ]);

            $result = $response->getBody();
            $result = json_decode($result, true);

        } catch (\Exception $e) {
            Log::debug('wxApiPostRequest:Error ' . $e->getMessage());
            return $e->getMessage();
        }
        return $result;
    }

    /**
     * get请求
     *
     * @param       $url
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    protected function wxHttpGet($url, $data = [])
    {
        $method = 'GET';
        $apiUrl = $this->getApiUrl($url, $data['query']);

        try {
            $authorization = $this->getAuthorization($apiUrl, $method, [], $data['paymentConfig']);
            $response      = $this->httpClient->request($method, $apiUrl, [
                'headers'     => [
                    'Content-Type'  => 'application/json',
                    'Accept'        => 'application/json',
                    'Authorization' => $authorization
                ],
                'http_errors' => false
            ]);
            $result        = $response->getBody();
            $result        = json_decode($result, true);
            Log::info('请求微信支付V3接口', [
                'url'            => $apiUrl,
                'method'         => $method,
                'requestBody'    => $data,
                'responseBody'   => $result,
                'responseHeader' => $response->getHeaders(),
            ]);
        } catch (\Exception $e) {
            Log::debug('wxApiGetRequest:Error ' . $e->getMessage());
            throw new Error(10000, $e->getMessage());
        }


        return $result;
    }

    /**
     * get请求
     *
     * @param       $url
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    protected function wxHttpRequest($url, $data = [], $method, $host = self::WX_PAYMENT_HOST)
    {
        $apiUrl = $this->getApiUrl($url, $data['query'], $host);

        $body = $data['body'] ?? [];

        try {
            $authorization = $this->getAuthorization($apiUrl, $method, $body, $data['paymentConfig']);
            $params        = [
                'headers'     => [
                    'Content-Type'  => 'application/json',
                    'Accept'        => 'application/json',
                    'Authorization' => $authorization
                ],
                'http_errors' => false
            ];

            if (!empty($body)) {
                $params['body'] = json_encode($body);
            }

            $response = $this->httpClient->request($method, $apiUrl, $params);
            $result   = $response->getBody();
            $result   = json_decode($result, true);
            Log::info('请求微信支付V3接口', [
                'url'            => $apiUrl,
                'method'         => $method,
                'requestBody'    => $data,
                'responseBody'   => $result,
                'responseHeader' => $response->getHeaders(),
            ]);

        } catch (\Exception $e) {
            Log::debug('wxApiRequest:Error ' . $e->getMessage());

            \App\Models\Error::log('wxApiRequest', $e);
            throw new Error(10000, $e->getMessage());
        }


        return $result;
    }

    /**
     * 获取url
     *
     * @param $url
     * @param $data
     *
     * @return string
     */
    private function getApiUrl($url, $data, $host = self::WX_PAYMENT_HOST)
    {
        $query = http_build_query($data);
        return empty($query) ? $host . $url : $host . $url . '?' . $query;
    }

    /**
     * 获取验签
     *
     * @param $url
     * @param $method
     * @param $body
     * @param $paymentConfig
     *
     * @return mixed
     */
    protected function getAuthorization($url, $method, $body, $paymentConfig)
    {
        return (new WxPayAuthService(new PaymentHandler()))->getWxPayAuthorization($url, $method, $body, $paymentConfig);
    }
}
