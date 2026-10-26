<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MissionCenterTaskProgress extends Base
{
    /**
     * 表名称
     * @var string
     */
    protected $table = 'mission_center_task_progress';

    // 主键
    protected $primaryKey = 'id';


    protected $fillable = [
        'id',
        'enterprise_id',
        'task_id',
        'mem_uid',
        'worker_id',
        'is_finished',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    /**
     * 是否完成
     * is_finished  1 未完成 2 已完成
     */
    const TASK_NOT_FINISHED = 1;
    const TASK_ALREADY_FINISHED = 2;


}