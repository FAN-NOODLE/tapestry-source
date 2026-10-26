<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class ShoppingGuideCashWithdrawal extends Base
{
    protected $table = 'shopping_guide_cash_withdrawal';

    protected $primaryKey = 'id';

    /**
     * 提现类型：1：待提现，2可提现
     */
    const CASH_TYPE_PRE_WITHDRAWAL = 1;
    const CASH_TYPE_WITHDRAWAL = 2;

    const CASH_TYPE = [
        self::CASH_TYPE_PRE_WITHDRAWAL => 'expect_cash_amount',
        self::CASH_TYPE_WITHDRAWAL => 'available_cash_amount'
        ];

    protected $fillable=[
        'id',
        'enterprise_id',
        'worker_id',
        'mem_uid',
        'available_cash_amount',
        'expect_cash_amount',
        'partner_id',
        'user_type',
        'frozen_amount',
        'cash_withdrawal'
    ];

    /**
     * 身份类型 1导购 2合伙人
     */
    const  USER_TYPE_WORKER = 1;
    const  USER_TYPE_PARTNER = 2;
}
