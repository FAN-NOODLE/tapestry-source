<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class MemberChangeRecord extends Base
{
    /**
     * 分配（待确认）
     * @var int
     */
    const UNALLOCATED = 2;

    /**
     * 表名
     * @var string
     */
    protected $table = 'member_change_record';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * 已确认交接
     */
    const CONFIRM_HANDOVER_TRUE = 1;

    /**
     * 待确认交接
     */
    const CONFIRM_HANDOVER_FALSE = 2;

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'old_mem_uid',
        'ec_userid',
        'mem_uid',
        'status'
    ];
}
