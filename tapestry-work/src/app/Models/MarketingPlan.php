<?php
/**
 * 营销计划表
 */
namespace App\Models;

class MarketingPlan extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'marketing_plan';

    // 主键
    protected $primaryKey = 'id';


    protected $fillable = [
        'id',
        'enterprise_id',
        'plan_title',
        'range_updated',
        'operator_id',
        'operator_name',
        'created_at',
        'updated_at',
        'deleted_at',
        'add_object_ids',
        'object_type'
    ];

    /**
     * 可见范围类型
     * object_type  1 区域 2  部门
     */
    const OBJECT_TYPE_STORE = 1;
    const OBJECT_TYPE_DEPARTMENT = 2;

    /**
     * 一对多 可见范围
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 11:18
     */
    public function range()
    {
        return $this->hasMany(MarketingPlanVisibleRange::class,'plan_id','id');
    }

    /**
     * 一对多 任务
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @author  zhaopengfei<zhaopengfei@vchangyi.com>
     * @date    2021/1/12 11:38
     */
    public function taskArray()
    {
        return $this->hasMany(MarketingPlanTask::class,'plan_id','id');
    }

    const plan_select = [
        //计划表
        "id as plan_id","plan_title","operator_name",
        "created_at","updated_at","object_type","plan_sn"
    ];


}
