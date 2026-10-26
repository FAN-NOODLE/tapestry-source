<?php

namespace App\Library;

class Json
{
    public static function encode($value, $option = JSON_UNESCAPED_UNICODE + JSON_UNESCAPED_SLASHES, $depth = 512)
    {
        return json_encode($value, $option, $depth);
    }

    public static function decode($data)
    {
        return json_decode($data, true);
    }
}
