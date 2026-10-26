<?php
/**
 * Created by PhpStorm.
 * User: guoxingong
 * Date: 2019/4/18
 * Time: 10:30
 */

namespace App\Library;

use Illuminate\Support\Facades\Request;
use Monolog\Formatter\JsonFormatter as BaseJsonFormatter;


class JsonFormatter extends BaseJsonFormatter
{
    public function format(array $record)
    {
        $request=Request::instance();
        $request_route=$request->method()." ".$request->path();
        $request_params=$request->all();
        // 这个就是最终要记录的数组，最后转成Json并记录进日志
        $newRecord = [
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'time' => $record['datetime']->format('Y-m-d H:i:s'),
            'total_time' => sprintf('%0.3f', microtime(true) - START_TIME),
            'message'    => $record['message'],
            'context'    => $record['context'],
            'app_name'   => env('APP_NAME'),
            'request_route' => $request_route,
            'request_params'=>$request_params,
        ];

        $json = $this->toJson($this->normalize($newRecord), true) . ($this->appendNewline ? "\n" : '');

        return $json;
    }
}