<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/12 0012
 * Time: 15:52
 */

namespace App\Library\DataLake;


use App\Exceptions\Error;
use App\Library\WechatNotice\Chat;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\ClientException;
use GuzzleHttp\Exception\ServerException;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class Request
{
    private $baseUrl = 'https://dl-api.lorealchina.com';

    private $printEcho = true;

    private $enterpriseId;

    function __construct($enterpriseId = '')
    {
        $this->enterpriseId = $enterpriseId;
    }

    public function setEnterpriseId($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
    }

    /**
     * @param Base $request
     * @return Response
     * @throws \Throwable
     */
    public function requestApi(Base $request)
    {
        return new Response([]);
//        if ($request->getThrottle()) {
//            $this->delayRequest($request);
//        }
        $client = $this->initClient();
        $request->notifyEnterpriseId($this->enterpriseId);
        try {

            $response = $client->request(
                'POST',
                $request->getFullUrl(),
                $request->getOptions()
            );

        } catch (ClientException|ServerException $guzzleError) {
            Chat::getInstance()->send('请求datalake接口出错', [
                'api'   => $request->getFullUrl(),
                'param' => $request->getOptions(),
                'error_info' => [
                    'code'    => $guzzleError->getCode(),
                    'msg'     => $guzzleError->getMessage(),
                    'context' => $guzzleError->getTraceAsString()
                ]
            ]);
            throw  $guzzleError;
        } finally {
            if (isset($response) && $response->getBody()) {
                $responseContent = json_decode($response->getBody(), true);
            } else {
                $responseContent = [];
            }
            Log::info('loreal crm response', [
                'api'      => $request->getFullUrl(),
                'param'    => $request->getOptions(),
                'response' => $responseContent
            ]);
            $this->echoLine($this->getErrorContext($request, $responseContent));
        };
        return new Response($responseContent);
    }


    private function delayRequest(Base $request)
    {
        $key = md5(json_encode($request->getOptions()));
        if (Redis::connection()->exists($key)) {
            sleep(Redis::connection()->ttl($key));
        }
        Redis::connection()->setex($key, 10, 1);
    }

    private function initClient()
    {
        return new Client();
    }


    /**
     * 测试用函数,打印请求和返回
     *
     * @param string $title
     * @param array $context
     */
    public function echoLine($context)
    {
        if ($this->printEcho && config('app.env') == 'dev') {
            $str = PHP_EOL;
            $str .= var_export($context, true);
            $str .= PHP_EOL;
            file_put_contents(__DIR__ . '/../dataLake_log.txt', $str, FILE_APPEND);
        }
    }

    public function getErrorContext(Base $request, $responseContext = [])
    {
        return [
            'api'      => $request->getFullUrl(),
            'param'    => $request->getOptions(),
            'response' => $responseContext
        ];
    }

}