<?php


namespace App\Models;

class SubscribeMessageRecord extends Base
{
    const TO_STORE        = 1;
    const SERVICE_CHANGE  = 2;
    const SERVICE_COMMENT = 3;
    protected $table = 'subscribe_message_record';

    protected $guarded = [];
}
