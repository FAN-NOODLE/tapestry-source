<?php


namespace App\Models;


class StaffOrderRelation extends Base
{
    /**
     * is_my_order 是否是员工自己的订单（1=是、2=否）
     */
    const IS_MY_ORDER_YES = 1;
    const IS_MY_ORDER_NO  = 2;


    protected $table = 'staff_order_relation';

    protected $guarded = [];

}