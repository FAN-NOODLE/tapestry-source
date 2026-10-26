<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MarketingTaskTagRelation extends Base
{
    /**
     * 表名称
     * @var string
     */
    protected $table = 'marketing_task_tag_relation';

    // 主键
    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'task_id',
        'tag_id',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

}