<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MissionCenterAppletsConfig extends Base
{
    /**
     * 表名称
     * @var string
     */
    protected $table = 'mission_center_applets_config';

    // 主键
    protected $primaryKey = 'id';

    /**
     * 1、培训任务 2、邀约任务
     */
    const TYPE_TRAINING_TASK = 1;
    const TYPE_INVITATION_TASK = 2;

    protected $fillable = [
        'id',
        'enterprise_id',
        'name',
        'app_id',
        'app_original_id',
        'app_image_url',
        'type',
        'created_at',
        'updated_at',
        'deleted_at',
    ];

}