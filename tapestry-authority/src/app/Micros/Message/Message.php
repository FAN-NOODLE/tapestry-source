<?php


namespace App\Micros\Message;


use ReflectionClass;
use ReflectionException;

class Message
{
    private $item;

    function getClassName($itemName)
    {
        try {
            $nameSpace = 'App\Micros\Message';
            // ReflectionClass 类, 报告了一个类的有关信息。
            $class      = new ReflectionClass($nameSpace.'\\'.$itemName);
            $this->item = $class->newInstance();
        } catch (ReflectionException $e) {
            $this->item = "";
        }

        return $this->item;
    }

    function send($params)
    {
        return $this->item->send($params);
    }
}
