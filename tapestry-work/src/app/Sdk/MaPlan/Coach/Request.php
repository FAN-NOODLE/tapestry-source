<?php

namespace App\Sdk\MaPlan\Coach;

use App\Exceptions\Error;
use App\Library\WechatNotice\Chat;
use App\Sdk\MaPlan\Coach\Api\AccessToken;
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
     * @var $username //获取token的账号
     */
    protected $username;
    /**
     * @var $password //获取token的密码
     */
    protected $password;

    /**
     * @var Client 请求类
     */
    protected $request;

    /**
     * @var $enterpriseId //企业的id
     */
    protected $enterpriseId;
    /**
     * @var $header //请求接口带的token请求头
     */
    protected $header = [];

    protected $originType=1;

    public function __construct(int $enterpriseId,$originType)
    {
        $this->request = new Client();
        $this->enterpriseId = $enterpriseId;
        $this->originType = $originType;
    }

    public static function Instance($enterpriseId,$originType): Request
    {
        if (empty(self::$instance)) {
            self::$instance = new self($enterpriseId,$originType);
        }
        return self::$instance;

    }

    /**
     * 接口请求
     * @param $uri
     * @param $params
     * @param array $header
     * @return array|mixed
     * @throws GuzzleException
     */
    public function sendPost($uri, $params, $header = [])
    {
        Log::info("调用MA接口开始");
        $options = $header;
        if (!empty($params)) {
            $this->setHeader($header);
            $options = $this->setOptions($params);
        }
        $responseContent = "";
        try {
            $response = $this->request->request("POST", $uri, $options);
            $responseContent = $response->getBody()->getContents();
            if ($response->getStatusCode() != 200) {
                throw new Error($response->getStatusCode(), $responseContent);
            }
            unset($options['crowd_result']);
            Log::info("请求MA接口返回", ["url" => $uri, "request" => $options, "response" => $responseContent]);
        } catch (Exception $e) {
            unset($options['crowd_result']);
            Log::info("请求MA接口返回异常" . json_encode(["url" => $uri, "request" => $options, "response" => $e->getMessage()]));
            Chat::getInstance()->send('请求MA接口出错', [
                'env' => config('app.env'),
                'api' => $uri,
                'param' => $params,
                'error_info' => $e->getMessage()
            ]);
        }
        Log::info("调用MA接口结束");
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
        $body = $this->getQueryParam($params);
        $header = $this->getHeader();
        return array_merge($header, $body);
    }

    protected function getQueryParam($params, $requestType = 'json'): array
    {
        if ($requestType == 'json') {
            return ['body' => json_encode($params)];
        } else if ($requestType == 'form_data') {
            return ['form_params' => $params];
        } else if ($requestType == 'xml') {
            return ['body' => $params];
        } else {
            return ['data' => $params];
        }

    }

    /**
     * @return array
     */
    protected function getHeader(): array
    {
        return $this->header;
    }

    protected function setHeader($header)
    {
        $this->header = $header;
        if ($this->originType==1) {
            $accessToken = (new AccessToken($this->enterpriseId,$this->originType))->getAccessToken();
            $this->header = ["headers" => [
                "Authorization" => "Bearer {$accessToken}",
                "Content-Type" => "application/json",
                "Cache-Control" => "no-cache"
            ]];
        }
    }

    /**
     * 获取accessToken
     * @param $uri
     * @return array
     * @throws GuzzleException
     */
    public function getAccessToken($uri): array
    {
        $accountInfo = Config::getInfo();
        $credentials = base64_encode("{$accountInfo['username']}:{$accountInfo["password"]}");
        $headers["headers"] = array(
            "Authorization" => "Basic {$credentials}",
            "Content-Type" => "application/x-www-form-urlencoded",
            "Cache-Control" => "no-cache"
        );
        return $this->sendPost($uri, [], $headers);
    }
}