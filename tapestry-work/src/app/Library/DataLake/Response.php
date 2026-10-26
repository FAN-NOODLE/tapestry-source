<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/17 0017
 * Time: 15:28
 */

namespace App\Library\DataLake;


class Response
{
    private $response = [];

    public function __construct(array $response)
    {
        $this->response = $response;
    }

    public function getBody()
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
        return 'DL返回未知的错误';
    }

    public function isOk()
    {
        return true;
    }

    public function getData()
    {
        return $this->response['data'] ?? [];
    }

}