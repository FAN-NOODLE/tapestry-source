<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/9 0009
 * Time: 10:24
 */

namespace App\Library\WechatNotice;


use Closure;
use GuzzleHttp\Client;
use GuzzleHttp\Handler\CurlHandler;
use GuzzleHttp\HandlerStack;
use GuzzleHttp\Middleware;
use GuzzleHttp\Psr7\Request;
use GuzzleHttp\Psr7\Response;
use GuzzleHttp\Exception\ConnectException;
use GuzzleHttp\Exception\RequestException;


class ClientFactory
{
    public static function create($host,$timeout=2):Client
    {
        $handlerStack = HandlerStack::create(new CurlHandler());
        $handlerStack->push(Middleware::retry(static::retryDecider(), static::retryDelay()));
        // 指定 handler
        $config = [
            'handler'  => $handlerStack,
            'base_uri' => $host,
            'timeout'  => $timeout,
        ];
        return new Client($config);
    }

    private static function retryDecider():Closure
    {
        return function (
            $retries,
            Request $request,
            Response $response = null,
            RequestException $exception = null
        ){
            if ($retries >= 5) {
                return false;
            }
            // 请求失败，继续重试
            if ($exception instanceof ConnectException) {
                return true;
            }
            if ($response) {
                // 如果请求有响应，但是状态码大于等于500，继续重试(这里根据自己的业务而定)
                if ($response->getStatusCode() >= 500) {
                    return true;
                }
            }
            return false;
        };
    }

    private static function retryDelay():Closure
    {
        return function ($numberOfRetries) {
            return 1000 * $numberOfRetries;
        };
    }

}