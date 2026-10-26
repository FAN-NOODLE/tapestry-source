<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MissionCenterVisibleRange extends Base
{
    /**
     * 表名称
     * @var string
     */
    protected $table = 'mission_center_visible_range';

    // 主键
    protected $primaryKey = 'id';

    /**
     * 1、按导购所属部门
     */
    const CHOICE_DEPARTMENT = 1;

    protected $fillable = [
        'id',
        'enterprise_id',
        'task_id',
        'object_type',
        'object_id',
        'created_at',
        'updated_at',
        'deleted_at',
        'mem_uid',
        'worker_id',
    ];

}