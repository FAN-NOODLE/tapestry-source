<?php

namespace App\Services;

use App\Models\EnterpriseLogSystem;
use Illuminate\Support\Facades\{Log, Request};

class EnterpriseLogService
{
    /**
     * 调试级别日志
     *
     * @param $message
     * @param $context mixed
     */
    public static function debug($message, $context = [])
    {
        self::write(__FUNCTION__, $message, $context);
    }

    /**
     * 一般级别日志
     *
     * @param $message
     * @param $context mixed
     */
    public static function info($message, $context = [])
    {
        self::write(__FUNCTION__, $message, $context);
    }

    /**
     * 错误级别日志
     *
     * @param $message
     * @param $context mixed
     */
    public static function error($message, $context = [])
    {
        self::write(__FUNCTION__, $message, $context);
    }

    /**
     * 紧急级别日志
     *
     * @param $message
     * @param $context mixed
     */
    public static function critical($message, $context = [])
    {
        self::write(__FUNCTION__, $message, $context);
    }

    /**
     * 写入DB
     *
     * @param $level string 级别
     * @param $message string 信息
     * @param $context mixed 内容
     */
    private static function write($level, $message, $context = [])
    {
        if (!is_array($context)) {
            $context = [$context];
        }

        $data = [
            'request_path' => Request::instance()->path(),
            'message'      => $message,
            'context'      => json_encode($context, JSON_UNESCAPED_UNICODE),
            'level'        => $level
        ];

        try {
            EnterpriseLogSystem::query()->create($data);
        } catch (\Exception $exception) {
            Log::critical('写入日志表异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'data'    => $data
            ]);
        }
    }
}
