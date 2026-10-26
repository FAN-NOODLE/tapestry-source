<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class Address extends Base
{
    const DEF_ADDR = 1; // 默认地址

    protected $table = 'user_address';

    public $primaryKey = 'id';

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
        'created_at',
        'deleted_at',
    ];
}
