<?php
/**
 * Created by PhpStorm.
 * User: guoxingong
 * Date: 2019/4/18
 * Time: 10:30
 */

namespace App\Library;

use Monolog\Formatter\JsonFormatter as BaseJsonFormatter;

class JsonFormatter extends BaseJsonFormatter
{
    public function format(array $record)
    {
        // 这个就是最终要记录的数组，最后转成Json并记录进日志
        $newRecord = [
            'level'      => $record['level_name'],
            'time'       => $record['datetime']->format('Y-m-d H:i:s'),
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'message'    => $record['message'],
            'context'    => $record['context'],
        ];

        return $this->toJson($this->normalize($newRecord), true).($this->appendNewline ? "\n" : '');
    }
}
