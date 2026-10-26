<?php

namespace App\Common;

/**
 * Created by PhpStorm.
 * User: long
 * Date: 18-11-27
 * Time: 下午4:32
 */
class Constant
{
    // 应用安装状态 1=已安装 2=未安装
    const ENTERPRISE_INSTALL_STATUS     = 1;
    const ENTERPRISE_UNINSTALLED_STATUS = 2;

    // 企业应用标示数组
    const APP_IDENTITY = [
        'coupon', // 礼劵
        'novice', // 新人优惠
        'sec_kill', // 秒杀限时购
        'group', // 拼团
        'activity', // 优惠
    ];

    /**
     * 默认页码
     */
    const DEFAULT_PAGE = 1;

    /**
     * 默认每页数据量
     */
    const DEFAULT_LIMIT = 10;
}
