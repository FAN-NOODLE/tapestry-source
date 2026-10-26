<?php

namespace App\Models;

class ActivityShareLog extends Base
{

    public    $primaryKey = 'id';
    protected $table      = 'activity_share_log';
    protected $fillable   = [
        'id',
        'enterprise_id',
        'activity_id',
        'openid',
        'source'
    ];
}