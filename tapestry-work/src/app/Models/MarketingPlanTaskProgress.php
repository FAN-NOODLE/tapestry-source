<?php
/**
 * 营销计划任务进度记录表
 */
namespace App\Models;

class MarketingPlanTaskProgress extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'marketing_plan_task_progress';

    // 主键
    protected $primaryKey = 'id';

    //完成方式(1、复制话术，2、转发给顾客，3、分享朋友圈  4、保存素材)
    const FINISHED_TYPE_COPY = 1;

    //完成方式(1、复制话术，2、转发给顾客，3、分享朋友圈  4、保存素材)
    const FINISHED_TYPE_FORWARD = 2;

    //完成方式(1、复制话术，2、转发给顾客，3、分享朋友圈  4、保存素材)
    const FINISHED_TYPE_SHARE = 3;

    //完成方式(1、复制话术，2、转发给顾客，3、分享朋友圈  4、保存素材)
    const FINISHED_TYPE_DOWNLOAD = 4;

    /**
     * 是否完成
     * is_finished  1 未完成 2 已完成
     */
    const TASK_NOT_FINISHED = 1;
    const TASK_ALREADY_FINISHED = 2;

    protected $fillable = [
        'enterprise_id',
        'task_id',
        'mem_uid',
        'worker_id',
        'finished_type',
        'is_finished'
    ];
}
