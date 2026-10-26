<?php

namespace App\Models;

class ActivityPrize extends Base
{

    public $primaryKey = 'id';
    protected $table = 'activity_prize';
    protected $fillable = [
        'id',
        'enterprise_id',
        'activity_id',
        'prize_name',
        'prize_short_name',
        'prize_num',
        'prize_level',
        'prize_ratio',
        'prize_status',
    ];
}