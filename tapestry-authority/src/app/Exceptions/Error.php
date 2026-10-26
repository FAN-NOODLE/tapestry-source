<?php

namespace App\Exceptions;

use Symfony\Component\HttpKernel\Exception\HttpException as BaseHttpException;

class Error extends BaseHttpException
{
    /**
     * @var array
     */
    public $data = [];

    /**
     * Error constructor.
     * @param int $code
     * @param string $message
     * @param int $statusCode
     * @param \Exception|null $previous
     * @param array $headers
     * @param array $data
     */
    public function __construct(
        int $code,
        string $message,
        array $data = [],
        int $statusCode = 500,
        \Exception $previous = null,
        array $headers = []
    )
    {
        $this->data = $data;
        parent::__construct($statusCode, $message, $previous, $headers, $code);
    }
}
