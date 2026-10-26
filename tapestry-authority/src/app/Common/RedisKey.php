<?php

namespace App\Common;

class RedisKey
{
    // 企业配置
    const ENTERPRISE_CONFIG = 'crs:shop:enterprise:config:list';
    // 企业 access_token
    const ACCESS_TOKEN_BY_ENTERPRISE_ID = 'crs:shop:access:token';

    // 小程序授权的账号类目id
    const MINI_WX_CATEGORY_IDS = 'crs:shop:mini_wx_category_ids:';

    const SEND_SUBSCRIBE_MESSAGE = 'crs:gw:shop:send:subscribe:message';


    // 订阅消息提醒
    const USER_SUBSCRIBE_NOTICE_STATUS = 'crs:authority:subscribe-notice-status:cache:uid:';


    // 企业订阅消息场景
    const ENTERPRISE_MESSAGE_SCENE = 'crs:authority:enterprise-message-scene:cache:eid:';

    // 企业用户召回配置策略
    const ENTERPRISE_MESSAGE_SCENE_COUNT = 'crs:authority:message-scene:cache:eid:';

    const ENTERPRISE_THIRD_PART_ACCOUNT_WX_CATEGORY = 'crs:authority:enterprise-with-wx-category:cache:eid:';

    const SUBSCRIBE_TOUCH_DATA_CACHE = 'crs:shop:index:subscribe:cache:eid:';

    const ACCESS_TOKEN_LOCK = 'access_token_lock:enterprise_id_%s:type_%s';

    const RESET_PASSWORD_TIME = 'crs:authority:reset-password:enterprise_id_%s:mobile_%s';
}
