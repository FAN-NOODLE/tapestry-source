<?php
/**
 * Created by PhpStorm.
 * User: 王鹏飞
 * Date: 2019/5/20
 * Time: 下午10:03
 */

namespace App\Models;


class ShopOrderTrace extends Base
{
    protected $table = 'shop_order_trace';

    protected $fillable = [
        'id',
        'order_id',
        'event',
        'event_desc',
        'from_order_status',
        'to_order_status',
        'created_at',
        'updated_at',
        'deleted_at',
    ];
}
