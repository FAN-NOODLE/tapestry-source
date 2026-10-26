<?php
/**
 * 导购小程序太阳码表
 */
namespace App\Models;

class MemberBannerRelation extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'member_banner_relation';

    // 主键
    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'task_id',
        'mem_uid',
        'worker_id',
        'background_url',
        'suncode_mediaid',
        'suncode_url',
        'suncode_param',
        'created_at',
        'updated_at',
        'deleted_at',
    ];
}
