<?php


namespace App\Micros\Message;


abstract class AbstractMessage
{
    abstract function send($params);
}
