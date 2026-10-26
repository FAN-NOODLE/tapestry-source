<?php
/**
 * 常量
 * Created by PhpStorm.
 * User: 刘瑞
 * Date: 2018/11/26
 * Time: 14:52
 */

namespace App\Common;

class Constant
{
    /**
     * 当前标示
     */
    const CLIENT = 'console';

    /**
     * 默认分页数
     */
    const PAGE = 1;

    /**
     * 默认分页数
     */
    const LIMIT = 10;

    /**
     * 成功码
     */
    const SUCCESS_CODE = 0;

    /**
     * 周
     */
    const WEEK = 1;

    /**
     * 今日
     */
    const DAY = 1;

    /**
     * 默认页码
     */
    const DEFAULT_PAGE = 1;

    /**
     * 默认每页数据量
     */
    const DEFAULT_SIZE = 10;

    /**
     * 默认导出最大值
     */
    const MAX_DEFAULT_SIZE = 999;

    /**
     * 是否设置为门店：是
     */
    const DEPARTMENT_IS_STORE_TRUE = 1;
    /**
     * 是否设置为门店：否
     */
    const DEPARTMENT_IS_STORE_FALSE = 0;

    /**
     * 员工在职状态：在职
     */
    const MEMBER_ACTIVE_TRUE = 1;

    /**
     * 员工在职状态：离职
     */
    const MEMBER_ACTIVE_FALSE = 0;

    /**
     * UC员工状态：创建
     */
    const MEMBER_STATUS_CREATED = 1;

    /**
     * UC员工状态：修改
     */
    const MEMBER_STATUS_UPDATED = 2;

    /**
     * UC员工状态：删除
     */
    const MEMBER_STATUS_DELETED = 3;

    /**
     * 离职员工客户是否分配完成状态：已完成
     */
    const CUSTOMER_HANDOVER_COMPLETED_TRUE = 1;

    /**
     * 离职员工客户是否分配完成状态：未完成
     */
    const CUSTOMER_HANDOVER_COMPLETED_FALSE = 2;

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
     * 订单类型：1 普通/新人，2 拼团，3 限购 2
     */
    const ORDER_TYPE_GROUP = 2;

    /**
     * 拼团状态：1 正在拼团，2 拼团成功，3 拼团失败, 4 拼团未付款订单失效
     */
    const GROUP_STATUS_FAIL_PAY = 4;

    /**
     * 省类型
     */
    const PROVINCE_TYPE = 1;

    /**
     * 市类型
     */
    const CITY_TYPE = 2;

    /**
     * 商品属性类型：商品
     */
    const ATTRIBUTE_TYPE_GOODS = 1;

    /**
     * 商品属性类型：赠品
     */
    const ATTRIBUTE_TYPE_GIFT = 2;

    /**
     * 商品属性类型：奖品
     */
    const ATTRIBUTE_TYPE_PRIZE = 3;

    /**
     * 商品属性类型：商品 + 赠品
     */
    const ATTRIBUTE_TYPE_GOODS_GIFT = 4;

    /**
     * 订单列表是否开启支付信息 1 是 2 否 1
     */
    const ORDER_PAYMENT_TRUE = 1;

    /**
     * 订单列表是否开启关联退款 1 是 2 否 1
     */
    const ORDER_REFUND_TRUE = 1;

    /**
     * 订单是否关联赠品信息 1：是，2：否 1
     */
    const ORDER_GIFT_TRUE = 1;

    /**
     * 是否排除SPU选择 1：排除，2：不排除
     */
    const IS_SPU_SELECTION = 1;
    const NOT_SPU_SELECTION = 2;

    /**
     * 以元为单位
     */
    const MONEY_UNIT = 100;

    // 配置设置字段备注
    const CONFIG_KEY_DESC = [
        'mall_open'                 => '好物圈开关',
        'monthly_sales'             => '商品详情页月销量开关',
        'coupon'                    => '优惠券是否显示开关',
        'activity'                  => '商品详情页-活动开关',
        'guarantee'                 => '商品详情页-保障开关',
        'wx_authentication'         => '商品详情页-微信认证开关',
        'seven_days_without_reason' => '商品详情页-七天无理由开关',
        'stock'                     => '商品详情页-商品详情库存开关',
        'mall_button_open'          => '是否开启前台好物推荐按钮（和好物圈开关同时打开才会启用小程序好物推荐接口）',
        'collect_bills_list'        => '购物袋-列表凑单提示开关',
        'category_buy_button'       => '分类页-加购按钮开关',
        'order_empower'             => '下单是否需要授权小程序',
        'collect_bills_full'        => '全场凑单提示开关',
        'direct_sale'               => '保障-官方直营',
        'one_false_pays_ten'        => '保障-假一赔十',
        'sf_free_postage'           => '保障-顺丰包邮',
        'other'                     => '其他',
        'view'                      => '围观',
        'barrage'                   => '弹幕',
        'specs'                     => '规格',
        'show_at_discount_price'    => '是否显示商品详情页优惠价',
        'show_at_mini_program_live' => '是否显示微信小程序直播',
        'category_icon_button'      => '小程序分类是否展示',
        'show_at_live'              => '直播总开关',
        'share_icon'                => '分享图标',
        'share_poster'              => '分享海报',
        'specs_select_distribution' => '规格选择配送',
        'goods_detail_distribution' => '商品详情配送',
        'selling_point'             => '商品卖点',
        'show_head_bar'             => '小程序端商品详情页顶部栏是否显示',
        'store_authentication'      => '店铺保障',
        'distribution'              => '配送方式展示',
        'activity_style'            => '促销风格',
        'comment'                   => '商品详情页-评论开关'
    ];

