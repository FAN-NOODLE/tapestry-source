<?php
namespace App\Library\Message;

class Response
{
    protected $response = [];

    public function __construct($response)
    {
        $this->response = $response;
    }

    public function getBody()
    {
        if(!is_array($this->response)){
            return [$this->response];
        }

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
        if(isset($this->response['msgid'])){
            return $this->response['msgid'];
        }
        return '返回未知的错误';
    }

    public function isOk()
    {
        return (int)$this->getCode() === 0;
    }

    public function getData()
    {
        return $this->response['data'] ?? null;
    }

}