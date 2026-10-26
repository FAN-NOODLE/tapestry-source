<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MissionCenterTask extends Base
{
    /**
     * 表名称
     * @var string
     */
    protected $table = 'mission_center_task';

    // 主键
    protected $primaryKey = 'id';


    protected $fillable = [
        'id',
        'enterprise_id',
        'task_title',
        'type',
        'start_time',
        'end_time',
        'before_start_time',
        'before_end_time',
        'fixed_time',
        'content',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    /**
     * 任务类型
     * type  1 培训任务 2 邀约任务 3 群发任务
     */
    const TYPE_TRAINING_TASK = 1;
    const TYPE_INVITATION_TASK = 2;
    const TYPE_GROUP_TASK = 3;

    /**
     * 一对多 可见范围
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @Date: 2021/5/19
     * @Time: 9:38
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function range()
    {
        return $this->hasMany(MissionCenterVisibleRange::class,'task_id','id');
    }


    /**
     * 一对多 任务进度
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @Date: 2021/5/19
     * @Time: 9:39
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function progress()
    {
        return $this->hasMany(MissionCenterTaskProgress::class,'task_id','id');
    }

    /**
     * 一对多 任务素材
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     * @Date: 2021/5/19
     * @Time: 9:39
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function resource()
    {
        return $this->hasMany(MissionCenterTaskResource::class, 'task_id', 'id');
    }




}