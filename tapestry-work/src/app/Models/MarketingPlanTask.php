<?php
/**
 * 营销计划任务表
 */

namespace App\Models;

class MarketingPlanTask extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'marketing_plan_task';
    protected $guarded = [];
    // 主键
    protected $primaryKey = 'id';

    //查询最早为6月前
    const QUERY_BETWEEN_MONTH_AGO = 6;
    //查询最早为一月后
    const QUERY_BETWEEN_MONTH_LATER = 1;
    const task_select = [
        //任务表
        "id as task_id", "task_title", "start_time", "end_time", "content", "before_start_time", "type",
        "before_end_time", "fixed_time", "operator_name", "object_type", "task_sn", "add_object_ids"
    ];

    /**
     * 任务类型
     * type  1 培训任务 2 邀约任务 3 群发任务
     */
    const TYPE_TRAINING_TASK = 1;
    const TYPE_INVITATION_TASK = 2;
    const TYPE_GROUP_TASK = 3;

    /**
     * 任务状态
     * type  1 未开始 2 进行中 3 已结束
     */
    const TASK_NOT_STARTED = 1;
    const TASK_IN_PROGRESS = 2;
    const TASK_COMPLETE = 3;

    public function resource()
    {
        return $this->hasMany(MarketingPlanTaskResource::class, 'task_id', 'id');
    }

    /**
     * 一对多 可见范围
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @Date: 2021/5/19
     * @Time: 18:10
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function range()
    {
        return $this->hasMany(MarketingPlanVisibleRange::class, 'task_id', 'id');
    }

    public function tags()
    {
        return $this->hasMany(MarketingTaskTagRelation::class,'task_id','task_id');
    }
}
