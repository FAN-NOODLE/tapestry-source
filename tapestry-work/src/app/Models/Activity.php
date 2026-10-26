<?php

namespace App\Models;

class Activity extends Base
{

    protected $table = 'activity';

    public $primaryKey = 'id';

    // 活动状态，0未开始，1进行中，2已结束，3暂停
    const StatusNotStarted = 0;
    const StatusRun        = 1;
    const StatusEnd        = 2;
    const StatusPause      = 3;

    protected $fillable = [
        'id',
        'enterprise_id',
        'activity_name',
        'activity_image',
        'activity_content',
        'start_time',
        'end_time',
        'check_time',
        'activity_status',
        'activity_num',
    ];
}