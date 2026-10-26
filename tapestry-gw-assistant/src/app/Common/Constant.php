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
     * 默认分页数
     *
     * @var int
     */
    const PAGE = 1;

    /**
     * 默认分页数
     *
     * @var int
     */
    const LIMIT = 10;

    /**
     * 登录token 有效时长
     * 单位：秒(s)
     *
     * @var int
     */
//    const LOGIN_TOKEN_EXPIRED =  1800;
    const LOGIN_TOKEN_EXPIRED = 2 * 60 * 60;

    /**
     * 是否是店长
     *
     * @var int
     */
    const IS_BOSS_YES = 1;

    /**
     * 员工状态(在职)
     *
     * @var int
     */
    const WORKER_ONLINE = 1;

    const WORKER_OFFLINE = 2;


    /**
     * 订单状态变更短信通知模板ID
     *
     * @var string
     */
    const ORDER_STATUS_CHANGE_SMS_CODE = '43951';

    /**
     * 待发货
     *
     * @var string
     */
    const ORDER_STATUS_WAIT_SEND = '40';

    /**
     * 订单状态描述
     * 订单状态：20:未支付;21:已取消;23:待支付;25:支付中;30:待确认(已支付);40:待发货;41:售后中(待发货,仅退款);60:待收货;61:售后中;80:已签收;81:售后中(已签收,退货退款);83:部分退款;100:已完成;200:退款中;201:退款失败;205:全部退款
     * @var array
     */
    const ORDER_STATUS_TEXT
        = [
            '0'   => '待支付',
            '20'  => '未支付',
            '21'  => '已取消',
            '23'  => '待支付',
            '25'  => '支付中',
            '30'  => '待确认',
            '40'  => '待发货',
            '41'  => '售后中',
            '60'  => '待收货',
            '80'  => '已签收',
            '81'  => '售后中',
            '83'  => '部分退款',
            '100' => '已完成',
            '200' => '退款中',
            '201' => '退款失败',
            '205' => '全部退款'
        ];

    /**
     * 订单配送类型
     *
     * @var array
     */
    const ORDER_DELIVERY_TYPE
        = [
            '1' => '需配送',
            '2' => '需自提'
        ];

    /**
     * @var string
     */
    const APP_CLIENT_NAME = 'store';

    /**
     * 账号格式
     *
     * @var string
     */
    const ACCOUNT_FORMAT = '+86-%s';

    /**
     * 自提事件
     *
     * @var string
     */
    const TAKE_GOODS_IN_STORE = 'take_goods_in_store';

    /**
     * 同意售后
     *
     * @var string
     */
    const REFUND_CONSOLE_AGREE = 'console_agree';

    /**
     * 拒绝申请售后
     *
     * @var string
     */
    const MANUAL_SEND_OUT_ORDER = 'manual_send_out_order';

    /**
     *  人工拒绝（驳回申请不符）
     *
     * @var string
     */
    const MANUAL_REFUSE_APPLY_REFUND = 'manual_refuse_apply_refund';

    /**
     * 售后同意退款
     *
     * @var string
     */
    const REFUND_MANUAL_AGREE_REFUND = 'manual_agree_refund';

    /**
     * 订单完成
     *
     * @var string
     */
    const ORDER_MANUAL_COMPLETE = 'manual_complete';

    /**
     * 确认发货
     *
     * @var string
     */
    const DELIVER_ORDER = 'store_send_out_order';

    /**
     * 需要通知的订单状态
     *
     * @var array
     */
    const SERVICE_ORDER_NOTICE_STATUS
        = [
            '40', '41', '81'
        ];

    /**
     * 推广订单需要通知的订单状态
     *
     * @var array
     */
    const PROMOTE_ORDER_NOTICE_STATUS = ['40'];

    /**
     * 推广订单退消息流转状态
     *
     * @var array
     */
    const PROMOTE_ORDER_STATUS_TRACE = [20,25];

    /**
     * 员工状态 1=在职，2=离职
     *
     * @var array
     */
    const MEMBER_WORK_STATUS
        = [
            1 => '在职',
            2 => '离职'
        ];

    /**
     * 导购类型-门店导购（店长或者店员）
     *
     * @var int
     */
    const MEMBER_GUIDE_TYPE_STORE = 1;

    /**
     * 导购类型-社群导购
     *
     * @var int
     */
    const MEMBER_GUIDE_TYPE_COMMUNITY = 2;

    /**
     * 导购类型-企业导购
     *
     * @var int
     */
    const MEMBER_GUIDE_TYPE_WORK = 3;

    /**
     * 店员
     *
     * @var int
     */
    const LOGIN_MANAGER_ROLE_STORE = 1;

    /**
     * 店长
     *
     * @var int
     */
    const LOGIN_MANAGER_ROLE_MANAGER = 2;

    /**
     * 企业导购
     *
     * @var int
     */
    const LOGIN_MANAGER_ROLE_WORK_MANAGER = 3;

    /**
     * 企业管理
     *
     * @var int
     */
    const LOGIN_MANAGER_ROLE_WORK_GUIDE = 4;

    /**
     * 短信模板
     *
     * @var array
     */
    const SMS_TEMPLATE
        = [
            'promote_store' => '746247',
            'promote_work'  => '746246',
            'service_store' => '746245'
        ];
    /**
     *
    746246
    【畅鲟智慧导购】您有一个新的推广订单“%orderNo%”，请进入畅鲟智慧导购小程序查看。
    746245
    【畅鲟智慧导购】您的门店“%storeName%”，有一个新的“%orderStatus%”服务订单“%orderNo%”等待处理，请尽快进入畅鲟智慧导购小程序操作，以免延误。
    746247
    【畅鲟智慧导购】您的门店“%storeName%”，有一个新的推广订单“%orderNo%”，请进入畅鲟智慧导购小程序查看。
     *
     */
    /**
     * 店员角色
     *
     * @var string
     */
    const LOGIN_ROLE_STORE_MEMBER = 'guide_store_member';

    /**
     * 店长角色
     *
     * @var string
     */
    const LOGIN_ROLE_STORE_MANAGER = 'guide_store_manager';

    /**
     * 企业角色
     *
     * @var string
     */
    const LOGIN_ROLE_ENTERPRISE = 'guide_enterprise_guide';

    /**
     * 社群导购角色
     *
     * @var string
     */
    const LOGIN_ROLE_COMMUNITY = 'guide_community_guide';

    /**
     * 企业管理角色
     *
     * @var string
     */
    const LOGIN_ROLE_ENTERPRISE_ADMIN = 'guide_enterprise_admin';

    /**
     * 合伙人
     *
     * @var string
     */
    const LOGIN_ROLE_STORE_PARTNER = 'guide_store_partner';

    /**
     * 角色列表
     *
     * @var string
     */
    const LOGIN_ROLE_LIST = [
        self::LOGIN_ROLE_STORE_MEMBER,
        self::LOGIN_ROLE_STORE_MANAGER,
        self::LOGIN_ROLE_ENTERPRISE,
        self::LOGIN_ROLE_COMMUNITY,
        self::LOGIN_ROLE_ENTERPRISE_ADMIN,
        self::LOGIN_ROLE_STORE_PARTNER
    ];

    /**
     * 订单类型
     * 1:推广订单 2:服务订单
     *
     * @var string
     */
    const EXTENSION_ORDER_TYPE = 1;
    const SERVICE_ORDER_TYPE = 2;


    /**
     * 订单类型
     * 1:待接订单 2:我的订单
     *
     * @var string
     */
    const RECEIVED_ORDER_TYPE = 1;
    const MY_ORDER_TYPE = 2;

    /**
     * 添加方式 1 导购名片 2 其他方式
     */
    const ADD_TYPE_BUSINESS_CARD = 1;
    const ADD_TYPE_BUSINESS_CARD_OTHER = 2;

    // 测试专用账号，无验证码限制
    const WX_LOGIN_SUPER_USER = [
        '15829712354',
        '18092066590',
        '15249279779',
        '15691849921'
    ];

    /**
     * 标准标识
     */
    const STANDARD_IDENTIFIER = 1;

    /**
     * 灰度标识
     */
    const RELEASE_IDENTIFIER = 2;

    /**
     * 标准(爱点击主体)标识
     */
    const STANDARD_ICLK_IDENTIFIER = 3;

}
