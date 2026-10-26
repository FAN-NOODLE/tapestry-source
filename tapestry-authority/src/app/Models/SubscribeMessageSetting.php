<?php

namespace App\Models;

class SubscribeMessageSetting extends Base
{
    // 模板类型：一次性订阅
    const TEMPLATE_TYPE_ONCE = 2;
    // 模板类型：长期订阅
    const TEMPLATE_TYPE_ALSO = 3;

    // 场景类型：直播订阅
    const SCENE_TYPE_LIVE = 1;

    // 场景类型：优惠券状态通知
    const SCENE_TYPE_COUPON_STATUS = 2;
    // 场景类型：订单待支付提醒
    const SCENE_TYPE_WAIT_PAY = 3;
    // 场景类型：订单发货提醒
    const SCENE_TYPE_WAIT_DELIVER = 4;
    // 场景类型：售后审核结果通知
    const SCENE_TYPE_REFUND_AUDIT_RESULT = 5;
    // 场景类型：发票
    const SCENE_TYPE_INVOICE = 6;

    // 场景类型：限时购预约
    const SCENE_TYPE_FLASH_SALE_SUBSCRIBE = 7;

    // 场景类型：助力抽奖预约
    const SCENE_TYPE_LUCK_SUBSCRIBE = 8;
    // 场景类型：助力抽奖活动进度提醒
    const SCENE_TYPE_LUCK_PLAN = 9;
    // 场景类型：助力抽奖开奖结果通知
    const SCENE_TYPE_LUCK_RESULT = 10;

    // 场景类型：

    // 场景类型：拼团进度提醒
    const SCENE_TYPE_GROUP_PLAN = 12;
    // 场景类型：拼团结果通知
    const SCENE_TYPE_GROUP_RESULT = 13;

    // 场景类型：砍价进度提醒
    const SCENE_TYPE_BARGAIN_PLAN = 14;
    // 场景类型：砍价失败提醒
    const SCENE_TYPE_BARGAIN_FAIL = 15;
    // 场景类型：砍价成功提醒
    const SCENE_TYPE_BARGAIN_SUCCESS = 16;

    // 场景类型：翻倍券失败提醒
    const SCENE_TYPE_DOUBLE_COUPON_FAIL = 17;
    // 场景类型：翻倍券成功提醒
    const SCENE_TYPE_DOUBLE_COUPON_SUCCESS = 18;

    // 场景类型：好友助力失败提醒
    const SCENE_TYPE_HELP_FAIL = 19;
    // 场景类型：好友助力完成提醒
    const SCENE_TYPE_HELP_COMPLETED = 20;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'subscribe_message_setting';

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];

    public function belongsTid()
    {
        return $this->belongsTo(SubscribeMessageTid::class, 'tid', 'tid');
    }

    public function category()
    {
        return $this->belongsTo(EnterpriseCategory::class, 'category_id', 'id');
    }
}
