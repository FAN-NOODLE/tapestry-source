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
        $request = Request::instance();
        $request_route = $request->method() . " " . $request->path();
        if (defined('LOG_ID')) {
            $request_id = LOG_ID;
        } else {
            $request_id = md5(uniqid()) . rand(10000, 99999);
        }
        // 这个就是最终要记录的数组，最后转成Json并记录进日志
        $newRecord = [
            'request_id'    => $request_id,
            'time'          => $record['datetime']->format('Y-m-d H:i:s'),
            'request_route' => $request_route,
            'message'       => $record['message'],
            'context'       => $record['context'],
        ];

        $json = $this->toJson($this->normalize($newRecord), true) . ($this->appendNewline ? "\n" : '');

        return $json;
    }
}