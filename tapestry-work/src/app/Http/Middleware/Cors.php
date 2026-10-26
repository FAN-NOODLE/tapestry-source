<?php
/*
 * Author:xx_lufei
 * Time:2016年11月3日13:15:51
 * Note:Access Control Headers.
 */

namespace App\Http\Middleware;


use Illuminate\Http\Request;
use Illuminate\Http\Response;

class Cors
{
    private $headers;

    public function handle(Request $request, \Closure $next)
    {
        $this->headers = [
            'Access-Control-Allow-Methods' => 'GET, POST, PUT, DELETE, OPTIONS',
            'Access-Control-Allow-Headers' => $request->header('Access-Control-Request-Headers'),
            'Access-Control-Allow-Credentials' => 'true',//允许客户端发送cookie
        ];

        //如果是复杂请求，先返回一个200，并allow该origin
        if ($request->isMethod('options')) {
            return $this->setCorsHeaders(new Response('OK', 200));
        }
        //如果是简单请求或者非跨域请求，则照常设置header
        $response = $next($request);
        $methodVariable = array($response, 'header');
        //这个判断是因为在开启session全局中间件之后，频繁的报header方法不存在，所以加上这个判断，存在header方法时才进行header的设置
        if (is_callable($methodVariable, false, $callable_name)) {
            return $this->setCorsHeaders($response);
        }
        return $response;
    }

    /**
     * @param $response
     * @return mixed
     */
    public function setCorsHeaders($response)
    {
        foreach ($this->headers as $key => $value) {
            $response->header($key, $value);
        }

        $response->header('Access-Control-Allow-Origin', $_SERVER['HTTP_ORIGIN'] ?? '*');

        return $response;
    }
}