<?php

namespace App\Models;

class MemberUserRelationVisit extends Base
{
    protected $table = 'member_user_relation_visit';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'mem_uid',
        'count',
        'result_time',
        'created_at',
        'updated_at',
    ];
}
