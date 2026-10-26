<?php

namespace App\Sdk\Active\CRM\KsCore;


class Response
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
        if(isset($this->response['errcode'])){
            $this->response['code'] = $this->response['errcode'];
        }
        return $this->response['code'] ?? 500;
    }

    public function getMsg()
    {
        if (isset($this->response['errmsg'])) {
            return $this->response['message'] = $this->response['errmsg'];
        }
        if (isset($this->response['message'])) {
            return $this->response['message'];
        }
        return 'CRM返回未知的错误';
    }

    public function isOk(): bool
    {
        return (int)$this->getCode() === 0;
    }

    public function getData(): ?array
    {
        return $this->response['data'] ?? [];
    }

    public function getDataReq(): array
    {
        return $this->response;
    }

}