<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/10/12
 * Time: 10:37 AM
 */

namespace App\Library\Uc;

use App\Exceptions\Error;
use App\Services\EnterpriseService;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\ClientException;
use Illuminate\Support\Facades\Log;

abstract class Abstracts
{
    /**
     * GuzzleHttp\Client
     *
     * @var Client|null
     */
    protected $client = null;

    /**
     * Http request header
     *
     * @var array
     */
    protected $headers = [
        'Content-Type' => 'application/json',
        'Referer' => 'crs.vchangyi.com'
    ];

    /**
     * 初始化
     *
     * Abstracts constructor.
     */
    public function __construct()
    {
        $this->client = new Client(['verify' => false]);
    }

    /**
     * 获取UC api地址
     *
     * @return mixed
     * @throws \Exception
     */
    protected function getUcApiUrl()
    {
        $apiUrl = config('client.api_url.uc_api');

        if (empty($apiUrl)) {
            throw new \Exception('UC API url error');
        }

        return $apiUrl;
    }

    /**
     * 发起POST请求
     *
     * @param $url
     * @param array $data
     * @return mixed
     * @throws Error
     */
    protected function post($url, $data = [])
    {
        Log::debug('UC api: '.$url);

        $data = empty($data) ? '{}' : json_encode($data);

        $options = [
            'headers' => $this->headers,
            'body' => $data
        ];

        try {
            $response = $this->client->post($url, $options);
            $data = $response->getBody();
        } catch (ClientException $clientException) {
            throw new Error(1000, $clientException->getMessage(), 500);
        }

        return $this->result($data);
    }

    /**
     * 发起GET请求
     *
     * @param $url
     * @param array $data
     * @return mixed
     */
    protected function get($url, $data = [])
    {
        $url = $url.'?'.http_build_query($data);
        Log::debug('UC api: '.$url);

        $options = [
            'headers' => $this->headers
        ];

        try {
            $response = $this->client->get($url, $options);
            $data = $response->getBody();
        } catch (ClientException $clientException) {
            throw new Error(1000, $clientException->getMessage(), 500);
        }

        return $this->result($data);
    }

    /**
     * @param $data
     * @return mixed
     * @throws Error
     */
    private function result($data)
    {
        Log::debug('UC Response: '.$data);
        $array = json_decode($data, true);
        $requestId = empty($array['requestId']) ? 'NOT REQUEST ID' : $array['requestId'];

        // 接口响应成功，直接返回结果
        if (!empty($array['code']) && $array['code'] === 'SUCCESS') {
            $result['data'] = empty($array['data']) ? [] : $array['data'];
            // 防止返回结果不是数组形式报错
            if (is_array($result)) {
                $result['requestId'] = $requestId;
            }
            return $result;
        }

        // 输出这个错误
        if (!empty($array['code']) && $array['code'] == '40029') {
            throw new Error(100000, $array['code'].":".$array['msg']);
        }

        // 输出这个错误
        if (!empty($array['code']) && $array['code'] == '60008') {
            throw new Error(100000, $array['code'].":".$array['msg']);
        }

        // 输出这个错误
        if (!empty($array['code']) && $array['code'] == '60104') {
            throw new Error(100000, $array['code'].":".$array['msg']);
        }

        // 输出这个错误
        if (!empty($array['code']) && $array['code'] == '40071') {
            throw new Error(100000, $array['code'].":已在企业微信后台设置了该标签，请点击同步标签");
        }

        // 提示性错误不以异常抛出，业务层做判断处理
        if (!empty($array['code']) && is_numeric($array['code'])) {
            return $data;
        }

        Log::error('UC error: '.json_encode($array));
        $error = '[UC] '.$array['msg'];

        // 详细字段报错信息
        if (!empty($array['data'][0]['field'])) {
            $error .= ' field: '.$array['data'][0]['field'].' message: '.$array['data'][0]['message'];
        }

        throw new Error(100000, $error);
    }

    /**
     * 上传图片
     * @param       $url
     * @param array $multipart
     *
     * @return mixed
     */
    protected function postMultipartUpload($url, $multipart= [])
    {
        $options = [
            'multipart' => $multipart
        ];

        try {
            $response = $this->client->post($url, $options);
            $data = $response->getBody();
        } catch (ClientException $clientException) {
            throw new Error(1000, $clientException->getMessage(), 500);
        }

        return $this->result($data);

    }

    /**
     * 查询获取应用标示
     *
     * @param $qyDomain
     * @return string
     * @throws \Exception
     */
    protected function getAppIdentifier($qyDomain)
    {
        $appIdentifier = 'contact';

        $pluginList = [
            'crs-shop',
            'crs-manu-customer',
            'crs_store_assistant_self'
        ];

        $detail = (new EnterpriseService())->getDomainInfo($qyDomain);
        if ($detail['contact_status'] != 1) {
            $usePluginList = (new EnterpriseWorkPlugin())->lists([
                'enterprise_id' => $detail['id']
            ]);
            $usePluginIdList = array_column($usePluginList, 'identifier');

            Log::debug('安装的标识', $usePluginIdList);

            foreach ($pluginList as $item) {
                if (in_array($item, $usePluginIdList)) {
                    $appIdentifier = $item;
                    break;
                }
            }
        }


        Log::debug('最终返回的标识', [$appIdentifier]);

        return $appIdentifier;
    }
}
