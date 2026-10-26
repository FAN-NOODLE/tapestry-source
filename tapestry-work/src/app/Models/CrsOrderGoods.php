<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/9/10
 * Time: 11:47
 */

namespace App\Models;

class CrsOrderGoods extends Base
{
    protected $table = 'crs_order_goods';

    protected $primaryKey = 'id';

    protected $fillable = [
        'id',
        'enterprise_id',
        'order_no',
        'sku_no',
        'spu_id',
        'sku_id',
        'spu_name',
        'type',
        'goods_type',
        'category_name',
        'sku_picture',
        'sku_attribute',
        'quantity',
        'original_price',
        'sell_price',
        'single_discount',
        'sub_price',
        'goods_pay_price',
        'created_at',
        'updated_at',
    ];
}
