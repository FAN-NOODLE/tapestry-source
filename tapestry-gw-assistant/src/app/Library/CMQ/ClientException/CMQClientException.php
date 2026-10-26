<?php
namespace App\Library\CMQ\ClientException;

/**
 * Created by PhpStorm.
 * Function: xxx
 * User: 龚国新
 * Date: 2019/10/22
 * Time: 21:01
 */

class CMQClientException extends CMQExceptionBase
{
    public function __construct($message, $code=-1, $data=array())
    {
        parent::__construct($message, $code, $data);
    }

    public function __toString()
    {
        return "CMQClientException  " .  $this->get_info();
    }
}