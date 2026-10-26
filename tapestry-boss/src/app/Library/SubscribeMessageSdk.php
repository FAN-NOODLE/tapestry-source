<?php

namespace App\Library;

use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;

class SubscribeMessageSdk
{
    /**
     * 订阅消息接口域名
     */
    const WXA_API_HOST = 'https://api.weixin.qq.com/wxaapi/newtmpl';

    /**
     * 微信公众平台接口域名
     */
    const CGI_BIN_HOST = 'https://api.weixin.qq.com/cgi-bin';

    protected $accessToken;

    public function __construct($enterpriseId, $appId = '', $appSecret = '')
    {
        $this->accessToken = (new AccessToken())->getAccessToken($enterpriseId, $appId, $appSecret);
    }

    /**
     * 获取小程序账号的类目
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getCategory()
    {
        $url = '/getcategory?access_token=' . $this->accessToken;

        return $this->request('GET', $url);
    }

    /**
     * 获取帐号所属类目下的公共模板标题
     *
     * @param array $params
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getPubTemplateTitleList(array $params)
    {
        $url = '/getpubtemplatetitles?access_token=' . $this->accessToken;

        return $this->request('GET', $url, $params);
    }

    /**
     * 获取模板标题下的关键词列表
     *
     * @param string $tid 模板标题id
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getPubTemplateKeyWordsById(string $tid)
    {
        $url = '/getpubtemplatekeywords?access_token=' . $this->accessToken . '&tid=' . $tid;

        return $this->request('GET', $url);
    }

    /**
     * 组合模板并添加至帐号下的个人模板库
     *
     * @param array $data
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function addTemplate(array $data)
    {
        $url = '/addtemplate?access_token=' . $this->accessToken;

        return $this->request('POST', $url, $data);
    }

    /**
     * 删除帐号下的个人模板
     *
     * @param string $priTmplId 要删除的模板id
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function deleteTemplate(string $priTmplId)
    {
        $url = '/deltemplate?access_token=' . $this->accessToken . '&priTmplId=' . $priTmplId;

        return $this->request('POST', $url);
    }

    /**
     * 获取当前帐号下的个人模板列表
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function getTemplateList()
    {
        $url = '/gettemplate?access_token=' . $this->accessToken;

        return $this->request('GET', $url);
    }

    /**
     * 发送订阅消息
     *
     * @param array $params
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function send(array $params)
    {
        $url = '/message/subscribe/send?access_token=' . $this->accessToken;

        return $this->request('POST', $url, $params, self::CGI_BIN_HOST);
    }

    /**
     * 接口请求
     *
     * @param        $method
     * @param        $url
     * @param        $params
     * @param string $host
     *
     * @return array|mixed|string
     * @throws \GuzzleHttp\Exception\GuzzleException
     */
    public function request($method, $url, $params = [], $host = self::WXA_API_HOST)
    {
        // GET请求，拼接参数
        if ('GET' == $method && !empty($params)) {
            $url = $url . '&' . http_build_query($params);
            $params = [];
        }

        try {
            $url = $host . $url;
            $result = (new Client())->request($method, $url, [
                'headers'     => [
                    'Content-Type' => 'application/json',
                    'Accept'       => 'application/json',
                ],
                'body'        => json_encode($params, JSON_UNESCAPED_UNICODE),
                'http_errors' => false
            ]);

            $response = $result->getBody()->getContents();
            $response = json_decode($response, true);
        } catch (\Exception $e) {
            Log::debug('【订阅消息】请求异常：' . $e->getMessage());

            return [];
        }

        return $response;
    }
}
