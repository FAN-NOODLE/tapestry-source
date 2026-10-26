<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/19
 * @Time: 9:31
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: MissionCenterTask.php
 */

namespace App\Models;


class MissionCenterTaskResource extends Base
{
    //1、图片(默认)，2、小程序海报 3、视频
    const RESOURCE_TYPE_PICTURE = 1;

    //1、图片(默认)，2、小程序海报 3、视频
    const RESOURCE_TYPE_APPLET_BANNER = 2;

    //1、图片(默认)，2、小程序海报 3、视频
    const RESOURCE_TYPE_VIDEO = 3;

    //素材是否更新:1、是  2、否
    const RESOURCE_UPDATED_YES = 1;

    //素材是否更新:1、是  2、否
    const RESOURCE_UPDATED_NO = 2;

    //小程序路径是否更新:1.是 2 否
    const APPLET_PATH_YES = 1;
    const APPLET_PATH_NO = 2;

    // 定时任务是否执行
    const  IS_CRONTAB_YES = 1;
    const  IS_CRONTAB_NO = 2;

    /**
     * 表名称
     * @var string
     */
    protected $table = 'mission_center_task_resource';

    // 主键
    protected $primaryKey = 'id';


    protected $fillable = [
        'id',
        'enterprise_id',
        'task_id',
        'type',
        'mediaid',
        'mediaid_updated_at',
        'resource_url',
        'applet_path',
        'resource_updated',
        'applet_card_title',
        'applet_card_cover_url',
        'is_crontab',
        'created_at',
        'updated_at',
        'deleted_at',
        'applet_path_update',
        'poster_url',
        'sun_code',
    ];

}