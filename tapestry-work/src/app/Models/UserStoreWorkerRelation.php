<?php

namespace App\Models;

class UserStoreWorkerRelation extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'user_store_worker_relation';

    /**
     * 白名单
     *
     * @var array
     */
    protected $guarded = [];

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'store_id',
        'user_id',
        'mem_uid',
        'worker_id',
    ];
}
