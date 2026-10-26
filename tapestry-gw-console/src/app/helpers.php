<?php

use App\Exceptions\Error;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use Illuminate\Support\Facades\Redis;
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
    function message($code, $args)
    {
        if (is_numeric($code)) {
            if (empty(config('code.' . $code)) && empty($args)) {
                $message = '未知的错误';
            } else if (!empty($args)) {
                $message = '%s';
            } else {
                $message = config('code.' . $code);
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

if (!function_exists('array_combine_by_key')) {
    /**
     * 根据数组指定的键对应的值, 作为新数组的键名
     *
     * @param array $arr 二维数组
     * @param string $key 键名
     *
     * @return array
     */
    function array_combine_by_key($arr, $key)
    {
        $keys = array_column($arr, $key);
        return array_combine($keys, $arr);
    }
}

if (!function_exists('formatPrice')) {
    /**
     * 格式化金额
     * @param $price
     * @return string
     */
    function formatPrice($price)
    {
        if (!is_numeric($price)) {
            $price = 0;
        }

        return number_format(bcdiv($price, 100, 2), 2);
    }
}

/**
 * redis批量删除
 *
 * @param $redisKey
 *
 * @return int
 */
if (!function_exists('getDownloadPath')) {
    function getDownloadPath()
    {
        return config('app.download_path');
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
 *  手机号脱敏
 */
if (!function_exists('hide_mobile')) {
    function hide_mobile($value)
    {
        return $value ? substr_replace($value, '****', 3, 4) : $value;
    }
}

if (!function_exists('fetchRepeatMemberInArray')) {

    function fetchRepeatMemberInArray($array): array
    {
        // 获取去掉重复数据的数组
        $uniqueArr = array_unique($array);
        // 获取重复数据的数组
        return array_diff_assoc($array, $uniqueArr);
    }
}