<?php


namespace App\Library\Ddc;


use App\Exceptions\Error;
use GuzzleHttp\Client;
use GuzzleHttp\Exception\RequestException;
use Illuminate\Support\Facades\Log;

class DdcMessage
{

    const PLATFORM = 'isg';

    const GUIDE_USER_BINDING_ADD = 'guide.user_binding_add';
    const GUIDE_USER_BINDING_UPDATE = 'guide.user_binding_update';
    const GUIDE_USER_BINDING_DELETE = 'guide.user_binding_delete';

    /**
     * @var DdcConfig
     */
    private $config;

    protected $client;


    /**
     * DDCMessage constructor.
     */
    public function __construct()
    {
        $this->config = DdcConfig::getInstance();
        $this->client = new Client();
    }

    public function push($enterpriseHash, $method, $data)
    {
        $params = [
            'method' => $method,
            'platform' => self::PLATFORM,
            'data' => $data
        ];

        $headers = [
            'appKey' => $this->config->getAuthKey(),
            'sign' => $this->generateSignature($params),
            'timestamp' => time(),
            'enterprise-hash' => $enterpriseHash,
            'x-request-id' => $_SERVER['x_request_id'] ?? md5(uniqid()) . rand(10000, 99999)
        ];
        $this->post($this->config->getHost() . '/micro-ddc/message/push', $params, $headers);
        return true;
    }

    protected function generateSignature(array $params)
    {
        //多维数组内排序
        $sortFunc = static function (array $params, callable $sortFunc): array {
            array_walk($params, static function(&$item) use ($sortFunc) {
                $item = is_array($item) ? $sortFunc($item, $sortFunc) : $item;
            });
            ksort($params);
            return $params;
        };
        $paramJsonString = json_encode($sortFunc($params, $sortFunc), JSON_HEX_AMP + JSON_UNESCAPED_UNICODE + JSON_UNESCAPED_SLASHES);
        return md5($this->config->getAuthSecret() . $paramJsonString . $this->config->getAuthSecret());
    }

    protected function post($url, $body = [], $headers = [])
    {
        Log::info("DDC request start", ['headers' => $headers, 'body' => $body]);
        $options = [
            'headers' => $headers,
            'body' => empty($body) ? '{}' : json_encode($body)
        ];
        $startTime = microtime(true);
        $responseData = null;
        try {
            $response = $this->client->post($url, $options);
            $responseData = $response->getBody();
            Log::info("DDC request success", [
                'httpCode' => $response->getStatusCode(),
                'responseData' => (string) $responseData
            ]);
        } catch (RequestException $e) {
            Log::error("DDC request exception", [
                'code' => $e->getCode(),
                'msg' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw new Error(1000, "DDC请求失败", 500, $e);
        } finally {
            $endTime = microtime(true);
            Log::info("DDC request time", [
                'startTime' => $startTime,
                'endTime' => $endTime,
                'duration' => sprintf("%.3f", $endTime - $startTime),
            ]);
        }

        return $this->getResponseData($responseData);
    }

    /**
     * 处理http返回结果
     * @param string $responseData
     * @return mixed
     */
    private function getResponseData(string $responseData)
    {
        $result = json_decode($responseData, true);
        if (!isset($result['code'])) {
            // 抛出异常
            throw new Error(100000, "返回结果code参数不存在");
        }
        if (200 != $result['code']) {
            // 抛出异常
            $message = $result['code'] . ":" . $result['message'];
            throw new Error(100000, $message);
        }
        return $result;
    }
}