<?php
/**
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/11/23
 * Time: 17:45
 */

namespace App\Models;


class ShopOrderGoods extends Base
{
    protected $table = 'shop_order_goods';

    public $fillable = [
        'id',
        'enterprise_id',
        'order_no',
        'spu_id',
        'sku_id',
        'spu_name',
        'category_name',
        'sku_picture',
        'type',
        'object_id',
        'spu_picture',
        'sku_attribute',
        'quantity',
        'sell_price',
        'activity_price',
        'sub_price',
        'sku_no',
        'goods_type',
        'created_at',
        'updated_at',
        'deleted_at'
    ];

}
