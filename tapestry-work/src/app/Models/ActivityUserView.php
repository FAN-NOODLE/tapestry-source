<?php

namespace App\Models;

class ActivityUserView extends Base
{

    public    $primaryKey = 'id';
    protected $table      = 'activity_user_view';
    protected $fillable   = [
        'id',
        'enterprise_id',
        'activity_id',
        'openid',
        'unionid',
        'user_name',
        'source',
        'share_num'
    ];
}