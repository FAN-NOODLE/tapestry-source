<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class GroupOrderGoods extends Base
{
    protected $table = 'group_order_goods';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'qw_group_id',
        'group_user_id',
        'union_id',
        'order_no',
        'spu_id',
        'sku_id',
        'spu_name',
        'category_name',
        'sku_picture',
        'spu_picture',
        'quantity',
        'refund_quantity',
        'sell_price',
        'sub_price',
        'created_at',
        'updated_at',
    ];
}