    // 商品状态转换
    const GOODS_STATUS = [
        '待发布',
        '上架',
        '下架',
        '待上架',
        '待下架',
    ];

    /**
     * 1 导购主动添加  2 支付即导购添加  3 社群添加
     */
    const ACTIVE_ADDITION = 1;
    const GUIDE_ADDITION = 2;
    const COMMUNITY_ADDITION = 3;

    /**
     * 订单是否包含商品信息 1：是，2：否 2
     */
    const ORDER_GOODS_FALSE = 2;

    /**
     * 批量更新商品
     */
    const BATCH_UPDATE_GOODS_FILENAME = '门店批量更新数据模板';
    const BATCH_UPDATE_GOODS_EXT = 'xlsx';
    const BATCH_UPDATE_LIMIT = 500;
    const BATCH_UPDATE_CACHE_TTL = 7200;
    const BATCH_UPDATE_GOODS_HEADER = [
        '门店商品ID',
        '门店商品名称',
        '门店售卖价',
    ];

    const IMPORTED_ERROR_GOODS_FILENAME = '上传失败商品列表';
    const IMPORTED_ERROR_GOODS_HEADER = [
        '门店商品ID',
        '门店商品名称',
        '门店售卖价',
        '原因',
    ];

    const BATCH_UPDATE_STORE_EXT = 'xlsx';
    const IMPORTED_ERROR_STORE_FILENAME = '上传失败门店列表';
    const BATCH_UPDATE_STORE_HEADER = [
        '门店全称',
        '门店简称',
        '门店类型',
        '店长姓名',
        '店长手机号',
    ];

    const IMPORTED_ERROR_STORE_HEADER = [
        '门店全称',
        '门店简称',
        '门店类型',
        '店长姓名',
        '店长手机号',
        '原因',
    ];

    /**
     *  商品状态(1：待上架，2：上架中，3：已下架)
     */
    const STATUS = [
        '',
        '待上架',
        '上架中',
        '已下架'
    ];

    /**
     * 有商城权限
     *
     * @var int
     */
    const HAS_SHOP_AUTH = 1;

    /**
     * 订单详情是否展示"补推Erp"按钮状态：1是，0否
     */
    const SHOW_PUSH_ERP_FALSE = 0;
    const SHOW_PUSH_ERP_TRUE  = 1;

    /**
     *活动名称映射
     */
    const ACTIVITY_NAME = [
        'group'      => ['name'   => '拼团',
                         'code'   => 'PT',
                         'header' => ['ID', '名称', '状态', '时间', '商品名称', '参团人数', '订单数', '机器参与次数', '路径', '物料']
        ]
        ,
        'flash-sale' => ['name'   => '限时购',
                         'code'   => 'XSG',
                         'header' => ['活动ID', '活动名称', '状态', '活动时间', '成交订单数', '成交客户数', '成交总订单金额']
        ]
        ,
        'haggle'     => ['name'   => '砍价',
                         'code'   => 'KJ',
                         'header' => ['活动ID', '活动名称', '状态', '活动时间', '发起人数', '帮砍人数', '成交订单数']
        ],
        'help'       => ['name'   => '好友助力',
                         'code'   => 'ZL',
                         'header' => ['活动ID', '活动名称', '状态', '活动时间', '奖品名称', '发起人数', '参与人数', '助力完成数', '成交订单数', '剩余助力奖品数']
        ]
    ];

    /**
     * 导购类型：1、门店导购，2、社群导购，3、企业导购
     */
    const STORE_GUIDE = 1;
    const GROUP_GUIDE = 2;
    const ENTERPRISE_GUIDE = 3;

    /**
     * 标识导购 1 标识 2 不标识
     */
    const GUIDE_TAG = 1;
    const GUIDE_NO_TAG = 2;


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

    const ROLE_LIST = [
        "guide_enterprise_admin", // 智慧导购企业管理
        "console_community_guide", // 管理后台社群导购
        "guide_community_guide",  // 智慧导购社群导购
        "console_enterprise_guide", // 管理后台企业导购
        "guide_enterprise_guide", // 智慧导购企业导购
        "console_store_manager", // 管理后台店长
        "guide_store_manager", // 智慧导购店长
        "guide_store_member" // 智慧导购店员
    ];

    /**
     * 批量改价队列
     */
    const MULTI_CHANGE_PRICE_QUEUE_NAME = 'multi_change_price_queue_name';

    /**
     * 未上架商品列表下载队列
     */
    const NOT_SALE_LIST_DOWNLOAD = 'not_sale_list_download';

    /**
     * 批量更改商品信息 excel数据保存key
     */
    public const GOODS_BATCH_CHANGE_SAVE_LIST = 'crs:console:goods-batch-change-save-list:%s:%s';
    /**
     * 批量更改商品信息 进度保存key
     */
    public const GOODS_BATCH_CHANGE_PROGRESS = 'crs:console:goods-batch-change-progress:%s:%s';

    /**
     * 批量更改商品信息 系统进度保存key
     */
    public const GOODS_BATCH_CHANGE_PROGRESS_SYSTEM = 'crs:console:goods-batch-change-progress-system:%s:%s';

    /**
     * 批量更改商品信息 结果及状态保存key
     *
     */
    public const GOODS_BATCH_CHANGE_RESULT = 'crs:console:goods-batch-change-result:%s:%s';

    /**
     * 企业配置信息
     */
    public const ENTERPRISE_CONFIG = 'crs:console:enterprise-config:%s';

    /**
     * 微信优惠券详情
     */
    public const WX_COUPON_DETAIL = 'crs:console:wx-coupon-detail:%s';
}
