<?php
/**
 * Created by Phpstorm.
 * User : zoulongbo
 * Date : 2019-10-31
 * Tome : 10:12
 */

namespace App\Library;

use App\Exceptions\Error;
use App\Services\EnterpriseLogService;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;

class OpenPlatformClient
{
    /**
     * 微信开放平台 代小程序业务相关接口host
     */
    const WX_OPEN_MINI_HOST =  'https://api.weixin.qq.com/wxa';

    /**
     * @var Client
     */
    protected $httpClient;

    public function __construct()
    {
        $this->httpClient = new Client();
    }

    /**
     * post请求
     * @param       $url
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    protected function httpPost($url, $data = [])
    {
        $method = 'POST';
        $apiUrl = $this->getApiUrl($url, $data['query']);

        try {
            $response      = $this->httpClient->request($method, $apiUrl, [
                'headers'     => [
                    'Accept'        => 'application/json',
                    'Content-Type'  => 'application/json'
                ],
                'body'        => json_encode($data['body']),
                'http_errors' => false
            ]);

            $result        = $response->getBody();
            $result        = json_decode($result, true);

        } catch (\Exception $e) {
            Log::debug('openMiniPostRequest:Error ' . $e->getMessage());
            EnterpriseLogService::error('OpenPlatformError', [
                'url'     => $url,
                'data'    => $data,
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine()
            ]);
            return $e->getMessage();
        }
        return $result;
    }

    /**
     * get请求
     * @param       $url
     * @param array $data
     *
     * @return mixed|\Psr\Http\Message\StreamInterface
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    protected function httpGet($url, $data = [])
    {
        $method = 'GET';
        $apiUrl = $this->getApiUrl($url, $data['query']);
        try {
            $response      = $this->httpClient->request($method, $apiUrl, [
                'headers'     => [
                    'Content-Type'  => 'application/json',
                    'Accept'        => 'application/json'
                ],
                'body'        => '{}',
                'http_errors' => false
            ]);
            $result        = $response->getBody();
            $result        = json_decode($result, true);
        } catch (\Exception $e) {
            Log::debug('openMiniGetRequest:Error ' . $e->getMessage());
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
    private function getApiUrl($url, $data)
    {
        $query = http_build_query($data);
        return empty($query) ? self::WX_OPEN_MINI_HOST . $url : self::WX_OPEN_MINI_HOST . $url . '?' . $query;
    }
}
