<?php

namespace App\Common;

class Constant
{
    /**
     * 默认页码
     */
    const DEFAULT_PAGE = 1;

    /**
     * 默认每页数据量
     */
    const DEFAULT_LIMIT = 10;

    /**
     * 数据不存在
     */
    const NOT_FOUND_CODE = 700404;

    /**
     * 小程序应用类型: 1商城小程序、2企业微信小程序 1
     */
    const APP_TYPE = 1;

    const STORE_MEMBER = 'guide_store_member';
    const GUIDE_STORE_MANAGER = 'guide_store_manager';
    const CONSOLE_STORE_MANAGER = 'console_store_manager';
    const GUIDE_ENTERPRISE_GUIDE = 'guide_enterprise_guide';
    const CONSOLE_ENTERPRISE_GUIDE = 'console_enterprise_guide';
    const GUIDE_COMMUNITY_GUIDE = 'guide_community_guide';
    const CONSOLE_COMMUNITY_GUIDE = 'console_community_guide';
    const GUIDE_ENTERPRISE_ADMIN = 'guide_enterprise_admin';
    const CONSOLE_ENTERPRISE_ADMIN = 'console_enterprise_admin';

    // 短信发送限制次数
    const MESSAGE_LIMIT_COUNT = 10;

    // 后台登录token盐值
    const TOKEN_SALT = 'REDACTED_FOR_SOURCE_SHARING';

    // 短信白名单
    const MESSAGE_MOBILE_WHITE_LIST = [
        '18829253910' => '徐利利',
        '18821751802' => '韩磊',
        '18535959118' => '段龙',
        '15249279779' => '邱阿朋',
        '15110188985' => '王乐',
        '18229066928' => '翟俊亮',
        '18591927513' => '赵鑫',
        '18829297422' => '穆萧',
        '15229392991' => '张阳阳',
        '18092066590' => '赵洋',
        '18616749050' => '姜涛',
    ];

}
