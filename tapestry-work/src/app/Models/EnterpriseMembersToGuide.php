<?php

namespace App\Models;

class EnterpriseMembersToGuide extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'enterprise_members_to_guide';

    /**
     * The attributes that are mass assignable.
     * @var array
     */
    protected $fillable = [
        'enterprise_id',
        'domain',
        'userid',
        'mem_uid',
        'is_exec',
        'result',
        'start_at',
        'end_at',
        'created_at'
    ];

}
