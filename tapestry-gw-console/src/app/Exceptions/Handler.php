<?php

namespace App\Exceptions;

use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
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
     * @param Exception $e
     *
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
     * @param Request $request
     * @param Exception $e
     *
     * @return Response | mixed
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
     *
     * @return JsonResponse
     */
    private function error($e)
    {

        if ($e instanceof ValidationException) {
            $statusCode = 422;
            $message    = current(current(array_values($e->errors())));
            $code       = $statusCode;
        } else {
            $statusCode = $e->getStatusCode();
            $message    = $e->getMessage();
            $code       = $e->getCode() ?? 450000;
        }
        $data = empty($e->data) ? (object)[] : $e->data;

        if ($statusCode == 404) {
            $message = '路由地址错误';
        }

        if ($statusCode == 405) {
            $message = '请求方式错误';
        }

        $response = [
            'code'       => $code,
            'message'    => $message,
            'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
            'data'       => $data,
        ];

        if (in_array(config('app.env'), ['dev', 'test'])) {
            $request = app(\Illuminate\Http\Request::class);
            $insert  = [
                'request_id' => $_SERVER['x_request_id'] ?? LOG_ID,
                'ip'         => $request->getClientIp(),
                'port'       => $request->getPort(),
                'method'     => $request->method(),
                'url'        => $request->url(),
                'params'     => json_encode($request->all(), JSON_UNESCAPED_UNICODE),
                'code'       => $code,
                'message'    => $message,
                'data'       => json_encode($data, JSON_UNESCAPED_UNICODE),
            ];
            Log::error('错误信息：', $insert);
        }

        return response()->json($response, $statusCode, [], JSON_UNESCAPED_UNICODE);
    }
}
