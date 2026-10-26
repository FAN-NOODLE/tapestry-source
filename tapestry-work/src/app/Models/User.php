<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

class User extends Base
{
    /**
     * 已授权
     * @var int
     */
    const IS_AUTH = 1;

    /**
     * 未成交
     * @var int
     */
    const IS_ADDITION = 2;

    /**
     * 已成交
     * @var int
     */
    const IS_ADDED = 3;

    /**
     * 已分配
     * @var int
     */
    const IS_ALLOCATED_YES = 1;

    /**
     * 未分配
     * @var int
     */
    const IS_ALLOCATED_NO = 2;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'user';

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
        'union_id',
        'open_id',
        'nick_name',
        'real_name',
        'real_name_index',
        'mobile',
        'birthday',
        'avatar',
        'country',
        'province',
        'city',
        'register_at',
        'auth_at',
        'type',
        'sex',
        'user_amount',
        'sync_source',
        'created_at',
        'updated_at',
    ];

    /**
     * @var array
     */
    public static $selectField = ['id', 'enterprise_id', 'union_id', 'open_id', 'nick_name', 'real_name', 'real_name_index',
                                  'mobile', 'birthday', 'avatar', 'country', 'province', 'city', 'register_at', 'auth_at', 'type', 'sex', 'user_amount', 'sync_source'];

    /**
     * 用户优惠券列表
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function coupons()
    {
        return $this->hasMany(UserCoupon::class, 'user_id', 'id');
    }
}
