<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Exceptions\Error;
use App\Library\Util\Logger;
use App\Library\WechatNotice\Chat;
use App\Sdk\Active\CRM\CoachCore\Auth\AuthSign;
use App\Sdk\Active\CrmHelper;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Throwable;

class RequestCore
{
    private $baseUrl = '';

    private $enterpriseId;

    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
        $this->baseUrl      = Config::getUrlCore();
    }


    /**
     * @param Base $request
     *
     * @return ResponseCore
     * @author: wangle@vchangyi.com
     * @Time: 2022/4/26
     */
    public function requestApiCore(Base $request): ResponseCore
    {
        $request->setHeaderCore($this->enterpriseId);
        if ($request->getThrottle()) {
            $this->delayRequest($request);
        }
        $client = $this->initClient();

        $url   = $this->setUrl($request);
        $param = $this->setOption($request);

        // 组装请求参数
        $body['Data']     = $param['data'];
        $body['Noncestr'] = Config::NONCESTR;
        $body['Ticks']    = time();
        // 生成sign
        $Signature = AuthSign::makeSign($body, Config::SecretValue);

        $body['Signature'] = $Signature;
        $params['body']    = json_encode($body);
        $params['headers'] = $param['headers'];

        $logParam = ['url' => $url, 'param' => $params, 'time' => date('Y-m-d H:i:s', time())];
        Logger::info('请求参数', $logParam);
        CrmHelper::echoLine(Config::LOG_NAME, '请求参数', $logParam);

        try {
            $response = $client->request('POST', $url, $params);
        } catch (Throwable $e) {
            $logArr = [
                'url'    => $this->getApi($request),
                'params' => $params,
                'code'   => (int)$e->getCode(),
                'msg'    => $e->getMessage()
            ];
            Log::error($request->getErrorStr(), $logArr);
            Chat::getInstance()->send('请求CRM接口出错', [
                'env'        => config('app.env'),
                'api'        => $url,
                'url'        => $this->getApi($request),
                'param'      => $param,
                'error_info' => $e->getMessage(),
            ]);
            throw new Error($e->getCode(), $e->getMessage());
        }

        $reqBody = json_decode($response->getBody(), true) ?? [];
        if (empty($reqBody)) {
            return new ResponseCore($reqBody);
        }

        Logger::info('返回参数', $reqBody);
        CrmHelper::echoLine(Config::LOG_NAME, '返回参数', $reqBody);

        Log::info('cdp-crm response', [
            'api'      => $this->setUrl($request),
            'param'    => $this->setOption($request),
            'response' => $reqBody
        ]);
        if ($response->getStatusCode() == 200) {
            return new ResponseCore($reqBody);
        } else {

            // 异常告警
            Chat::getInstance()->send('请求CRM接口出错', [
                'env'        => config('app.env'),
                'api'        => $url,
                'param'      => $param,
                'error_info' => $reqBody
            ]);

        }
        $responseErr = $reqBody;

        throw new Error($responseErr['code'], $responseErr['message']);

    }


    private function delayRequest(Base $request)
    {
        $key = $this->enterpriseId . ':' . md5(json_encode($request->getQueryParam()));
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
        $body   = $request->handleQueryParam();
        $header = $request->getHeader();
        return array_merge($header, $body);
    }


}