<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 17:24
 */

namespace App\Library\DataLake\Exception;


use App\Exceptions\Error;
use App\Library\DataLake\Response;

class QueryDataLakeBusinessError extends Error
{
    public function __construct(
        Response $response,
        $message='查询DataLake业务异常',
        \Exception $previous = null,
        array $headers = []
    )
    {
        parent::__construct(100000, $message, 500,$previous, $headers,$response->getBody());
    }
}