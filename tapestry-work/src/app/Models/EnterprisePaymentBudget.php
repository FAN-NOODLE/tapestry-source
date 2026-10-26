<?php


namespace App\Models;


class EnterprisePaymentBudget extends Base
{

    /**
     * 充值
     */
    const INCOME_TYPE_IN = 1;

    /**
     * 扣款
     */
    const INCOME_TYPE_OUT = 2;

    // 1=充值 2=活动扣款 3=活动返款
    const RECHARGE_TYPE = 1;
    const EVENT_DEDUCTION_TYPE = 2;
    const ACTIVITY_REFUND_TYPE = 3;

    protected $table = 'enterprise_payment_budget';

    protected $guarded = [];


}