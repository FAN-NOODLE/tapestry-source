<?php

declare(strict_types=1);

namespace App\Library;

use App\Traits\Singleton;
use Closure;
use GuzzleHttp\{HandlerStack, Middleware};
use GuzzleHttp\Client;
use GuzzleHttp\Exception\{ConnectException, RequestException};
use GuzzleHttp\Handler\CurlHandler;
use GuzzleHttp\Psr7\{Request, Response};
use Illuminate\Support\Facades\Log;
use Symfony\Component\HttpKernel\Exception\HttpException;

/**
 * guzzle统一请求类
 * Class GuzzleClient
 * @package App\Library
 */
class GuzzleClient
{
    use Singleton;

    /**
     * @var Client
     */
    private $guzzleClient;

    /**
     * GuzzleClient constructor.
     * @param string $host
     * @param array $middleware
     * @param int $timeout
     */
    public function __construct(string $host, $middleware = [], $timeout = 2)
    {
        // 创建 Handler
        $handlerStack = HandlerStack::create(new CurlHandler());
        if ($middleware) {
            foreach ($middleware as $name => $value) {
                $handlerStack->push($value, $name);
            }
        }
        // 创建重试中间件，指定决策者为 $this->retryDecider(),指定重试延迟为 $this->retryDelay()
        $handlerStack->push(Middleware::retry($this->retryDecider(), $this->retryDelay()));
        // 指定 handler
        $config             = [
            'handler'  => $handlerStack,
            'base_uri' => $host,
            'timeout'  => $timeout,
        ];
        $this->guzzleClient = new Client($config);
    }


    /**
     * 表单请求
     * @param $uri
     * @param array $data
     * @param array $options
     * @return string
     */
    public function requestPost(string $uri, array $data = [], array $options = []): string
    {
        try {
            $options['headers']     = [
                "Content-Type" => "application/x-www-form-urlencoded;charset=utf-8",
            ];
            $options['form_params'] = $data;
            return $this->guzzleClient->post($uri, $options)->getBody()->getContents();
        } catch (RequestException $exception) {
            return $exception->getResponse()->getBody()->getContents();
        }
    }

    /**
     * json请求
     * @param string $uri
     * @param $body
     * @param bool $encode
     * @param array $option
     * @return string
     */
    public function requestPostJson(string $uri, $body, bool $encode = true, array $option = []): string
    {
        return $this->request("post", $uri, $body, $encode, $option);
    }

    /**
     * @param string $uri
     * @param $body
     * @param bool $encode
     * @param array $option
     *
     * @return string
     *
     * @author 邱阿朋 <qiuapeng@vchangyi.com>
     * @date   2021/01/14 17:26
     */
    public function requestPatchJson(string $uri, $body, bool $encode = true, array $option = []): string
    {
        return $this->request("patch", $uri, $body, $encode, $option);
    }

    /**
     * get请求
     * @param $uri
     * @param array $options
     * @return string
     */
    public function requestGet($uri, $options = []): string
    {
        try {
            return $this->guzzleClient->get($uri, $options)->getBody()->getContents();
        } catch (RequestException $exception) {
            return $exception->getResponse()->getBody()->getContents();
        }
    }

    public function request(string $method, string $uri, $body, bool $encode = true, array $option = []): string
    {
        try {
            $header = [
                "Content-Type" => "application/json",
                "Accept"       => "application/json",
            ];
            if ($option && isset($option['headers'])) {
                $header = array_merge($option['headers'], $header);
            }
            $body               = $encode ? json_encode($body) : $body;
            $options['headers'] = $header;
            $options['body']    = $body;
            if ($method == "post") {
                return $this->guzzleClient->post($uri, $options)->getBody()->getContents();
            }
            if ($method == "patch") {
                return $this->guzzleClient->patch($uri, $options)->getBody()->getContents();
            }
            throw new HttpException(\Illuminate\Http\Response::HTTP_METHOD_NOT_ALLOWED, "请求方式错误");
        } catch (RequestException $exception) {
            return $exception->getResponse()->getBody()->getContents();
        }
    }

    /**
     * 返回一个匿名函数, 匿名函数若返回false 表示不重试，反之则表示继续重试
     * @return Closure
     * @author  邱阿朋 <qiuapeng@vchangyi.com>
     * @date    2020/6/9 12:44
     */
    private function retryDecider(): Closure
    {
        return function (
            $retries,
            Request $request,
            Response $response = null,
            RequestException $exception = null
        ) {
            // 超过最大重试次数，不再重试
            if ($retries >= 5) {
                Log::debug("超过最大重试次数,不再重试");
                return false;
            }
            // 请求失败，继续重试
            if ($exception instanceof ConnectException) {
                Log::debug("请求失败,继续重试{$request->getUri()}", [$exception]);
                return true;
            }
            if ($response) {
                // 如果请求有响应，但是状态码大于等于500，继续重试(这里根据自己的业务而定)
                if ($response->getStatusCode() >= 500) {
                    Log::debug("请求有响应,但是状态码大于等于500,继续重试");
                    return true;
                }
            }
            return false;
        };
    }


    /**
     * 返回一个匿名函数，该匿名函数返回下次重试的时间（毫秒）
     * @return Closure
     * @author  邱阿朋 <qiuapeng@vchangyi.com>
     * @date    2020/6/9 12:45
     */
    private function retryDelay(): Closure
    {
        return function ($numberOfRetries) {
            return 1000 * $numberOfRetries;
        };
    }
}