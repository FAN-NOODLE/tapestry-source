<?php

namespace App\Traits;

use Illuminate\Support\Facades\Log as RequestRecordLog;

trait RequestRecord
{
    /**
     * @param array $message
     * @param array $content
     * @param null|\Exception $exception
     */
    public function wxRecord($message, $content, $exception = null)
    {
        if (isset($message['message']) || (isset($message['errmsg']) && $message['errmsg'] !== 'ok')) {
            $message = $message['errmsg'] ?? $message['message'];
            $this->errWxRecord($message, $content, $exception);
        } else {
            $this->sucWxRecord($content);
        }
    }

    private function sucWxRecord($data = [])
    {
        $data = array_merge($data, [
            'time' => date('Y-m-d H:i:s')
        ]);

        RequestRecordLog::info('请求微信成功', $data);
    }

    /**
     * @param string $message
     * @param array $data
     * @param null|\Exception $exception
     */
    private function errWxRecord($message = '', $data = [], $exception = null)
    {
        $data = array_merge($data, [
            'time' => date('Y-m-d H:i:s')
        ]);

        if ($exception) {
            $data['data']['result'] = [
                'msg'  => $exception->getFile(),
                'file' => $exception->getFile(),
                'line' => $exception->getLine(),
                'code' => $exception->getCode(),
            ];
        }

        RequestRecordLog::error('请求微信异常:' . $message, $data);
    }
}