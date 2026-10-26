<?php

namespace App\Models;

class UserLoginLog extends BaseModel
{
    protected $table = 'user_login_log';

    protected $fillable = [
        'id',
        'enterprise_id',
        'mem_uid',
        'mem_userid',
        'login_day_time',
        'created_at',
    ];

}
