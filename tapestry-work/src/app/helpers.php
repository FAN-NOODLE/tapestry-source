<?php

use GuzzleHttp\Client;
use App\Exceptions\Error;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use App\Services\EnterpriseService;
use App\Common\Constant;

if (!function_exists('success')) {
    function success(array $data = [])
    {
        $response = $data;
        Log::debug('Response Success: ' . json_encode($response, JSON_UNESCAPED_UNICODE));

        return $response;
    }
}

if (!function_exists('successJson')) {
    function successJson(array $data = [], $httpCode = 200)
    {
        return response()->json([
            'code'       => 200,
            'message'    => 'success',
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => !empty($data) ? $data : (object)$data
        ], $httpCode);
    }
}

//if (! function_exists('error')) {
//    function error($code, $args = [], $statusCode = 500)
//    {
//        $message = empty(config('code.'.$code)) ? '未知的错误' : config('code.'.$code);
//
//        if (!empty($args)) {
//            $message = vsprintf($message, $args);
//        }
//
//        throw new Error($code, $message, $statusCode);
//
//    }
//}


/**
 * 服务器抛错
 */
if (!function_exists('error')) {
    function error($codeOrMessage, $args = [], $httpStatus = 500)
    {
        if (is_numeric($codeOrMessage)) {
            $message = empty(config('code.' . $codeOrMessage)) ? '未知的错误' : config('code.' . $codeOrMessage);

            if (!empty($args)) {
                $message = vsprintf($message, $args);
            }

            throw new Error($codeOrMessage, $message, $httpStatus);

        } else {
            $e = new Error(100000, $codeOrMessage, $httpStatus);

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
}

if (!function_exists('message')) {
    function message($code)
    {
        return config('code.' . $code);
    }
}

if (!function_exists('paginates')) {
    function paginates($paginator)
    {
        return [
            'data' => $paginator->items(),
            'meta' => [
                'total' => $paginator->total(),
                'page'  => $paginator->currentPage(),
                'size'  => (int)$paginator->perPage(),
            ]
        ];
    }
}

if (!function_exists('request')) {
    /**
     * Get an instance of the current request or an input item from the request.
     *
     * @param array|string $key
     * @param mixed        $default
     *
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

/**
 * 最近几个月
 *
 * @param int $month
 *
 * @return array
 */
if (!function_exists('getMonth')) {
    function getMonth($month = 6)
    {
        $month--;
        $arr = array();
        for ($i = 0; $i <= $month; ++$i) {
            $t     = strtotime("-$i month");
            $arr[] = [
                'date'  => date('m', $t),
                'start' => date('Y-m-01', $t),
                'end'   => date('Y-m-t', $t)
            ];

        }
        return array_reverse($arr);
    }
}

/**
 * 获取前几个周的周一开始日期
 *
 * @param int $month
 *
 * @return array
 */
if (!function_exists('getWeek')) {
    function getWeek($week = 6)
    {
        $arr = array();
        for ($i = 1; $i <= $week; ++$i) {
            $s     = strtotime("-$i monday", time());
            $arr[] = [
                'date'  => date('W', $s),
                'start' => date('m.d', $s),
                'end'   => date('m.d', strtotime('+6 days', $s))
            ];
        }
        return array_reverse($arr);
    }
}

/**
 * 提示性错误
 */
if (!function_exists('notice')) {
    function notice($code, $args = [], $data = [], $statusCode = 200, $previous = null, $headers = [])
    {
        $message = empty(config('code.' . $code)) ? '未知的错误' : config('code.' . $code);

        if (!empty($args)) {
            $message = vsprintf($message, $args);
        }

        throw new Error($code, $message, $statusCode, $previous, $headers, $data);
    }
}

/**
 * 根据数组指定的键对应的值, 作为新数组的键名
 *
 * @param array  $arr 二维数组
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

/**
 * 获取orm模型完整sql
 */
if (!function_exists('getSql')) {
    // 获取执行sql
    function getSql($query)
    {
        $sql = str_replace("?", "'%s'", $query->toSql());

        return vsprintf($sql, $query->getBindings());
    }
}

/**
 * 数组分页
 *
 * @param int   $size 每页数量
 * @param int   $page 页码
 * @param array $list 列表数组
 *
 * @return array
 */
if (!function_exists('pageArray')) {
    function pageArray($page = 1, $size = 15, $list = [])
    {

        // 判断当前页面是否为空 如果为空就表示为第一页面
        $page = (empty($page)) ? '1' : $page;

        $totals = count($list);
        // 计算总页数
        $countpage = ceil($totals / $size);

        if ($page > $countpage) {
            $page = $countpage;
        }

        // 计算每次分页的开始位置
        $start = ($page - 1) * $size;

        $pagedata = array_slice($list, $start, $size);

        return $pagedata;
    }
}


/**
 * 获取时间段查询23:59分格式
 *
 * @param int $type 时间段类型 1=今日 2=最近7天 3=最近30天 5=昨日
 *
 * @return array
 */
if (!function_exists('getWhereTimeData')) {
    function getWhereTimeData($type = 0)
    {
        $startTime = $endTime = '';

        switch ($type) {
            case 1:
                $startTime = date('Y-m-d 00:00:00');
                break;
            case 2:
                $startTime = date('Y-m-d', strtotime('-6 days'));
                break;
            case 3:
                $startTime = date('Y-m-d', strtotime('-29 days'));
                break;
            case 5:
                $startTime = date('Y-m-d', strtotime('-1 days'));
                $endTime   = date('Y-m-d 23:59:59', strtotime('-1 days'));
                break;
        }

        // 如果时间类型存在
        if (!empty($startTime)) {
            $startTime = strtotime($startTime);
            if (empty($endTime)) {
                $endTime = strtotime(date('Y-m-d 23:59:59', time()));
            } else {
                $endTime = strtotime($endTime);
            }
        }

        return [$startTime, $endTime];
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

    if (!function_exists('switchAppHostByDomain')) {
        /**
         * 切换企业域名
         *
         * @param $domain
         *
         * @return array|bool
         * @throws Exception
         */
        function switchAppHostByDomain($domain)
        {
            if (empty($domain)) {
                return false;
            }

            // 获取企业版本信息
            $enterprise = (new EnterpriseService())->getDomainInfo($domain);

            if ($enterprise['version'] != Constant::NORMAL_VERSION) {
                // 获取网关域名
                $originConfig = config('client.api_url.shop');
                if (strpos($originConfig, 'gw-shop') == false) {
                    $gwShop    = config('client.api_url.shop') . '/gw-shop';
                    $gwConsole = config('client.api_url.console') . '/gw-console';
                    // 重写SHOP域名
                    config(['client.api_url.shop' => $gwShop]);
                    config(['client.api_url.console' => $gwConsole]);
                    config(['client.api_url.build' => 'http://49.235.249.11:5678']);
                    if (in_array(getenv('APP_ENV'), ['dev', 'test'])) {
                        \Illuminate\Support\Facades\Log::info($enterprise['id'] . '企业切换域名', config('client'));
                    }
                }
            }

            return $enterprise;
        }
    }

    if (!function_exists('switchAppHost')) {
        /**
         * 切换企业域名
         *
         * @param $param
         *
         * @return bool
         * @throws Exception
         */
        function switchAppHost($param)
        {
            if (empty($param['enterprise_id']) && empty($param['hash'])) {
                return false;
            }

            // 获取企业版本信息
            $enterprise = (new \CrsSdk\Boss\Enterprise())->detail($param);

            if ($enterprise['data']['version'] != Constant::NORMAL_VERSION) {

                // 获取网关域名
                $originConfig = config('client.api_url.shop');
                if (strpos($originConfig, 'gw-shop') == false) {
                    $gwShop    = config('client.api_url.shop') . '/gw-shop';
                    $gwConsole = config('client.api_url.console') . '/gw-console';
                    // 重写SHOP域名
                    config(['client.api_url.shop' => $gwShop]);
                    config(['client.api_url.console' => $gwConsole]);
                }
            }

            return true;
        }
    }
}

if (!function_exists('msectime')) {
    /**
     * 返回当前的毫秒时间戳
     */
    function msectime()
    {
        list($msec, $sec) = explode(' ', microtime());
        $msectime = (float)sprintf('%.0f', (floatval($msec) + floatval($sec)) * 1000);

        return $msectime;
    }
}

if (!function_exists('idGenerate')) {
    function idGenerate()
    {
        $url    = config('client.api_url.uid');
        $client = (new Client())->get($url, [
            'http_errors' => false
        ])->getBody()->getContents();

        return (int)$client;
    }
}

if (!function_exists('uids')) {
    function uids($num)
    {
        $url    = config('client.api_url.uids');
        $client = (new Client())->get($url .'/batch?num=' . $num, [
            'http_errors' => false
        ])->getBody()->getContents();
        $res = json_decode($client, true);
        if (isset($res['data']) && isset($res['data']['ids'])) {
            return  $res['data']['ids'];
        } else {
            throw new Error(100000, '批量获取分布式id失败:' . $client);
        }
    }
}
/**
 * [std_class_object_to_array 将对象转成数组]
 *
 * @param [stdclass] $stdclassobject [对象]
 *
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

if (!function_exists('successful')) {
    function successful(array $data = [], $httpCode = 200)
    {
        return response()->json([
            'code'       => 200,
            'message'    => 'success',
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => !empty($data) ? $data : (object)$data
        ], $httpCode);
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

/**
 * 基本写入日志
 */
if (!function_exists('writeLog')) {
    function writeLog($keyword, Exception $error)
    {
        Log::error("---" . $keyword . "---" .
            "---错误码---" . $error->getCode() .
            "---失败原因---" . $error->getMessage() .
            "---文件名称---" . $error->getFile() .
            "---文件行数---" . $error->getLine()
        );
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

if (!function_exists('handleDateTimeISO')) {
    function handleDateTimeISO($dateTimeISO)
    {
        return $dateTimeISO ? date('Y-m-d H:i:s', strtotime($dateTimeISO)) : '-';
    }
}

if (!function_exists('handleDateTime')) {
    function handleDateTime($dateTime): string
    {
        return $dateTime ?? '-';
    }
}

if (!function_exists('handleDateTimeUTC')) {
    function handleDateTimeUTC($dateTime): string
    {
        return $dateTime ? str_replace(array('T', 'Z'), ' ', $dateTime) : '-';
    }
}

if (!function_exists('handleDateTimeUTCToDay')) {
    function handleDateTimeUTCToDay($dateTime): string
    {
        return $dateTime ? date("Y-m-d", strtotime(str_replace(array('T', 'Z'), ' ', $dateTime))) : '-';
    }
}

if (!function_exists('uniqueMultiArrayByField')) {

    // OR 代表 任何一个重复就算作重复  AND 代表 两个都重复才算重复
    function uniqueMultiArrayByField($arr, $fields, $rel = 'OR')
    {
        $tmpArr    = [];
        $resultArr = [];
        $errArr    = [];
        if ($rel == 'AND' || count($fields) == 1) {
            foreach ($arr as $val) {
                $tmpStr = '';
                foreach ($fields as $field) {
                    $tmpStr .= $val[$field] . '|';
                }
                if (!in_array($tmpStr, $tmpArr)) {
                    $resultArr[] = $val;
                    $tmpArr[]    = $tmpStr;
                } else {
                    $errArr[] = $val;
                }
            }
        } else if ($rel == 'OR') {
            foreach ($arr as $val) {
                $exist = false;
                foreach ($fields as $field) {
                    $varField = 'var_' . $field;
                    if (!isset($$varField)) $$varField = [];
                    if (in_array($val[$field], $$varField)) {
                        $exist = true;
                        break;
                    }
                    array_push($$varField, $val[$field]);
                }

                if (!$exist) {
                    $resultArr[] = $val;
                } else {
                    $errArr[] = $val;
                }
            }
        }


        return [$resultArr, $errArr];
    }
}

if (!function_exists('arrayDiffTwoDimension')) {
    //二维数组取差集
    function arrayDiffTwoDimension($arr1, $arr2)
    {
        $arr = [];
        foreach ($arr1 as $key => $value) {
            if (!in_array($value, $arr2)) {
                $arr[] = $value;
            }
        }
        return $arr;
    }
}

if (!function_exists("array_xml")) {
    /**
     * 数组转xml
     *
     * @param array $data
     *
     * @return string
     */
    function array_xml(array $data): string
    {
        $xml = "<xml>";
        foreach ($data as $key => $val) {
            $xml .= "<$key><![CDATA[$val]]></$key>";
        }
        $xml .= "</xml>";
        return $xml;
    }
}


/**
 * 删除多维数组中指定key键以及value值
 *
 * @param array $arr 数组
 * @param string  $key  键（key）
 *
 * @return array
 */
if (!function_exists("arrayRemoveByKey")) {

    function arrayRemoveByKey($arr, $key): array
    {
        $new_array = [];
        foreach ($arr as $k => $v) {
            if (!is_array($v)) {
                if ($k == $key) {
                    $new_array = $v;
                } else {
                    $new_array[$k] = $v;
                }
            } else {
                // 判断是一维数组
                if (count($v) == count($v, 1)) {
                    /* 判断当前元素的值是不是数组，是的话就递归下去，重复查找这个函数。 */
                    if ($k == $key) {
                        $new_array = arrayRemoveByKey($v, $key);;
                    } else {
                        $new_array[$k] = arrayRemoveByKey($v, $key);
                    }
                } else {  // 判断是二维数组
                    if ($k == $key) {
                        $new_array[] = arrayRemoveByKey($v, $key);
                    } else {
                        $new_array[$k] = arrayRemoveByKey($v, $key);
                    }
                }
            }
        }
        return $new_array;
    }
}
