<?php

namespace App\Common;

class RedisKey
{
    // 企业信息 redis 缓存清理
    const ENTERPRISE_HASH_INFO = 'shopEnterpriseInfo:';
    // 企业hash详情缓存
    const ENTERPRISE_DOMAIN_INFO = 'crs:enterprise:domain:info:';
    // 企业信息缓存
    const ENTERPRISE_ID_INFO = 'crs:enterprise:id:info:';
    // 企业应用缓存
    const ENTERPRISE_WORK_PLUGIN = 'crs:enterprise:work:plugin:';
    // 企业配置缓存
    const ENTERPRISE_CONFIG_INFO = 'crs:enterprise:config:info:';
    // 企业corpId详情缓存
    const ENTERPRISE_CORP_ID_INFO = 'crs:enterprise:corpid:info:';
    // 导购业绩配置
    const ENTERPRISE_GUIDE_SETTING = 'crs:gw-shop:enterprise_guide_setting:';
    // 海报二维码缓存
    const QR_CODE = 'crs:work:qrCode:';
    // 小程序授权信息缓存
    const REDIS_SHOP_AUTH_PREFIX = 'shopAuthInfo:';
    // 小程序商城主色调
    const SHOP_COLOR = 'crs:work:shopColor:';
    // 门店列表缓存
    const STORE_LIST = 'crs:work:store:list:';
    // store 门店员工信息
    const STORE_WORKER_INFO = 'crs:store:worker:info';
    // store 用户门店导购关联信息
    const USER_STORE_GUIDE_INFO = 'crs:user:store:guide:info';

    // 顾客关系绑定成功消息发送标示
    const CUSTOMER_RELATION_BOUND_MESSAGE_FLAG = 'crs:work:customer-relation-bound-message-flag:';

    // 企业顾客统计数据初始化标识
    const ENTERPRISE_CUSTOMER_STATISTIC_INIT_FLAG = 'crs:work:enterprise-customer-statistic-init-flag:';

    // 企业顾客累计统计数据更新标识
    const ENTERPRISE_CUMULATIVE_CUSTOMER_STATISTIC_UPDATE_FLAG = 'crs:work:enterprise-cumulative-customer-statistic-update-flag:';

    // 同步业务架构至企微架构
    const SYNC_STORE_INFO_BY_CHAT_TIME = 'crs:work:store:sync-store-info-by-chat-time:';

    // 同步组织架构至业务架构
    const SYNC_STORE_INFO_BY_DEPARTMENT_TIME = 'crs:work:store:sync-store-info-by-department-time:';

    // UC token缓存
    const UC_API_ACCESS_TOKEN = 'crs:work:uc:token';

    // 职位上次更新时间
    const UPDATE_STORE_MANAGER_BY_POSITION_TIME = 'crs:work:member:update-store-manager-by-position-time:';

    // 在职继承参数缓存，避免重复提交
    const CUSTOMER_INHERIT_PARAMS = 'crs:work:customer:inherit:params:';

    // 在职继承异步缓存，等待处理后继续提交
    const CUSTOMER_INHERIT_WAIT_TIME = 'crs:work:customer:inherit:wait:time:';

}
