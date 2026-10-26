<?php
/**
 * Created by PhpStorm.
 * User: guoxingong
 * Date: 2019/4/18
 * Time: 10:30
 */

namespace App\Library;

use App\Exceptions\Error;
use CrsSdk\HttpRunner;
use GuzzleHttp\Client;


class TencentLbs
{
    const URL = 'https://apis.map.qq.com';
    const KEY = 'I3KBZ-P4MKP-T7WDN-VVUQG-YGYYO-DQFU6';
    // 签名校验参数
    const SK = 'xmodqucF4HB6Xi71rzrHUweLkIZYZn9';

    protected $httpClient = null;

    /**
     * 成功的http状态
     */
    const SUCCESS_CODE = 200;
    public function __construct()
    {
        $this->httpClient = new Client();
    }

    /**
     * 逆地址解析
     * @param $lon
     * @param $lat
     * @return array|mixed
     * @throws \Exception
     * author: zhaojie<zhaojie@vchangyi.com>
     * create: 2020/5/12 12:08
     */
    public function getAddressByLocation($lon, $lat)
    {
        $path = "/ws/geocoder/v1?";
        $params = [
            'key' => self::KEY,
            'location' => "{$lat},{$lon}"
        ];
        $paramUrl = $this->getSignPath($path, $params);
        $params['sig'] = $paramUrl;
        $result = $this->httpGet($paramUrl, []);
        return $result;
    }

    /**
     * get请求 带sign的path生成
     * @param $path
     * @param $params
     * @return string
     * author: zhaojie<zhaojie@vchangyi.com>
     * create: 2020/5/11 17:51
     */
    private function getSignPath($path, $params)
    {
        ksort($params);
        $paramStr = urldecode(http_build_query($params , '' , '&'));
        $sign = md5($path . $paramStr . self::SK);
        return $path . $paramStr . '&sig=' . $sign;
    }

    /**
     * 发起http请求
     *
     * @param       $url
     * @param array $data
     *
     * @return mixed
     * @throws \Exception
     */
    protected function httpPost($url, $data = [])
    {
        return $this->request('POST', $url, $data);
    }

    /**
     * 发起http请求
     *
     * @param       $url
     * @param array $data
     *
     * @return mixed
     * @throws \Exception
     */
    protected function httpGet($url, $data = [])
    {
        return $this->request('GET', $url, $data);
    }

    /**
     * @param       $method
     * @param       $url
     * @param array $data
     *
     * @return mixed
     * @throws \Exception
     */
    public function request($method, $url, $data = [])
    {
        $apiUrl = self::URL . $url;
        $data = empty($data) ? '{}' : json_encode($data);
        $headers = [
            'Content-Type'   => 'application/json'
        ];
        $httpRunner = new HttpRunner($method, $apiUrl, $headers, $this->httpClient, $data);
        $response = $httpRunner->execute();
        $result = $response->getBody()->getContents();
        $result = json_decode($result, true);

        if (is_null($result)) {
            throw new Error(80000, '[' . $method . '] 服务接口返回数据异常');
        }

        if (isset($result['status']) && $result['status'] !== 0) {
            $code = $result['status'];
            $msg  = $result['message'];
            throw new Error($code, $msg);
        }

        if (isset($result['status']) && $result['status'] === 0) {
            $result = $result['result'] ?? [];
        }
        return $result;
    }
}