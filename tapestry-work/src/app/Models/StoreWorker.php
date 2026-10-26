<?php

namespace App\Models;

class StoreWorker extends Base
{
    /**
     * 是否是店长：1=是，2=否
     */
    const BOSS_YES = 1;
    const BOSS_NO = 2;

    /**
     * 员工状态：1=在职，2=离职
     */
    const IN_SERVICE = 1;
    const OUT_SERVICE = 2;

    /**
     * 员工来源 企业微信
     */
    const SOURCE_WORK = 1;
    /**
     * 表名称
     * @var string
     */
    protected $table = 'store_worker';

    /**
     * The attributes that are mass assignable.
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'user_id',
        'store_id',
        'name',
        'nickname',
        'mobile',
        'sms_mobile',
        'mem_uid',
        'mem_userid',
        'avatar',
        'is_boss',
        'qr_code',
        'worker_no',
        'worker_status',
        'source',
        'store_sun_code',
        'type',
        'contact_qr_code_id',
        'created_at'
    ];

}
