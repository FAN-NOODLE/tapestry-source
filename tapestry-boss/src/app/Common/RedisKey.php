<?php

namespace App\Common;

class RedisKey
{
    /**
     * 模型单记录缓存 存json
     * %s 模型名
     * %s where条件json
     */
    const MODEL_ARRAY = 'crs:boss:%s:%s';

    /**
     * 小程序授权信息 存json
     * %s 企业ID
     */
    const AUTH_APP_INFO = 'crs:boss:app:%s';

    /**
     * 商城授权信息
     * %s 企业hash
     */
    const AUTH_APP_INFO_HASH = 'crs:boss:app-hash:%s';

    // 企业配置
    const ENTERPRISE_CONFIG = 'crs:shop:enterprise:config:list';
    // 企业 access_token
    const ACCESS_TOKEN_BY_ENTERPRISE_ID = 'crs:shop:access:token';
    // shop小程序授权
    const SHOP_AUTH_DETAIL = 'crs:shop:mini:auth:detail';
    // 用户登录信息缓存
    const REDIS_SHOP_AUTH_PREFIX = 'shopAuthInfo:';
    // 企业信息 redis 缓存清理
    const ENTERPRISE_HASH_INFO = 'shopEnterpriseInfo:';
    // 企业hash详情缓存
    const ENTERPRISE_DOMAIN_INFO = 'crs:enterprise:domain:info:';
    // 企业corpId详情缓存
    const ENTERPRISE_CORP_ID_INFO = 'crs:enterprise:corpid:info:';

    /**
     * shop端 新人优惠 状态 redis缓存清理
     */
    const REDIS_APP_NOVICE = 'shopNoviceStatusInfo:';

    /**
     * shop端 拼团 状态 redis缓存清理
     */
    const REDIS_APP_GROUP = 'shopGroupStatusInfo:';

    /**
     * shop端 秒杀 状态 redis缓存清理
     */
    const REDIS_APP_SECKILL = 'shopNoviceStatusInfo:';

    /**
     * shop端 企业信息 redis 缓存清理
     */
    const REDIS_ENTERPRISE_INFO = 'shopEnterpriseInfo:';

    /**
     * 企业付费信息是否更新
     */
    const REDIS_ENTERPRISE_PAYMENT_HAS_NEW = 'enterprisePaymentHasNew';

    /**
     * 小程序授权信息（新）
     */
    const MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID = 'miniProgramAuthInfo:';

    /**
     *
     */
    const ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID = 'enterpriseWorkPlugin:';

}
