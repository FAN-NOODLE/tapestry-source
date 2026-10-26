<?php

namespace App\Sdk\Active\CRM\CoachCore;


class ResponseCore
{
    private $response = [];

    public function __construct(array $response)
    {
        $this->response = $response;
    }

    public function getBody(): array
    {
        return $this->response;
    }

    public function getCode()
    {
        if (isset($this->response['ErrorCode'])) {
            $this->response['code'] = $this->response['ErrorCode'] == 0 ? 200 : $this->response['ErrorCode'];
        }
        return $this->response['code'] ?? 500;
    }

    public function getMsg()
    {
        if (isset($this->response['Text'])) {
            return $this->response['message'] = $this->response['Text'];
        }
        if (isset($this->response['Message'])) {
            return $this->response['message'] = $this->response['Message'];
        }
        return '返回未知的错误';
    }

    public function getText()
    {
        return $this->response['Text'];
    }

    public function isOk(): bool
    {
        return (int)$this->getCode() === 200;
    }

    public function getData(): ?array
    {
        if (isset($this->response['ResultData'])) {
            return $this->response['data'] = $this->response['ResultData'];
        }
        return $this->response['data'] ?? [];
    }

}