<?php

namespace App\Library\Express;

use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class Express
{
    /**
     * @var string domain
     */
    protected static $domain;

    /**
     * @var string 密钥
     */
    protected static $key;

    /**
     * Express constructor.
     */
    public function __construct()
    {
        self::$key    = env('LOGISTICS_KEY', '3417610874232448');
        self::$domain = env('LOGISTICS_DOMAIN', 'https://logistics-test.vchangyi.com/');
    }

    const LOGISTICS_KEY = [
        'sf'           => '顺丰速运',
        'huaqikuaiyun' => '华企快运'
    ];


    /**
     * 物流查询
     *
     * @param $com string 快递编码
     * @param $no string 快递单号
     * @param $senderPhone string 寄件人手机号后四位，顺丰快递需要提供senderPhone和receiverPhone其中一个
     * @param $receiverPhone string 收件人手机号后四位，顺丰快递需要提供senderPhone和receiverPhone其中一个
     *
     * @return array
     * @throws GuzzleException
     */
    public function queryExpress($com, $no, $senderPhone = '', $receiverPhone = '')
    {
        $url    = self::$domain . '/exp/index';
        $params = [
            'com'           => $com,
            'no'            => $no,
            'senderPhone'   => $senderPhone,
            'receiverPhone' => $receiverPhone,
            'key'           => self::$key
        ];

        return self::httpRequest('POST', $url, $params);
    }

    /**
     * 发起请求
     *
     * @param $method string 请求方式
     * @param $uri string 请求地址
     * @param $params array 请求数据
     *
     * @return array
     * @throws GuzzleException
     * @throws Exception
     */
    private static function httpRequest($method, $uri, $params = [])
    {
        if ($method === 'POST') {
            $response = (new Client())->request($method, $uri, [
                'http_errors' => false,
                'headers'     => [
                    'Content-Type' => 'application/x-www-form-urlencoded'
                ],
                'form_params' => $params
            ]);
        } else {
            $params   = http_build_query($params);
            $response = (new Client())->request($method, $uri . '?' . $params);
        }

        $result = $response->getBody()->getContents();
        $result = json_decode($result, true);
        Log::info('快递查询数据请求信息', [
            'request'  => $params,
            'response' => $result
        ]);

        if (is_null($result) || $response->getStatusCode() !== 200) {
            error('第三方服务异常:' . $result['reason']);
        }

        return [
            'message' => $result['reason'] ?? '',
            'data'    => $result['result'] ?? '',
            'code'    => $result['error_code'] ?? 0
        ];
    }
}