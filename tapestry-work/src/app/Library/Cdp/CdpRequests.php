<?php

namespace App\Library\Cdp;

use App\Library\Util\CurlUtil;

/**
 * Class CdpRequests.
 */
class CdpRequests
{
    private $domain;

    public function __construct()
    {
        $this->domain = 'http://ckapi.vchangyi.com/ckapi/useract'; // 接口请求地址
    }

    /**
     * Get 请求
     *
     * @param array $params
     *
     * @return mixed
     * @throws
     */
    public function doGet(array $params = [])
    {
        // 请求地址
        $url = $this->getUrl();

        // 组装Headers
        $headers = [
            'Content-Type: application/json;charset=UTF-8',
        ];

        return CurlUtil::curlGet($url, $headers, json_encode($params));
    }

    /**
     * POST 请求
     *
     * @param array $params
     *
     * @return mixed
     * @throws
     */
    public function doPost(array $params = [])
    {
        // 请求地址
        $url = $this->getUrl();

        // 组装Headers
        $headers = [
            'Content-Type: application/json;charset=UTF-8',
        ];

        return CurlUtil::curlPost($url, $headers, json_encode($params));
    }

    /**
     * 拼接接口域名
     * @return string
     */
    private function getUrl()
    {
        return $this->domain;
    }
}
