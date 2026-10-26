<?php

namespace App\Sdk\CDP\Coach;

use App\Library\WechatNotice\Chat;
use Exception;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;


class Request
{
    /**
     * @var $instance //实例类本身 self
     */
    public static $instance;

    /**
     * @var $header //请求接口带的token请求头
     */
    protected $request;

    public function __construct()
    {
        $this->request = new Client();
    }

    public static function Instance(): Request
    {
        if (empty(self::$instance)) {
            self::$instance = new self();
        }
        return self::$instance;

    }

    /**
     * 接口请求
     *
     * @param string $apiPath
     * @param array $params
     * @return array|mixed
     * @throws
     */
    public function sendPost(string $apiPath, array $params)
    {
        $config = Config::getConfig();
        Log::info("调用CDP接口开始");
        $url = $config['domain'] . $apiPath;
        $postParams = $this->setOptions($params, $config['token']);
        try {
            $response = $this->request->request("POST", $url, $postParams);
            $responseContent = $response->getBody()->getContents();
            Log::info("请求CDP接口返回", ["url" => $url, "request" => $postParams, "response" => $responseContent]);
        } catch (Exception $e) {
            Log::info("请求CDP接口返回异常" . json_encode(["url" => $url, "request" => $postParams, "response" => $e->getMessage()]));
            Chat::getInstance()->send('请求CDP接口出错', [
                'env' => config('app.env'),
                'api' => $url,
                'param' => $postParams,
                'error_info' => $e->getMessage()
            ]);
        }
        Log::info("调用CDP接口结束");
        if (!empty($responseContent)) {
            return json_decode($responseContent, true) ?? [];
        }
        return [];
    }

    public function sendEasPost(string $url, array $payload, array $config): array
    {
        $postParams = RecommendationFeed::buildOptions($payload, $config['authorization']);
        $responseContent = '';
        Log::info('调用EAS推荐接口开始', ['url' => $url, 'request' => $payload]);
        try {
            $response = $this->request->request('POST', $url, $postParams);
            $responseContent = $response->getBody()->getContents();
            Log::info('调用EAS推荐接口返回', ['url' => $url, 'response' => $responseContent]);
        } catch (Exception $e) {
            Log::error('调用EAS推荐接口异常', ['url' => $url, 'error_info' => $e->getMessage()]);
        }
        Log::info('调用EAS推荐接口结束');
        if ($responseContent !== '') {
            return json_decode($responseContent, true) ?? [];
        }
        return [];
    }


    /**
     * @param $params
     * @param $token
     * @return array
     */
    protected function setOptions($params, $token): array
    {
        $header = [
            "headers" => ["token" => "REDACTED_FOR_SOURCE_SHARING" . $token,
                "Content-Type" => "application/json",
                "Cache-Control" => "no-cache"]
        ];
        return array_merge($header, ['body' => json_encode($params)]);
    }

}
