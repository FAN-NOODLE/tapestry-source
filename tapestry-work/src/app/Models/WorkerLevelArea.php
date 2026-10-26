<?php

namespace App\Models;

class WorkerLevelArea extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'worker_level_area';

    /**
     * The attributes that are mass assignable.
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'worker_id',
        'area_json'
    ];

}
