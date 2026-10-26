<?php

namespace App\Models;

class ActivityReserveRecord extends Base
{
    // 是否核销，0未核销,1已核销,已过期
    const IsCheckYes = 1;
    const IsCheckNo = 0;
    const CheckWait = 2;
    //白名单列表
    const WHITELISTS = [
        18735931251,
        16621128726,
        13992092140,
        18271185339
    ];
    //管理导购白名单
    const MANAGEWHITELISTS = [
        18735931251,
        16621128726,
        13992092140,
        18271185339
    ];
    public $primaryKey = 'id';


    protected $table = 'activity_reserve_record';
    protected $fillable = [
        'id',
        'enterprise_id',
        'activity_id',
        'unionid',
        'openid',
        'user_name',
        'avatar',
        'mobile',
        'qr_code',
        'activity_code',
        'reserve_time',
        'is_check',
        'check_time',
        'member_id',
        'member_name',
        'version'
    ];


}