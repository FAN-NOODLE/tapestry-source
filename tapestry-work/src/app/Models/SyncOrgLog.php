<?php

namespace App\Models;

class SyncOrgLog extends Base
{

    //UC同步中
    const STATUS_UC_PROCESSING = 0;
    //UC同步完成
    const STATUS_UC_OK = 1;


    const SOURCE_TYPE_MANUAL = 2;
    const SOURCE_TYPE_AUTO = 1;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'sync_org_log';


    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = ['id'];
}
