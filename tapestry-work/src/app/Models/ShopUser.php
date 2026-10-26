<?php

namespace App\Models;


class ShopUser extends Base
{
    /**
     * 客户状态 0、未授权，1、已授权，2、未成交，3、已成交
     */
    const UN_AUTH  = 0;

    const AUTH     = 1;

    const UN_TRADE = 2;

    const TRADED   = 3;

    protected $table = 'shop_user';


    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'enterprise_id',
        'account_id',
        'union_id',
        'open_id',
        'third_member_no',
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
        'sex',
        'level_id',
        'level_update_at',
        'user_amount'
    ];

    public function addresses()
    {
        return $this->hasMany(ShopUserAddress::class, 'user_id', 'id');
    }
}

