<?php

namespace App\Sdk\Active\CRM\CoachCore;


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
        return $this->response['code'] ?? 500;
    }

    public function getMsg()
    {
        if(isset($this->response['msg'])){
            return $this->response['msg'];
        }
        if(isset($this->response['message'])){
            return $this->response['message'];
        }
        if(isset($this->response['responseDesc'])){
            return $this->response['responseDesc'];
        }
        return 'CDP返回未知的错误';
    }

    public function isOk(): bool
    {
        return (int)$this->getCode() === 200;
    }

    public function getData(): ?array
    {
        return $this->response['data'] ?? [];
    }

    public function getResCode()
    {
        return $this->response['responseCode'] ?? 500;
    }

    public function isOkCode(): bool
    {
        return (int)$this->getResCode() === 200;
    }

    public function getResData(): ?array
    {
        return $this->response ?? [];
    }

    public function getResSuccessCode(): int
    {
        return $this->response['resultCode'] == "SUCCESS" ? 200 : 500;
    }

    public function isSuccessCode(): bool
    {
        return $this->getResSuccessCode() === 200;
    }

}