<?php

namespace App\Common;

class RedisKey
{
    /**
     * 登录企业存储
     *
     * @var string
     */
    const LOGIN_MANAGER_ENTERPRISES = 'crs:assistant:login_manager_enterprises';

    /**
     * 数据罗盘部分缓存
     */
    const DATA_COMPASS_SUMMARY = 'crs:assistant:data_compass_summary';

    /**
     * 绩效榜单明星榜单缓存
     */
    const PERFORMANCE_LIST_STAR = 'crs:assistant:performance_list_star:';

    /**
     * 绩效榜单SW明星榜单缓存
     */
    const PERFORMANCE_LIST_SW_STAR = 'crs:assistant:performance_list_sw_star:';
    /**
     * 绩效榜单KS明星榜单缓存
     */
    const PERFORMANCE_LIST_KS_STAR = 'crs:assistant:performance_list_ks_star:';
}
