<?php
/**
 * Created By PhpStorm
 * @Date: 2021/4/14
 * @Time: 10:26
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: HRSendMessageError.php
 */

namespace App\Library\Message\Exception;


use App\Exceptions\Error;
use App\Library\Message\Response;

class HRSendMessageError extends Error
{
    public function __construct(Response $response, $message='调用短信接口异常', \Exception $previous = null, array $headers = [])
    {
        parent::__construct(100000, $message, $response->getBody(), $response->getStateCode(), $previous, $headers);
    }
}