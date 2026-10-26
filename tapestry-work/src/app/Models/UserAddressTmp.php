<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/10
 * Time: 16:28
 */

namespace App\Models;

class UserAddressTmp extends Base
{
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'user_address_tmp';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
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
        'updated_at',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'user_id', 'province', 'city', 'area', 'address',
        'province_code', 'city_code', 'area_code', 'zip_code', 'default'];
}
