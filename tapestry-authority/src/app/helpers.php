<?php

use App\Exceptions\Error;
use GuzzleHttp\Client;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Redis;
use Illuminate\Validation\ValidationException;
use App\Models\Enterprise;
use App\Models\EnterpriseOpenPlatform;

if (!function_exists('success')) {
    function success(array $data = [], $httpCode = 200)
    {
        $requestId = $_SERVER['x_request_id'] ?? LOG_ID;

        return response()->json([
            'code' => 200,
            'message' => 'success',
            'request_id' => $requestId,
            'data' => !empty($data) ? $data : (object)$data
        ], $httpCode);
    }
}

if (!function_exists('error')) {
    function error($code, $args = [], $data = [], $httpCode = 500)
    {
        list($code, $message) = message($code, $args);

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
    function notice($code, $args = [], $data = [], $httpCode = 422)
    {
        list($code, $message) = message($code, $args);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        $e = new Error($code, $message, $data, $httpCode);

        $trace = current(debug_backtrace());
        $line  = $trace['line'];
        $file  = $trace['file'];

        $ref     = new \ReflectionClass($e);
        $linePro = $ref->getProperty('line');
        $linePro->setAccessible(true);
        $linePro->setValue($e, $line);

        $filePro = $ref->getProperty('file');
        $filePro->setAccessible(true);
        $filePro->setValue($e, $file);
        throw $e;
    }
}

if (!function_exists('validate')) {
    /**
     * @param $params
     * @param $rules
     * @param $message
     *
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
    function message($code, $args)
    {
        if (is_numeric($code)) {
            $message = empty(config('code.' . $code)) ? '未知的错误' : config('code.' . $code);
        } else {
            $config = config('notice.' . $code);
            $arr = explode('|', $config);
            $message = $arr[1];
            $code = $arr[0];
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
     *
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

if (!function_exists('redis')) {
    /**
     * redis实例
     *
     * @param null $name
     * @return \Illuminate\Redis\Connections\Connection
     */
    function redis($name = null)
    {
        return Redis::connection($name);
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
        $url = config('client.api_url.uid');
        $client = (new \GuzzleHttp\Client())->get($url, [
            'http_errors' => false,
            'verify' => false,
        ])->getBody()->getContents();

        return (int)$client;
    }
}

if (!function_exists('codeRand')) {
    /**
     * 生成随机密码或随机验证码等
     *
     * @param int $length
     * @param string $type alpha | alpha-number | alpha-number-sign...
     *
     * @return false|string
     */
    function codeRand($length = 6, $type = 'alpha-number')
    {
        $code_arr = array(
            'alpha' => 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ',
            'number' => '0123456789',
            'sign' => '#$%@*-_',
        );

        $type_arr = explode('-', $type);

        $chars = '';

        foreach ($type_arr as $t) {
            if (!array_key_exists($t, $code_arr)) {
                trigger_error("Can not generate type ($t) code");
            }
            $chars .= $code_arr[$t];
        }

        $chars = str_shuffle($chars);
        $number = $length > strlen($chars) - 1 ? strlen($chars) - 1 : $length;
        return substr($chars, 0, $number);
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
        $client = new Client();
        $uriEnv = env('SHORT_LINKS_URI', 'http://micro-track.vchangyi.com/micro-tracker/api/v1/trace/');
        $baseUri = $uriEnv . $traceId;

        $args = [];
        foreach ($params as $key => $val) {
            $args[] = [
                'name' => (string)$key,
                'value' => (string)$val
            ];
        }

        try {
            $headers['Content-Type'] = 'application/json';
            if (env('APP_ENV') == 'prod') {
                $headers['Accept'] = 'application/json';
//                $headers['Host']   = 'crs-api.vchangyi.com';

            }
            $result = $client->request('POST', $baseUri, [
                'headers' => $headers,
                'body' => json_encode(['args' => $args]),
                'http_errors' => false
            ]);
            $content = json_decode($result->getBody(), true);

            return $content['id'];
        } catch (\Exception $e) {
            throw new Error($e->getCode(), '短链接加密失败');
        }
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
        $client = new Client();
        $uriEnv = env('SHORT_LINKS_URI', 'http://micro-track.vchangyi.com/micro-tracker/api/v1/trace/');
        $baseUri = $uriEnv . $traceId;

        $headers['Content-Type'] = 'application/json';
        if (env('APP_ENV') == 'prod') {
            $headers['Accept'] = 'application/json';
            $headers['Host'] = 'crs-api.vchangyi.com';

        }
        try {
            $result = $client->request('GET', $baseUri, [
                'headers' => $headers,
                'http_errors' => false
            ]);
            return json_decode($result->getBody(), true);
        } catch (\Exception $e) {
            throw new Error($e->getCode(), '短链接解密失败');
        }
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
        $keyPrefix = $reflection->getConstant($keyPrefix);

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

/**
 * 删除多个redis缓存
 * 传说 redis 对 :  支持好像更好
 */
if (!function_exists('delRedis')) {
    function delRedis($keys)
    {
        if (!is_array($keys) || in_array('*', $keys)) {
            return false;
        }
        $redis = Redis::connection();
        $redis->del($keys);
        return true;
    }
}

if (!function_exists('exportToCsv')) {
    function exportToCsv($fileName = '', $headArr = [], $data = [])
    {
        //文件夹
        $filedir = storage_path() . '/';
        if (!is_dir($filedir)) mkdir($filedir, 777, true);
        //路径
        $fileName .= '.csv';
        $filePath = $filedir . $fileName;

        //生成csv文件
        $fp = fopen($filePath, 'w');

        $header = implode(',', $headArr) . PHP_EOL;
        $content = '';
        foreach ($data as $line => $val) {
            $content .= implode(',', $val) . PHP_EOL;
        }
        $csv = $header . $content;
        // 写入数据
        fwrite($fp, $csv);
        //关闭文件句柄
        fclose($fp);
        return;
    }
}

if (!function_exists('redis')) {
    /**
     * redis实例
     *
     * @param null $name
     * @return \Illuminate\Redis\Connections\Connection
     */
    function redis($name = null)
    {
        return Redis::connection($name);
    }
}

if (!function_exists('changeOpenPlatformConfig')) {
    function changeOpenPlatformConfig($enterpriseId, $openPlatformServer, $openPlatformId = 0, $initConfig = false)
    {
        $config = false;
        if ($openPlatformId == 0) {
            $enterprise = Enterprise::query()
                ->where('id', $enterpriseId)
                ->get(['open_platform_id'])
                ->first();
            if (empty($enterprise)) {
                error(600002);
            }
            if ($enterprise->open_platform_id == 0) {
                if (!$initConfig) {
                    error(620004);
                }
                $config = EnterpriseOpenPlatform::query()
                    ->orderBy('created_at', 'DESC')
                    ->first();
                // 关联配置
                Enterprise::query()->where('id', $enterpriseId)->update(['open_platform_id' => $config->id]);
            } else {
                $openPlatformId = $enterprise->open_platform_id;
            }
        }
        if (!$config && $openPlatformId != 0) {
            $config = EnterpriseOpenPlatform::query()
                ->where('id', $openPlatformId)
                ->orderBy('created_at', 'DESC')
                ->first();
        }
        if (empty($config)) {
            error(620004);
        }


        $openPlatformServer->config['app_id'] = $config->app_id;
        $openPlatformServer->config['secret'] = $config->secret;
        $openPlatformServer->config['token'] = $config->token;
        $openPlatformServer->config['aes_key'] = $config->aes_key;
        return true;
    }
}