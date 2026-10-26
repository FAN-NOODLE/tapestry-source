<?php

namespace App\Models;

class WorkSunCode extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'work_sun_code';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'worker_id',
        'sun_code',
    ];
}
