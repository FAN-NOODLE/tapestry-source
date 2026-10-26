<?php

namespace App\Models;

class Enterprise extends Base
{
    /**
     * 企业商城版本 1:标准版（旧版）,2:中心商城版(有门店), 3:中心商城, 4:直营连锁门店
     */
    const NORMAL_VERSION = 1;
    const VERSION_CENTER_SHOP_STORE = 2;
    const VERSION_CENTER_SHOP = 3;
    const VERSION_DIRECT_CHAIN_STORE = 4;

    // 启用
    const Enable = 1;

    // 企业可用系统定义
    const SYSTEM = [
        1 => '管理后台',
        2 => '智慧导购小程序',
        //3 => '天擎CDP'
    ];

    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise';

    protected $guarded = [];

    public static function getDetail($eID)
    {
        return self::query()->where('id', '=', $eID)->first();
    }
}