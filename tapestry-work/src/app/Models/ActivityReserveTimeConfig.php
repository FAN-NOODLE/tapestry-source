<?php

namespace App\Models;

class ActivityReserveTimeConfig extends Base
{

    public $primaryKey = 'id';


    protected $table    = 'activity_reserve_time_config';
    protected $fillable = [
        'id',
        'enterprise_id',
        'activity_id',
        'reserve_time',
        'reserve_num',
        'reserved_num',
    ];

}