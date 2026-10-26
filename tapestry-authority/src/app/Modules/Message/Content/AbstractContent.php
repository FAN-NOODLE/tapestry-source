<?php

namespace App\Modules\Message\Content;

abstract class AbstractContent
{
    /**
     * @var array
     */
    protected $message;

    public function __construct($message)
    {
        $this->message = $message;
    }
}
