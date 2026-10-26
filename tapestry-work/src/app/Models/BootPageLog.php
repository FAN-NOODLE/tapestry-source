<?php
/**
 * 引导页浏览记录表
 */
namespace App\Models;

class BootPageLog extends Base
{

    /**
     * 表名称
     * @var string
     */
    protected $table = 'boot_page_logs';

    // 主键
    protected $primaryKey = 'id';

    protected $fillable = [
        'enterprise_id',
        'version_code',
        'worker_id'
    ];
}
