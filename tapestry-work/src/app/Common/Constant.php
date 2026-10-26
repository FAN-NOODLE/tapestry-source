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
    // 任务完成状态 1=未完成 2=已完成
    const TASK_EXECUTING_STATUS   = 1;
    const TASK_EXE_SUCCESS_STATUS = 2;

    // 任务是否提醒 0=未提醒 1=已提醒
    const TASK_Not_REMIND_STATUS = 0;
    const TASK_REMIND_STATUS     = 1;

    // 任务显示状态 1=显示 2=隐藏
    const TASK_SHOW_STATUS = 1;
    const TASK_HIDE_STATUS = 2;

    /**
     * 每次uc获取人数
     */
    const UC_MEMBER_NUMBER = 100;

    /**
     * UC店长
     */
    const UC_MANAGER = '店长';
    /**
     * 按订单金额
     */
    const ORDER_SORT = 1;

    /**
     * 顾客是否会员，1-会员，2-非会员
     */
    const IS_MEMBER  = 1;
    const NOT_MEMBER = 0;

    /**
     * 按订单数量排序
     */
    const ORDER_NUMBER = 2;

    /**
     * 今日
     */
    const DAY = 1;

    /**
     * 本周
     */
    const WEEK = 2;

    /**
     * 本月
     */
    const MONTH = 3;

    /**
     * 默认页码
     */
    const DEFAULT_PAGE = 1;

    /**
     * 默认每页数据量
     */
    const DEFAULT_SIZE = 10;

    /**
     * 每页最大数据量
     */
    const DEFAULT_SIZE_MAX = 9999;

    /**
     * 离职员工客户分配状态：未分配
     */
    const CUSTOMER_HANDOVER_STATUS_FALSE = 1;

    /**
     * 离职员工客户分配状态：已分配
     */
    const CUSTOMER_HANDOVER_STATUS_TRUE = 2;

    /**
     * 客户分配确认状态：已确认
     */
    const CUSTOMER_HANDOVER_ACK_STATUS_TRUE = 1;

    /**
     * 客户分配确认状态：待确认
     */
    const CUSTOMER_HANDOVER_ACK_STATUS_FALSE = 2;

    /**
     * 离职员工客户是否分配完成状态：已完成
     */
    const CUSTOMER_HANDOVER_COMPLETED_TRUE = 1;

    /**
     * 离职员工客户是否分配完成状态：未完成
     */
    const CUSTOMER_HANDOVER_COMPLETED_FALSE = 2;

    /**
     * 是否专属导购：是
     */
    const MEMBER_IS_EXCLUSIVE_TRUE = 1;

    /**
     * 是否专属导购：否
     */
    const MEMBER_IS_EXCLUSIVE_FALSE = 2;

    /**
     * 查询出所有成员包括删除人员
     */
    const MEM_ALL_UC = 1;

    /**
     * UC的删除人员
     */
    const MEM_DEALETE_UC = 3;

    /**
     * crs小程序开通
     */
    const CRS_OPEN = 1;

    /**
     * crs小程序关闭
     */
    const CRS_CLOSE = 2;

    /**
     * 商城小程序开通
     */
    const STORE_OPEN = 1;

    /**
     * 商城小程序关闭
     */
    const STORE_CLOSE = 2;

    /**
     * 分表表名
     */
    const TABLE_NAME = 'enterprise_trend_data_';

    /**
     * 客户状态 0、未授权，1、已授权，2、未成交，3、已成交
     */
    const UN_AUTH  = 0;
    const AUTH     = 1;
    const UN_TRADE = 2;
    const TRADED   = 3;

    /**
     * 1 导购主动添加  2 支付即导购添加  3 社群添加
     */
    const ACTIVE_ADDITION    = 1;
    const GUIDE_ADDITION     = 2;
    const COMMUNITY_ADDITION = 3;

    /**
     * 时间段类型 1=今日 2=最近7天 3=最近30天 4 指定日期区间 5=昨日
     */
    const ORDER_TODAY_TYPE       = 1;
    const ORDER_LAST_SEVEN_DAYS  = 2;
    const ORDER_LAST_THIRTY_DAYS = 3;
    const ORDER_APPOINT_DAYS     = 4;
    const ORDER_YESTERDAY_DAYS   = 5;

    /**
     * 业绩排行类型 1=按金额数 2=按订单数
     */
    const ORDER_ACHIEVEMENT_MONEY = 1;
    const ORDER_ACHIEVEMENT_NUM   = 2;

    /**
     * 业绩排行排序 1=整序 2=倒叙
     */
    const ORDER_ACHIEVEMENT_ASC  = 1;
    const ORDER_ACHIEVEMENT_DESC = 2;


    /**
     * 默认页码
     */
    const PAGE = 1;

    /**
     * 最大获取
     */
    const MAX_SIZE = 1000;

    /**
     * 默认每页数据量
     */
    const SIZE = 10;

    /**
     * 企业商城版本 1-标准版（旧版） 2-中心商城有门店版 3-中心商城无门店版 4-官方直营版
     */
    const NORMAL_VERSION       = 1;
    const CENTER_STORE_VERSION = 2;
    const CENTER_VERSION       = 3;
    const NO_CENTER_VERSION    = 4;

    // 上架中
    const STATUS_UP = 2;

    // 普通商品
    const NORMAL_GOODS = 1;

    // 活动商品
    const ACTIVITY_GOODS = 2;

    /**
     * 1 商品 2 优惠券 6 好友助力-付邮 7 好友助力-买赠  12 限时购
     * 13 助力抽奖 14 拼团 15 砍价
     */
    const MARKETING_GOODS           = 1;
    const MARKETING_COUPON          = 2;
    const MARKETING_HELPER_POST     = 6;
    const MARKETING_HELPER_BUY_GIVE = 7;
    const MARKETING_SECKILL         = 12;
    const MARKETING_LOTTERY         = 13;
    const MARKETING_GROUP           = 14;
    const MARKETING_BARGAIN         = 15;

    /**
     * 1 发送 2下载
     */
    const OPERATE_TYPE_SEND     = 1;
    const OPERATE_TYPE_DOWNLOAD = 2;

    /**
     * 小程序应用类型: 1商城小程序、2企业微信小程序 1
     */
    const APP_TYPE = 1;

    /**
     * 1 快递 2 自提 3 快递+自提
     */
    const GOODS_EXPRESS                     = 1;
    const GOODS_SELF_EXTRACTING             = 2;
    const GOODS_SELF_EXTRACTING_AND_EXPRESS = 3;

    /**
     * 吊牌价 1 需要
     */
    const SIGN_PRICE = 1;

    /**
     * 顾客查询条件：订单数最小值和最大值
     */
    const CUSTOMER_ORDER_NUM_MIN = 0;
    const CUSTOMER_ORDER_NUM_MAX = 999;
    /**
     * 顾客查询条件：订单金额最小值和最大值
     */
    const CUSTOMER_ORDER_AMOUNT_MIN = 0;
    const CUSTOMER_ORDER_AMOUNT_MAX = 999999;

    /**
     * 数据库选择常量
     */
    const DB_WORK       = 'work';
    const DB_GW_SHOP    = 'gw_shop';
    const DATABASE_CONF = [
        'prod' => [
            self::DB_WORK    => 'crs_prod_work',
            self::DB_GW_SHOP => 'crs_prod_gw_shop',
        ],
        'test' => [
            self::DB_WORK    => 'crs_test_work',
            self::DB_GW_SHOP => 'crs_test_gw_shop',
        ],
        'dev'  => [
            self::DB_WORK    => 'crs_test_work',
            self::DB_GW_SHOP => 'crs_test_gw_shop',
        ]
    ];

    /**
     * 导购类型：1、门店导购，2、社群导购，3、企业导购  4.企业管理 5、分销合伙人
     */
    const STORE_WORKER_TYPE_STORE      = 1;
    const STORE_WORKER_TYPE_COMMUNITY  = 2;
    const STORE_WORKER_TYPE_ENTERPRISE = 3;
    const STORE_WORKER_TYPE_MANAGER    = 4;
    const STORE_WORKER_TYPE_PARTNER    = 5;
    /**
     * 数据库中is_xxx常量：1=是, 2=否
     */
    const YES = 1;
    const NO  = 2;

    /**
     * 1 内购 2 导购
     */
    const STAFF_TYPE = 1;
    const GUIDE_TYPE = 2;

    /**
     * 员工状态：1=在职，2=离职
     */
    const IN_SERVICE  = 1;
    const OUT_SERVICE = 2;

    /**
     * 企业标识常量 boss.enterprise_work_plugin.identifier
     */
    const ENTERPRISE_IDENTIFIER_STORE_ASSISTANT = 'crs_store_assistant';

    /**
     * 订单状态
     */
    const ORDER_STATUS_WAIT_DELIVER = 40;

    /**
     * UC接口exchange-info-list参数限制个数
     */
    const UC_EXCHANGE_INTERFACE_PARAM_MAX = 50;

    /**
     * 列表默认返回数据格式
     *
     * @var array
     */
    const LIST_DEFAULT_JSON = [
        'data' => [],
        'meta' => [
            'page'  => 1,
            'size'  => 0,
            'total' => 0
        ]
    ];


    /**
     * 导购顾客数据排行榜类型
     */

    const GUIDE_CUSTOMER_RANK_CUSTOMER     = 1;
    const GUIDE_CUSTOMER_RANK_NEW_CUSTOMER = 2;
    const GUIDE_CUSTOMER_RANK_ACCESS       = 3;
    const GUIDE_CUSTOMER_RANK_DEAL         = 4;

    /**
     * 一天的秒数
     */
    const ONE_DAY_SEC = 86400;

    /**
     * 导购统计排行榜数
     */
    const GUIDE_RANK_NUM = 50;

    /**
     * 企业关联数据范围类型：0-全部门店 1-区域 2-门店  3-商城  4 部门
     */
    const DATA_RANGE_ALL_STORE  = 0;
    const DATA_RANGE_AREA       = 1;
    const DATA_RANGE_STORE      = 2;
    const DATA_RANGE_MALL       = 3;
    const DATA_RANGE_DEPARTMENT = 4;

    /**
     * 业绩归属  1 自己 2 他人
     */
    const ACHIEVEMENT_ASCRIPTION_SELF   = 1;
    const ACHIEVEMENT_ASCRIPTION_OTHERS = 2;

    /**
     * 添加方式 1 导购名片 2 其他方式
     */
    const ADD_TYPE_BUSINESS_CARD       = 1;
    const ADD_TYPE_BUSINESS_CARD_OTHER = 2;

    /**
     * 全员分销 是否上新  1 是 2 否
     */
    const NEW_GOODS_YES = 1;
    const NEW_GOODS_NO  = 2;

    // 分包之后的小程序路径
    const MARKETING_PATH_NEW = [
        Constant::MARKETING_GOODS           => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_COUPON          => 'other/receive-coupons/receive-coupons',
        Constant::MARKETING_HELPER_POST     => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_HELPER_BUY_GIVE => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_SECKILL         => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_LOTTERY         => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_GROUP           => 'goods/goods-detail/goods-detail',
        Constant::MARKETING_BARGAIN         => 'goods/goods-detail/goods-detail',
    ];

    /**
     * 门店员工异步生成二维码和太阳码
     */
    const STORE_WORKER_GENERATE_QR_CODE_JOB_USE_WORK = 'store_worker_generate_qr_code_job_use_work';

    /**
     * 门店相关角色
     *
     * @var array
     */
    const GUIDE_IN_STORE_ROLE = [
        self::GUIDE_STORE_MEMBER,
        self::GUIDE_STORE_MANAGER,
    ];

    /**
     * 无门店类角色
     *
     * @var array
     */
    const GUIDE_NOT_STORE_ROLE = [
        self::GUIDE_ENTERPRISE_GUIDE,
        self::GUIDE_COMMUNITY_GUIDE,
        self::GUIDE_ENTERPRISE_ADMIN,
    ];

    // 内置角色
    const GUIDE_STORE_MEMBER       = 'guide_store_member';
    const GUIDE_STORE_MANAGER      = 'guide_store_manager';
    const GUIDE_ENTERPRISE_GUIDE   = 'guide_enterprise_guide';
    const GUIDE_COMMUNITY_GUIDE    = 'guide_community_guide';
    const GUIDE_ENTERPRISE_ADMIN   = 'guide_enterprise_admin';
    const CONSOLE_ENTERPRISE_ADMIN = 'console_enterprise_admin';

    /**
     * 导购初始化角色名
     *
     * @var array
     */
    const GUIDE_DEFAULT_ROLE = [
        self::GUIDE_STORE_MEMBER,
        self::GUIDE_STORE_MANAGER,
        self::GUIDE_ENTERPRISE_GUIDE,
        self::GUIDE_COMMUNITY_GUIDE,
        self::GUIDE_ENTERPRISE_ADMIN,
    ];

    /**
     * 普通用户角色
     *
     * @var int
     */
    const LOGIN_IDENTITY_USER = 1;

    /**
     * 导购角色
     *
     * @var int
     */
    const LOGIN_IDENTITY_GUIDE = 2;

    /**
     * 合伙人角色
     * @var int
     */
    const LOGIN_IDENTITY_PARTNER = 10;

    /**
     * /**
     * 员工购买商品有佣金(1有佣金)
     *
     * @var int commission_type
     */
    const GUIDE_COMMISSION_TYPE_YES = 1;

    /**
     * 客户分享统计业绩(2 计算业绩)
     *
     * @var int fission_type
     */
    const GUIDE_FISSION_TYPE = 2;

    /**
     * 合伙人状态
     */
    const PARTNER_STATUS_ENABLE  = 1;
    const PARTNER_STATUS_DISABLE = 0;

    /**
     * 会员统计
     */
    // 是否是会员(是)
    const IS_MEMBER_STATISTICS = 1;

    // 是否是会员（否）
    const NON_MEMBER_STATISTICS = 2;

    /**
     * add_way表示添加客户的来
     *
     *  0    未知来源
     * 1    扫描二维码
     * 2    搜索手机号
     * 3    名片分享
     * 4    群聊
     * 5    手机通讯录
     * 6    微信联系人
     * 7    来自微信的添加好友申请
     * 8    安装第三方应用时自动添加的客服人员
     * 9    搜索邮箱
     * 201    内部成员共享
     * 202    管理员/负责人分配
     */
    const ADD_WAY_UNKNOWN_SOURCE      = 0;
    const ADD_WAY_SCAN_QR_CODE        = 1;
    const ADD_WAY_SEARCH_MOBILE       = 2;
    const ADD_WAY_CARD_SHARE          = 3;
    const ADD_WAY_GROUP_CHAT          = 4;
    const ADD_WAY_MOBILE_ADDRESS_BOOK = 5;
    const ADD_WAY_WECHAT_CONTACT      = 6;
    const ADD_WAY_WECHAT_FRIENDS      = 7;
    const ADD_WAY_OTHER_APP           = 8;
    const ADD_WAY_SEARCH_EMAIL        = 9;
    const ADD_WAY_IASIDE_SHARE        = 201;
    const ADD_WAY_MANAGER_ALLOCATION  = 202;

    /**
     *  其他方式
     */
    const ADD_WAY_OTHER = -1;

    /**
     * 类型数组
     */
    const ADD_WAY_TYPE = [
        self::ADD_WAY_UNKNOWN_SOURCE      => '未知来源',
        self::ADD_WAY_SCAN_QR_CODE        => '扫描二维码',
        self::ADD_WAY_SEARCH_MOBILE       => '搜索手机号',
        self::ADD_WAY_CARD_SHARE          => '名片分享',
        self::ADD_WAY_GROUP_CHAT          => '群聊',
        self::ADD_WAY_MOBILE_ADDRESS_BOOK => '手机通讯录',
        self::ADD_WAY_WECHAT_CONTACT      => '微信联系人',
        self::ADD_WAY_WECHAT_FRIENDS      => '来自微信的添加好友申请',
        self::ADD_WAY_OTHER_APP           => '安装第三方应用时自动添加的客服人员',
        self::ADD_WAY_SEARCH_EMAIL        => '搜索邮箱',
        self::ADD_WAY_IASIDE_SHARE        => '内部成员共享',
        self::ADD_WAY_MANAGER_ALLOCATION  => '管理员/负责人分配',
    ];


}
