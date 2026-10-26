<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class TaskMemberRelation extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'task_member_relation';

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
        'task_id',
        'mem_uid',
        'enterprise_id',
        'store_id'
    ];
}
