<?php
/**
 * Created by PhpStorm.
 * User: zq
 * Date: 19-4-3
 * Time: 下午2:43
 */

namespace App\Models;

class ShopUserAddress extends Base
{
    protected $table = 'shop_user_address';

    protected $fillable = [
        'id',
        'enterprise_id',
        'user_id',
        'province',
        'city',
        'area',
        'address',
        'province_code',
        'city_code',
        'area_code',
        'zip_code',
        'default',
        'delivery_name',
        'delivery_mobile',
        'is_auth_address'
    ];
}