<?php

namespace App\Exceptions;

use Exception;
use Illuminate\Support\Facades\App;
use Illuminate\Support\Facades\Log;
use Illuminate\Validation\ValidationException;
use Illuminate\Auth\Access\AuthorizationException;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Laravel\Lumen\Exceptions\Handler as ExceptionHandler;
use Symfony\Component\HttpKernel\Exception\HttpException;

class Handler extends ExceptionHandler
{
    /**
     * A list of the exception types that should not be reported.
     *
     * @var array
     */
    protected $dontReport = [
        AuthorizationException::class,
        HttpException::class,
        ModelNotFoundException::class,
        ValidationException::class,
    ];

    /**
     * Report or log an exception.
     * This is a great spot to send exceptions to Sentry, Bugsnag, etc.
     *
     * @param \Exception $e
     * @return void
     * @throws Exception
     */
    public function report(Exception $e)
    {
        parent::report($e);
    }

    /**
     * Render an exception into an HTTP response.
     *
     * @param \Illuminate\Http\Request $request
     * @param \Exception $e
     * @return \Illuminate\Http\Response | mixed
     */
    public function render($request, Exception $e)
    {
        if ($e instanceof Error) {
            return $this->error($e);
        }

        if ($e instanceof HttpException) {
            return $this->error($e);
        }

        if ($e instanceof ValidationException) {
            return $this->error($e);
        }

        // 系统级异常输出json
        $response = [
            'code'       => $e->getCode() === 0 ? 100000 : $e->getCode(),
            'message'    => $e->getMessage(),
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => null,
        ];

        // 非生产环境返回异常信息
        if (App::environment() !== 'prod') {
            $response['error'] = [
                'message' => $e->getMessage(),
                'file'    => $e->getFile(),
                'line'    => $e->getLine(),
                'code'    => $e->getCode()
            ];
        }

        return response()->json($response, 500);
    }

    /**
     * 自定义错误输出
     *
     * @param Error $e
     * @return \Illuminate\Http\JsonResponse
     */
    private function error($e)
    {
        $statusCode = ($e instanceof ValidationException) ? 500 : $e->getStatusCode();
        $message = ($e instanceof ValidationException) ? current(current(array_values($e->errors()))) : $e->getMessage();

        $code = empty($e->getCode()) ? 100000 : $e->getCode();
        $data = empty($e->data) ? [] : $e->data;

        $response = [
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'code'       => $code,
            'message'    => $message,
            'data'       => $data
        ];

        Log::error('Response Error', $response);

        return response()->json($response, $statusCode);
    }
}
