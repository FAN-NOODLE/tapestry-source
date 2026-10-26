<?php

use App\Exceptions\Error;
use GuzzleHttp\Client;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Redis;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;

if (!function_exists('success')) {
    function success(array $data = [], $httpCode = 200)
    {
        return response()->json([
            'code'       => 200,
            'message'    => 'success',
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => !empty($data) ? $data : (object)$data
        ], $httpCode);
    }
}

if (!function_exists('error')) {
    function error($code, $msg = '', $args = [], $data = [], $httpCode = 200)
    {
        list($code, $message) = message($code, $args, $msg);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        throw new Error($code, $message, $data, $httpCode);
    }
}

/**
 * 提示性错误
 */
if (!function_exists('notice')) {
    function notice($code, $args = [], $data = [], $httpCode = 200)
    {
        list($code, $message) = message($code, $args);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        throw new Error($code, $message, $data, $httpCode);
    }
}

if (!function_exists('validate')) {
    /**
     * @param $params
     * @param $rules
     * @param $message
     * @throws ValidationException
     */
    function validate($params, $rules, $message)
    {
        $validator = Validator::make($params, $rules, $message);
        if ($validator->fails()) {
            throw new ValidationException($validator);
        }
    }
}

if (!function_exists('message')) {
    function message($code, $args, $msg = '')
    {
        if (is_numeric($code)) {
            $message = '未知的错误';
            if (!empty(config('code.' . $code)))   {
                $message = config('code.' . $code);
            } elseif (!empty($msg)) {
                $message = $msg;
            } else {
                $message = '未知的错误';
            }
        } else {
            $config  = config('notice.' . $code);
            $arr     = explode('|', $config);
            $message = $arr[1];
            $code    = $arr[0];
        }

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        return [$code, $message];
    }
}


if (!function_exists('request')) {
    /**
     * Get an instance of the current request or an input item from the request.
     *
     * @param array|string $key
     * @param mixed $default
     * @return Request|string|array
     */
    function request($key = null, $default = null)
    {
        if (is_null($key)) {
            return app('request');
        }

        if (is_array($key)) {
            return app('request')->only($key);
        }

        $value = app('request')->__get($key);

        return is_null($value) ? value($default) : $value;
    }
}

/**
 * 根据数组指定的键对应的值, 作为新数组的键名
 *
 * @param array $arr 二维数组
 * @param string $key 键名
 *
 * @return array
 */

if (!function_exists('array_combine_by_key')) {
    function array_combine_by_key($arr, $key)
    {
        $keys = array_column($arr, $key);
        return array_combine($keys, $arr);
    }
}

if (!function_exists('idGenerate')) {
    function idGenerate()
    {
        $url    = config('client.api_url.uid');
        $client = (new \GuzzleHttp\Client())->get($url, [
            'http_errors' => false
        ])->getBody()->getContents();

        return (int)$client;
    }
}

/**
 * 时间差计算
 *
 * @param float $startTime
 * @param float $endTime
 *
 * 返回 单位为秒(s)
 *
 * @return float
 */
if (!function_exists('exec_time')) {
    function exec_time(float $startTime, float $endTime = 0) : float
    {
        if ($endTime == 0) {
            $endTime = microtime(true);
        }

        $execTime = $endTime - $startTime;
        return round($execTime, 2);
    }
}

/**
 * 统计redis键名管理
 *
 * @param       $keyPrefix string 前缀统一管理在 Common/RedisKey.php
 * @param array ...$args array 可变参数
 *
 * @return string
 */
if (!function_exists('redisKey')) {

    function redisKey($keyPrefix, ...$args)
    {
        $reflection = new ReflectionClass(\App\Common\RedisKey::class);
        $keyPrefix  = $reflection->getConstant($keyPrefix);

        if (empty($keyPrefix)) {
            error('900002');
        }
        $ext = '';

        if (!empty($args)) {

            $result = [];
            array_walk_recursive($args, function ($value, $key) use (&$result) {
                if (!is_numeric($key)) {
                    $value = $key . '_' . $value;
                }
                array_push($result, $value);
            });

            $ext .= implode(':', $result);

        }
        return rtrim(rtrim($keyPrefix, ':') . ':' . $ext, ':');
    }
}

/**
 * 从redis中获取或存储数据
 *
 * @param      $key
 * @param null $data
 *
 * @return bool
 */
