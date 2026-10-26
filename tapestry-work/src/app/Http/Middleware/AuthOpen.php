<?php

namespace App\Http\Middleware;

use Closure;
use App\Exceptions\Error;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;

class AuthOpen
{

    /**
     * 开放接口签名校验
     *
     * @param Request $request
     * @param Closure $next
     * @return JsonResponse|mixed
     */
    public function handle($request, Closure $next)
    {
        // 接收参数
        $input = $request->input();
        //把当前参数中的路径排除
        if(isset($input['s']) && strpos($input['s'], '/') !== false){
            unset($input['s']);
        }

        Log::info('work-open', [
            '请求地址' => $request->path(),
            '请求方式' => $request->method(),
            '请求参数' => $input,
            '请求IP' => $request->getClientIp(),
            'header参数' => $request->header()
        ]);


        $timestamp = $request->header('timestamp');
        $appKey = $request->header('appKey');
        $sign = $request->header('sign');
        if (!$timestamp || !$appKey || !$sign) {
            error(110007, [], 401);
        }
        // 时间戳有效期为10分钟
        if ($timestamp + 600 < time()) {
            error(100005, [], 401);
        }

        try {
            // 根据企业Hash读取企业信息
            if (!isset($input['enterprise_hash']) || empty($input['enterprise_hash'])) {
                error(110002, [], 401);
            }
            $enterpriseInfo = (new Enterprise())->detail(['hash' => $input['enterprise_hash']])['data'];
        } catch (\Exception $e) {
            throw new Error(100000, $e->getMessage(), 500);
        }

        if ($sign !== 'work.test' && env('APP_ENV') != 'prod') {
            if (!isset($appKey) || empty($appKey)) {
                error(110007, [], 401);
            }
            $secretConfig = config('open.sign_config');
            if (!isset($secretConfig[$appKey])) {
                throw new Error(100000, 'app_key_error', 401);
            }
            $secret = $secretConfig[$appKey];
            /*$input  = array_filter($input, function ($v) {
                return $v !== '';
            }, ARRAY_FILTER_USE_BOTH);*/

            // 排序
            $this->ksortArr($input);
            // 按照规则拼接键和值
            $newSign = md5(json_encode($input) . $timestamp . $secret);

//            dd(json_encode($input) . $timestamp . $secret, $newSign);
            // 验签失败
            if ($newSign !== $sign) {
                throw new Error(100000, 'sign_error', 401);
            }

        }
        //传递企业信息
        $request->attributes->add(['enterprise_info' => $enterpriseInfo]);

        return $next($request);
    }


    /**
     * 按照ASCII字典序排序
     *
     * @param $arr
     * @return mixed
     */
    private function ksortArr($arr)
    {
        ksort($arr);
        foreach ($arr as $k => $v) {
            if (is_array($v)) {
                $arr[$k] = $this->ksortArr($v);
            }
        }
        return $arr;
    }
}
