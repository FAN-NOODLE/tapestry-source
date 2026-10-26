<?php
/*
 * This file is part of the car/chedianai_bc.
 *
 * (c) chedianai <i@chedianai.com>
 *
 * This source file is subject to the MIT license that is bundled
 * with this source code in the file LICENSE.
 */

namespace App\Library\Util;

use Illuminate\Support\Facades\Log;

class CurlUtil
{
    /**
     * CURL Post
     *
     * @param        $url
     * @param array  $headers
     * @param string $params
     *
     * @return array
     */
    public static function curlPost($url, $headers = [], $params = '')
    {
        // 响应Body
        $httpInfo = array();
        // 初始化curl
        $ch = curl_init();

        Log::info('请求参数: ' . var_export($params, true) . PHP_EOL);

        curl_setopt($ch, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_1);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 30);
        curl_setopt($ch, CURLOPT_TIMEOUT, 30);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
        // 证书相关检查
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE);
        curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);

        if ($headers) {
            curl_setopt($ch, CURLOPT_HEADER, 0);
            curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        }

        curl_setopt($ch, CURLOPT_POST, TRUE);
        curl_setopt($ch, CURLOPT_POSTFIELDS, $params);
        curl_setopt($ch, CURLOPT_URL, $url);

        $response = curl_exec($ch);

        Log::info('响应数据: ' . var_export($response, true) . PHP_EOL);

        if ($response === FALSE) {
            return [];
        }

        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE); // 获取响应HttpCode
        $httpInfo = array_merge($httpInfo, curl_getinfo($ch)); // 获取响应HttpBody
        curl_close($ch);

        Log::info('响应HttpCode: ' . var_export($httpCode, true) . PHP_EOL);
        Log::info('响应HttpBody: ' . var_export($httpInfo, true) . PHP_EOL);
        return json_decode($response, true);
    }

    /**
     * CURL Get
     *
     * @param        $url
     * @param array  $headers
     * @param string $params
     *
     * @return array
     */
    public static function curlGet($url, $headers = [], $params = '')
    {
        // 响应Body
        $httpInfo = array();
        // 初始化curl
        $ch = curl_init();

        Log::info('请求参数: ' . var_export($params, true) . PHP_EOL);

        curl_setopt($ch, CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_1);
        curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 30);
        curl_setopt($ch, CURLOPT_TIMEOUT, 30);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, TRUE);
        // 证书相关检查
        curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, FALSE);
        curl_setopt($ch, CURLOPT_SSL_VERIFYHOST, FALSE);

        if ($headers) {
            curl_setopt($ch, CURLOPT_HEADER, 0);
            curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);
        }

        if ($params) {
            if (is_array($params)) {
                $params = http_build_query($params);
            }
            curl_setopt($ch, CURLOPT_URL, $url . '?' . $params);
        } else {
            curl_setopt($ch, CURLOPT_URL, $url);
        }

        $response = curl_exec($ch);

        Log::info('响应数据: ' . var_export($response, true) . PHP_EOL);

        if ($response === FALSE) {
            return [];
        }

        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE); // 获取响应HttpCode
        $httpInfo = array_merge($httpInfo, curl_getinfo($ch)); // 获取响应HttpBody
        curl_close($ch);

        Log::info('响应HttpCode: ' . var_export($httpCode, true) . PHP_EOL);
        Log::info('响应HttpBody: ' . var_export($httpInfo, true) . PHP_EOL);
        return json_decode($response, true);
    }
}
