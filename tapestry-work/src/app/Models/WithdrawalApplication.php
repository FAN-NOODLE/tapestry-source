<?php

namespace App\Models;

class WithdrawalApplication extends Base
{
    // 待审批
    const WAIT_APPLICATION_STATUS = 1;

    // 提现成功
    const WITHDRAW_SUCCESS_STATUS = 2;

    // 提现失败
    const WITHDRAW_FAIL_STATUS = 3;

    // 驳回
    const WITHDRAW_REJECT_STATUS = 4;

    //转账中
    const WITHDRAW_IN_STATUS = 5;
    /**
     * 1:导购申请 2:合伙人申请
     */
    const APP_TYPE_STAFF = 1;
    /**
     * 1:导购申请 2:合伙人申请
     */
    const APP_TYPE_PARTNER = 2;

    //状态值
    const  STATUS = [
        1 => '待审批',
        2 => '提现成功',
        3 => '提现失败',
        4 => '已驳回',
        5 => '转账中'
    ];

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'withdrawal_application';

    // 主键
    protected $primaryKey = 'id';


    /**
     * 表字段
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'operator_id',
        'worker_id',
        'application_amount',
        'request_initiated_at',
        'request_handle_at',
        'status',
        'reason',
        'created_at',
        'updated_at',
        'deleted_at',
    ];


}