if (!function_exists('getOrSetData')) {
    function getOrSetData($key, $data = null, $expire = 3600)
    {
        if ($data) {
            if ($expire == null) {
                Redis::set($key, json_encode($data));
            } else {
                Redis::setex($key, $expire, json_encode($data));
            }
            return true;
        }

        $res = Redis::get($key);

        if ($res !== null) {
            return json_decode($res, true);
        }
        return false;
    }
}

if (!function_exists('delRedisKey')) {
    function delRedisKey($redisKey)
    {
        // 删除
        Redis::del($redisKey);
    }
}

/**
 * 短链接解密
 *
 * @return int|string
 */
if (!function_exists('decryptShortLink')) {
    function decryptShortLink($traceId = 0)
    {
        $client  = new \GuzzleHttp\Client();
        $uriEnv  = env('SHORT_LINKS_URI', 'http://micro-track.vchangyi.com/micro-tracker/api/v1/trace/');
        $baseUri = $uriEnv . $traceId;

        try {
            $result = $client->request('GET', $baseUri, [
                'headers'     => [
                    'Accept'       => 'application/json',
                    'Content-Type' => 'application/json',
                ],
                'http_errors' => false
            ]);
            return json_decode($result->getBody(), true);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), '短链接解密失败');
        }
    }
}

/**
 * 短链接解密兼容不同路径
 *
 * @return int|string
 */
if (!function_exists('decryptShortLinkByTest')) {
    function decryptShortLinkByTest($traceId = 0)
    {
        $client  = new \GuzzleHttp\Client();
        $uriEnv  = 'http://micro-track.vchangyi.com/micro-tracker/api/v1/trace/';
        $baseUri = $uriEnv . $traceId;

        try {
            $result = $client->request('GET', $baseUri, [
                'headers'     => [
                    'Accept'       => 'application/json',
                    'Content-Type' => 'application/json',
                ],
                'http_errors' => false
            ]);
            return json_decode($result->getBody(), true);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), '短链接解密失败');
        }
    }
}

/**
 * [std_class_object_to_array 将对象转成数组]
 * @param [stdclass] $stdclassobject [对象]
 * @return [array] [数组]
 */

if (!function_exists('objToArray')) {
    function objToArray($stdclassobject)
    {
        $_array = is_object($stdclassobject) ? get_object_vars($stdclassobject) : $stdclassobject;
        foreach ($_array as $key => $value) {
            $value       = (is_array($value) || is_object($value)) ? objToArray($value) : $value;
            $array[$key] = $value;
        }

        return $array;
    }
}
if (!function_exists('curl_get_https')) {
    function curl_get_https($url){
        $curl = curl_init(); // 启动一个CURL会话
        curl_setopt($curl, CURLOPT_URL, $url);
        curl_setopt($curl, CURLOPT_HEADER, 0);
        curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);
        curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, false); // 跳过证书检查
        curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, false);  // 从证书中检查SSL加密算法是否存在
        $tmpInfo = curl_exec($curl);     //返回api的json对象
        //关闭URL请求
        curl_close($curl);
        return $tmpInfo;    //返回json对象
    }
}

/**
 * 短链接加密
 *
 * @return int|string
 */
if (!function_exists('encryptShortLink')) {
    function encryptShortLink($params = [], $traceId = 0)
    {
        $client  = new Client();
        $uriEnv  = env('SHORT_LINKS_URI', 'http://micro-track.vchangyi.com/micro-tracker/api/v1/trace/');
        $baseUri = $uriEnv . $traceId;

        $args = [];
        foreach ($params as $key => $val) {
            if (is_array($val)) {
                $val = json_encode($val);
            }
            $args[] = [
                'name'  => (string)$key,
                'value' => (string)$val
            ];
        }

        try {
            $result  = $client->request('POST', $baseUri, [
                'headers'     => [
                    'Accept'       => 'application/json',
                    'Content-Type' => 'application/json',
                ],
                'body'        => json_encode(['args' => $args]),
                'http_errors' => false
            ]);
            $content = json_decode($result->getBody(), true);

            return $content['id'];
        } catch (\Exception $e) {
            throw new Error($e->getCode(), '短链接加密失败');
        }
    }
}

if (!function_exists('successJson')) {
    function successJson(array $data = [], $httpCode = 200)
    {
        return response()->json([
            'code'       => 200,
            'message'    => 'success',
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => $data
        ], $httpCode);
    }
}
