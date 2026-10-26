<?php

namespace App\Sdk\Active\CRM\CoachCore;

use App\Exceptions\Error;
use App\Library\WechatNotice\Chat;
use App\Sdk\Active\CRM\CoachCore\Auth\TokenCache;
use App\Sdk\Active\CrmHelper;
use Exception;
use GuzzleHttp\Client;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Throwable;

class Request
{

    private $baseUrl = '';
    private $tokenUrlReq = '';

    private $log;

    private $xCountryCode;

    private $enterpriseId;

    private $enterpriseConfig;

    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
        $this->enterpriseConfig = CrmHelper::getEnterpriseConfig($enterpriseId);

        // coach 和 sw 配置
        if (isset($this->enterpriseConfig['name'])) {

            if ($this->enterpriseConfig['name'] == 'coach') {
                $this->baseUrl = Config::getUrl();
                $this->tokenUrlReq = Config::getTokenUrl();
                $this->log = Config::LOG_NAME;
                $this->xCountryCode = 'CN';
            } elseif ($this->enterpriseConfig['name'] == 'sw') {
                $this->baseUrl = ConfigSW::getUrl();
                $this->tokenUrlReq = ConfigSW::getTokenUrl();
                $this->log = ConfigSW::LOG_NAME;
                $this->xCountryCode = 'SW';
            }
        }
    }


    /**
     * @param Base $request
     * @param string $requestType
     * @param bool $retry
     * @return Response
     * @throws Exception
     * @author: wangle@vchangyi.com
     * @Time: 2022/5/27
     */
    public function requestApi(Base $request, string $requestType = 'json', $retry = false): Response
    {
        $request->setHeader($this->enterpriseId, $this->xCountryCode);
        if ($request->getThrottle()) {
            $this->delayRequest($request);
        }
        $client = $this->initClient();

        $url = $this->setUrl($request);
        $param = $this->setOption($request, $requestType);
        $logParam = ['url' => $url, 'param' => $param, 'time' => date('Y-m-d H:i:s', time())];
        Log::info('请求参数', $logParam);
        CrmHelper::echoLine($this->log, '请求参数', $logParam);
        $requestStart = microtime(true);
        try {
            $response = $client->request('POST', $url, $param);
        } catch (Throwable $e) {
            $requestEnd = microtime(true);
            $subTime = $requestEnd - $requestStart;
            Chat::getInstance()->send('请求CRM接口出错', [
                'env' => config('app.env'),
                'api' => $url,
                'param' => $param,
                'error_info' => $e->getMessage(),
                "sub_time" => $subTime
            ]);
            throw new Error($e->getCode(), $e->getMessage() . "初次请求时间为" . $subTime);
        }
        if ($response->getStatusCode() == 401 && !$retry) {// 401 状态 token过期，做一次重试机制
            $this->clearTokenCache();
            return $this->requestApi($request, $requestType, true);
        }
        if ($requestType == 'xml') {
            $reqBody = simplexml_load_string($response->getBody()->getContents(), 'SimpleXMLElement', LIBXML_NOCDATA);
            $reqBodyXml = json_encode($reqBody);  //objecct转json
            $reqBodyArray = json_decode($reqBodyXml, true); //json转array
            Log::info('cdp-crm response', [
                'api' => $this->setUrl($request),
                'param' => $this->setOption($request),
                'response' => $reqBodyArray
            ]);
            CrmHelper::echoLine(Config::LOG_NAME, '返回参数', $reqBodyArray);
            $reqBodyArray = arrayRemoveByKey($reqBodyArray, "@attributes");
            return new Response($reqBodyArray);
        }
        $reqBody = json_decode($response->getBody(), true) ?? [];
        CrmHelper::echoLine($this->log, '返回参数', $reqBody);
        if (empty($reqBody)) {
            return new Response($reqBody);
        }
        Log::info('cdp-crm response', [
            'api' => $this->setUrl($request),
            'param' => $this->setOption($request),
            'response' => $reqBody
        ]);
        if ($this->isTokenExpiredResponse($reqBody) && !$retry) {
            Log::warning('cdp-crm token expired response, retry request', [
                'enterprise_id' => $this->enterpriseId,
                'api' => $this->setUrl($request),
                'response' => $reqBody
            ]);
            $this->clearTokenCache();
            return $this->requestApi($request, $requestType, true);
        }
        if ($response->getStatusCode() != 200) {
            $requestEnd = microtime(true);
            $subTime = $requestEnd - $requestStart;
            $reqBody['message'] = $reqBody['message'] . "执行总时长" . $subTime;
            Chat::getInstance()->send('请求CRM接口出错', [
                'env' => config('app.env'),
                'api' => $url,
                'param' => $param,
                'error_info' => $reqBody
            ]);
            $responseErr = $reqBody;
            throw new Error($responseErr['code'] ?? 10000, $responseErr['message']);
        }
        return new Response($reqBody);

    }

    private function isTokenExpiredResponse(array $response): bool
    {
        $code = (string)($response['code'] ?? $response['responseCode'] ?? '');
        $message = (string)($response['message'] ?? $response['msg'] ?? $response['responseDesc'] ?? '');
        $message = strtolower($message);

        if (in_array($code, ['401', '40101'], true)) {
            return true;
        }

        return strpos($message, 'token') !== false
            && (
                strpos($message, 'expired') !== false
                || strpos($message, 'invalid') !== false
                || strpos($message, 'not exists') !== false
                || strpos($message, 'not exist') !== false
            );
    }

    private function clearTokenCache(): void
    {
        new TokenCache($this->enterpriseId, true);
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

    private function setOption(Base $request, $requestType = 'json'): array
    {
        $body = $request->getQueryParam($requestType);
        $header = $request->getHeader();
        return array_merge($header, $body);
    }


    public function requestTokenApi(Base $request): Response
    {
        if ($request->getThrottle()) {
            $this->delayRequest($request);
        }
        $client = $this->initClient();
        $url = $this->setTokenUrl($request);
        $param = $this->setOption($request);
        $param['connect_timeout'] = 3;
        $param['timeout'] = 8;
        $logParam = ['url' => $url, 'param' => $param, 'time' => date('Y-m-d H:i:s', time())];
        Log::info('请求参数', $logParam);
        CrmHelper::echoLine($this->log, '请求参数', $logParam);

        try {
            $response = $client->request('POST', $url, $param);
        } catch (Throwable $e) {
            Chat::getInstance()->send('请求CRM获取Token接口出错', [
                'env' => config('app.env'),
                'api' => $url,
                'param' => $param,
                'error_info' => $e->getMessage()
            ]);
            throw new Error($e->getCode(), $e->getMessage());
        }
        $reqBody = json_decode($response->getBody(), true) ?? [];
        if (empty($reqBody)) {
            return new Response($reqBody);
        }
        Log::info('返回参数', $reqBody);
        Log::info('cdp-crm response', [
            'api' => $this->setUrl($request),
            'param' => $this->setOption($request),
            'response' => $reqBody
        ]);
        CrmHelper::echoLine($this->log, '返回参数', $reqBody);
        if ($response->getStatusCode() != 200) {
            Chat::getInstance()->send('请求CRM获取Token接口出错', [
                'env' => config('app.env'),
                'api' => $url,
                'param' => $param,
                'error_info' => $reqBody
            ]);
            $responseErr = $reqBody;
            throw new Error($responseErr['code'], $responseErr['message']);
        }

        return new Response($reqBody);
    }

    private function setTokenUrl(Base $request): string
    {
        return $this->tokenUrlReq . $this->getTokenApi($request);
    }

    private function getTokenApi(Base $request): string
    {
        if (config('app.env') == 'prod') {
            return $request->getTokenApiProd();
        }
        return $request->getTokenApiUat();
    }


}
