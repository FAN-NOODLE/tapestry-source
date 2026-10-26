<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class TaskMemberRecord extends Base
{
    /*
     * 状态类型：1:未完成，2:已完成
     */
    const CUSTOMER_IN_THE_ROUGH_STATUS = 1;
    const CUSTOMER_COMPLETED_STATUS = 2;

    /**
     * 任务显示状态 1=显示 2=隐藏
     */
    const TASK_SHOW_STATUS = 1;
    const TASK_HIDE_STATUS = 2;

    /**
     * 表名
     * @var string
     */
    protected $table = 'task_member_record';

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
        'task_id',
        'mem_uid',
        'customer_id',
        'coupon_ids',
        'task_type',
        'display_state',
        'status',
        'updated_at'
    ];
}
