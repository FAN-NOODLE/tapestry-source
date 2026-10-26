<?php

namespace App\Sdk\EcHub;

use App\Exceptions\Error;
use App\Library\WechatNotice\Chat;
use Exception;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\GuzzleException;
use Illuminate\Support\Facades\Log;


class Request
{
    /**
     * @var $instance //实例类本身 self
     */
    public static $instance;


    /**
     * @var Client 请求类
     */
    protected $request;

    /**
     * @var $enterpriseId //企业的id
     */
    protected $enterpriseId;

    protected $url;


    public function __construct(int $enterpriseId)
    {
        $this->request = new Client();
        $this->enterpriseId = $enterpriseId;
    }

    public static function Instance($enterpriseId): Request
    {
        if (empty(self::$instance)) {
            self::$instance = new self($enterpriseId);
        }
        return self::$instance;

    }


    public function setUrl(string $url): Request
    {
        $baseUrl=Config::getInfo()["baseUrl"];
        $this->url=$baseUrl;
        switch ($url) {
            case "saAssignment":
                $this->url.="/api/hub/weclient/contact/bind";
                break;
            case "saUnAssignment":
                $this->url.="/api/hub/weclient/contact/unbind";
                break;
            case "url-scheme":
                $this->url.="/api/hub/v1/wx/url-scheme";
                break;
            case "qr-code":
                $this->url.="/api/hub/v1/wx/qr-code";
                break;
        }
        return $this;
    }

    /**
     * 接口请求
     * @param $uri
     * @param $params
     * @return array|mixed
     * @throws GuzzleException
     */
    public function sendPost($params)
    {
        Log::info("调用EC_HUB接口开始");
        $options=[];
        if (!empty($params)) {
            $options = $this->setOptions($params);
        }
        $responseContent = "";
        try {
            $response = $this->request->request("POST", $this->url, $options);
            $responseContent = $response->getBody()->getContents();
            if ($response->getStatusCode() != 200) {
                throw new Error($response->getStatusCode(), $responseContent);
            }
            Log::info("请求EC_HUB接口返回", ["url" => $this->url, "request" => $options, "response" => $responseContent]);
        } catch (Exception $e) {
            Log::info("请求EC_HUB接口返回异常", ["url" => $this->url, "request" => $options, "response" => $e->getMessage()]);
            Chat::getInstance()->send('请求EC_HUB接口出错', [
                'env' => config('app.env'),
                'api' => $this->url,
                'param' => $params,
                'error_info' => $e->getMessage()
            ]);
        }
        Log::info("调用EC_HUB接口结束");
        if (!empty($responseContent)) {
            return json_decode($responseContent, true) ?? [];
        }
        return [];
    }

    /**
     * @param $params
     * @return array
     */
    protected function setOptions($params): array
    {
        $body=$this->getQueryParam($params);
        $header=$this->getHeader($params);
        $data["headers"]=$header;
        return array_merge($data, $body);
    }

    protected function getQueryParam($params, $requestType = 'json'): array
    {
        if ($requestType == 'json') {
            return ['json' => $params];
        } else if ($requestType == 'form_data') {
            return ['form_params' => $params];
        } else if ($requestType == 'xml') {
            return ['body' => $params];
        } else {
            return ['data' => $params];
        }

    }

    /**
     * @param $params
     * @return array
     */
    protected function getHeader($params): array
    {
        $nowTime=time();
        $notice=$this->generateRandomNotice();
        $data=[
            "x-hub-clientid"=>Config::getInfo()['appkey'],
            "x-hub-store"=>'coach',
            "x-hub-nonce"=>$notice,
            "x-hub-timestamp"=>$nowTime,
            "x-hub-version"=>"2.0",
        ];
        $signature=$this->generateSignature($data, $params);
        $header=$data;
        $header['x-hub-sign']=$signature;
        return $header;
    }


    protected function generateRandomNotice($length = 16) {
        $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
        $randomString = '';

        for ($i = 0; $i < $length; $i++) {
            $randomString .= $characters[rand(0, strlen($characters) - 1)];
        }

        return strtoupper($randomString);
    }

    protected function generateSignature($data,$body) {
        $secret = Config::getInfo()['appsecret'];
        $content=json_encode($body);
        $str="clientid={$data["x-hub-clientid"]}&key=$secret&store={$data["x-hub-store"]}&nonce={$data["x-hub-nonce"]}&timestamp={$data["x-hub-timestamp"]}&version={$data["x-hub-version"]}&data={$content}";
        return strtoupper(md5($str));
    }


}