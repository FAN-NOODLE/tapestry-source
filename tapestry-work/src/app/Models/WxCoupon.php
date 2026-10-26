<?php

namespace App\Models;

class WxCoupon extends Base
{
    /**
     * 表名
     * @var string
     */
    protected $table = 'wx_coupon';

    /**
     * 主键
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'shop_wx_coupon_id',
        'enterprise_id',
        'coupon_name',
        'coupon_value',
        'coupon_minimum'
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'shop_wx_coupon_id', 'enterprise_id',
        'coupon_name', 'coupon_value', 'coupon_minimum'];
}
