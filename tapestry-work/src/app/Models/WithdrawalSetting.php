<?php

namespace App\Models;

class WithdrawalSetting extends Base
{

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'withdrawal_setting';

    // 主键
    protected $primaryKey = 'id';


    protected $fillable = [
        'id',
        'transfer_no',
        'enterprise_id',
        'type',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

    //合伙人
    const PARTNER_STATUS = 5;

}