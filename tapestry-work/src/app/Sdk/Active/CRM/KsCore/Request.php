<?php

namespace App\Sdk\Active\CRM\KsCore;

use App\Exceptions\Error;
use App\Library\WechatNotice\Chat;
use App\Sdk\Active\CRM\KsCore\Auth\AuthSign;
use App\Sdk\Active\CrmHelper;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use App\Library\Util\Logger;

class Request
{
    private $baseUrl = '';

    public function __construct()
    {
        $this->baseUrl = Config::getUrl();
    }


    public function requestApi(Base $request, $method = "GET", $params = []): Response
    {
        if ($request->getThrottle()) {
            $this->delayRequest($request);
        }
        $client = $this->initClient();

        if ($method == 'POST') {
            $request->setType();

            $url                           = $this->setUrl($request);
            $param                         = $this->setOption($request);
            $param['form_params']['appid'] = Config::getAppId();
            // 拼接url
            $urlString = AuthSign::splicingString($param['form_params']);
            // 生成sign
            $sign = AuthSign::makeSign($param['form_params'], Config::getAppKey());

            $url .= $urlString . '&sign=' . $sign;
        } else {
            $url                      = $this->setUrl($request);
            $param                    = $this->setOption($request);
            $param['params']['appid'] = Config::getAppId();
            // 拼接url
            $urlString = AuthSign::splicingString($param['params']);
            // 生成sign
            $sign   = AuthSign::makeSign($param['params'], Config::getAppKey());
            $url    .= $urlString . '&sign=' . $sign;
            $params = $param;

        }
        $logParam = ['url' => $url, 'param' => $params, 'time' => date('Y-m-d H:i:s', time())];
        Logger::info('请求参数', $logParam);
        CrmHelper::echoLine(Config::LOG_NAME, '请求参数', $logParam);

        try {
            $response = $client->request(
                $method,
                $url,
                $params
            );

        } catch (\Throwable $e) {

            $logArr = [
                'url'    => $this->getApi($request),
                'params' => $request->getQueryParam(),
                'code'   => (int)$e->getCode(),
                'msg'    => $e->getMessage()
            ];
            Log::error($request->getErrorStr(), $logArr);

            throw new Error($e->getCode(), $e->getMessage());
        }

        $reqbody = json_decode($response->getBody(), true) ?? [];
        if (empty($reqbody)) {
            return new Response($reqbody);
        }
        Logger::info('返回参数', $reqbody);
        CrmHelper::echoLine(Config::LOG_NAME, '返回参数', $reqbody);

        Log::info('cdp-crm response', [
            'api'      => $this->setUrl($request),
            'param'    => $this->setOption($request),
            'response' => $reqbody
        ]);

        if ($response->getStatusCode() == 200) {
            return new Response($reqbody);
        } else {

            // 异常告警
            Chat::getInstance()->send('请求CRM接口出错', [
                'env'        => config('app.env'),
                'api'        => $url,
                'param'      => $param,
                'error_info' => $reqbody
            ]);
        }
        $responseErr = $reqbody;

        throw new Error($responseErr['code'], $responseErr['message']);

    }


    private function delayRequest(Base $request)
    {
        $key = md5(json_encode($request->getQueryParam()));
        if (Redis::connection()->exists($key)) {
            sleep(Redis::connection()->ttl($key));
//                throw new Error(10000,'请求API太频繁');
        }
        Redis::connection()->setex($key, 10, 1);
    }


    private function initClient(): Client
    {
        return new Client();
    }

    private function setUrl(Base $request): string
    {
        return $this->baseUrl . $this->getApi($request);
    }

    private function getApi(Base $request): string
    {
        if (config('app.env') == 'prod') {
            return $request->getApiProd();
        }
        return $request->getApiUat();
    }

    private function setOption(Base $request): array
    {
        $body   = $request->getQueryParam();
        $header = $request->getHeader();
        return array_merge($header, $body);
    }


}