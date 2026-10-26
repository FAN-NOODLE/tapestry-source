<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class GroupOrder extends Base
{
    protected $table = 'group_order';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'qw_group_id',
        'group_user_id',
        'union_id',
        'order_no',
        'order_pay_time',
        'refund_money',
        'pay_money',
        'order_money',
        'created_at',
        'updated_at',
    ];
}
