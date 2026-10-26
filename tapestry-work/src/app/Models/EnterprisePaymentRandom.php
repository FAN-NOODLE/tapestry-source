<?php


namespace App\Models;


class EnterprisePaymentRandom extends Base
{

    /**
     * 是否占用 1=否 2=是
     */
    const IS_RECEIVE_NO = 1;
    const IS_RECEIVE_YES = 2;

    protected $table = 'enterprise_payment_random';

    protected $guarded = [];

}