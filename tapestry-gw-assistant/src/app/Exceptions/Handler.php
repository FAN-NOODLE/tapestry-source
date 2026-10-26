<?php

namespace App\Exceptions;

use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Response;
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
     * @param Request   $request
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

        return parent::render($request, $e);
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
        $statusCode = ($e instanceof ValidationException) ? 422 : $e->getStatusCode();
        $message    = ($e instanceof ValidationException) ?
            current(current(array_values($e->errors()))) : $e->getMessage();

        $code = empty($e->getCode()) ? 650000 : $e->getCode();
        $data = empty($e->data) ? (object)[] : $e->data;

        if ($statusCode == 404) {
            $message = '路由地址错误';
        }

        if ($statusCode == 405) {
            $message = '请求方式错误';
        }

        $response = [
            'id'      => $_SERVER['x_request_id'] ?? LOG_ID,
            'code'    => $code,
            'status'  => $statusCode,
            'message' => $message,
            'error'   => 'ERROR',
            'data'    => $data
        ];

        Log::error('Response Error: ' . json_encode($response, JSON_UNESCAPED_UNICODE));
        return response()->json($response, $statusCode, [], JSON_UNESCAPED_UNICODE);
    }
}
