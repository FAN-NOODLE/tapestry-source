<?php

use App\Exceptions\Error;
use App\Models\Enterprise;
use App\Services\EnterpriseLogService;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use App\Models\EnterpriseOpenPlatform;

if (!function_exists('success')) {
    function success(array $data = [])
    {
        return response()->json($data);
    }
}

if (!function_exists('successV2')) {
    function successV2(array $data = [], $httpCode = 200)
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
    /**
     * @param $code
     * @param array $args
     * @param int $statusCode
     */
    function error($code, $args = [], $statusCode = 500)
    {
        $message = empty(config('code.' . $code)) ? '未知的错误' : config('code.' . $code);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        throw new Error($code, $message, $statusCode);
    }
}

if (!function_exists('message')) {
    function message($code)
    {
        return config('code.' . $code);
    }
}

if (!function_exists('switchAppHost')) {
    /**
     * 切换企业域名
     *
     * @param $idOrHash
     * @return array|bool
     */
    function switchAppHost($idOrHash)
    {
        if (empty($idOrHash)) {
            return false;
        }

        // 获取企业版本信息
        if (is_numeric($idOrHash)) {
            $enterprise = Enterprise::query()->where('id', '=', $idOrHash)->first();
        } else {
            $enterprise = Enterprise::query()->where('hash', '=', $idOrHash)->first();
        }

        if ($enterprise->version != Enterprise::NORMAL_VERSION) {
            config(['client.api_url.build' => 'http://49.235.249.11:5678']);
            // 获取网关域名
            $originConfig = config('client.api_url.shop');
            if (strpos($originConfig, 'gw-shop') == false) {
                $gwShop = config('client.api_url.shop') . '/gw-shop';
                $gwConsole = config('client.api_url.console') . '/gw-console';
                // 重写SHOP域名
                config(['client.api_url.shop' => $gwShop]);
                config(['client.api_url.console' => $gwConsole]);
            }

            if (in_array(getenv('APP_ENV'), ['dev', 'test'])) {
                EnterpriseLogService::debug('企业切换域名:' . $enterprise->id, config('client'));
            }
        }

        return $enterprise->toArray();
    }
}

if (!function_exists('multi_collect')) {
    /**
     * array 转 collect.
     *
     * @param $array
     * @return \Illuminate\Support\Collection
     */
    function multi_collect($array)
    {
        if (!is_array($array)) {
            return $array;
        }

        foreach ($array as $k => $v) {
            if (is_array($v)) {
                $array[$k] = multi_collect($v);
            }
        }

        return collect($array);
    }
}

/**
 * 提示性错误
 */
if (!function_exists('notice')) {
    function notice($code, $args = [], $data = [], $statusCode = 422, $previous = null, $headers = [])
    {
        $message = empty(config('code.' . $code)) ? '未知的错误' : config('code.' . $code);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        throw new Error($code, $message, $statusCode, $previous, $headers, $data);
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
 * 删除单个redis缓存
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

if (!function_exists('request')) {
    /**
     * Get an instance of the current request or an input item from the request.
     *
     * @param array|string $key
     * @param mixed $default
     * @return \Illuminate\Http\Request|string|array
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

if (!function_exists('getSql')) {
    // 获取执行sql
    function getSql($query)
    {
        $sql = str_replace("?", "'%s'", $query->toSql());

        return vsprintf($sql, $query->getBindings());
    }
}

if (!function_exists('getPaginateMeta')) {
    function getPaginateMeta($obj)
    {
        return [
            'meta' => [
                'total' => $obj->total(),
                'size'  => $obj->perPage(),
                'page'  => $obj->currentPage()
            ]
        ];
    }
}

if (!function_exists('arraySlice')) {
    /**
     * 分页函数
     * （效率高于array_slice函数）
     * array_slice底层实现：从第1个开始，一个个遍历，直到定位到开始位置，即偏移量（offset），然后取出后面limit个。
     * 这样随着offset的增大，遍历的次数也越来越大，遍历的次数是成线性增长的。
     * arraySlice是直接定位到开始位置offset，向后取limit个。
     *
     * @param array $list 数据列表
     * @param int $offset 开始位置
     * @param int $limit 每页数量
     * @return array
     */
    function arraySlice(array $list, int $offset, int $limit)
    {
        // 初始化待返数据
        $data = [];
        // 步长
        $step = 1;
        // 总数据量
        $count = count($list);
        for ($i = $offset; $i < $count; $i++) {
            if ($step <= $limit) {
                $data[] = $list[$i];
                $step++;
            }
        }

        return $data;
    }
}

/**
 * 统计redis键名管理
 *
 * @param       $keyPrefix string 前缀统一管理在 Common/RedisKey.php
 * @param array ...$args array 可变参数
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

if (!function_exists('codeRand')) {
    /**
     * 生成随机密码或随机验证码等
     *
     * @param int $length
     * @param string $type alpha | alpha-number | alpha-number-sign...
     * @return false|string
     */
    function codeRand($length = 6, $type = 'alpha-number')
    {
        $code_arr = array(
            'alpha'  => 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ',
            'number' => '0123456789',
            'sign'   => '#$%@*-_',
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

if (!function_exists('changeOpenPlatformConfig')) {
    function changeOpenPlatformConfig($enterpriseId, $openPlatformServer, $openPlatformId = 0, $initConfig = false)
    {
        Log::info('changeOpenPlatformConfig：' . $enterpriseId);
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
                    ->orderBy('id', 'DESC')
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
                ->first();
        }
        if (empty($config)) {
            error(620004);
        }

        Log::info('changeOpenPlatformConfig-2：' . var_export($config, true));

        $openPlatformServer->config['app_id'] = $config->app_id;
        $openPlatformServer->config['secret'] = $config->secret;
        $openPlatformServer->config['token'] = $config->token;
        $openPlatformServer->config['aes_key'] = $config->aes_key;
        return true;
    }
}
