<?php

namespace App\Models;

class WithdrawalAutomaticReview extends Base
{

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'withdrawal_automatic_review';

    // 主键
    protected $primaryKey = 'id';

    /**
     * 员工身份  1、导购  2、合伙人
     */
    const MEM_TYPE_WORKER = 1;

    const MEM_TYPE_PARTNER = 2;


    protected $fillable = [
        'id',
        'auto_status',
        'enterprise_id',
        'auto_at',
        'last_at',
        'created_at',
        'updated_at',
        'deleted_at'
    ];
}