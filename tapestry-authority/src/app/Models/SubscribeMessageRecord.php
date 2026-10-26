<?php

namespace App\Models;

class SubscribeMessageRecord extends Base
{
    /**
     * 允许一次
     */
    const TYPE_ONCE = 1;

    /**
     * 总是允许
     */
    const TYPE_ALSO = 2;

    /**
     * 用户订阅状态
     * 0未订阅（默认）
     * 1订阅成功
     * 2已订阅
     */
    const SUBSCRIBE_FALSE = 0;
    const SUBSCRIBE_TRUE = 1;
    const SUBSCRIBED = 2;

    /**
     * 订阅对象：1 视频直播
     * 对应的object_id为直播id
     */
    const OBJECT_TYPE_VIDEO = 1;

    /**
     * 订阅对象：2 优惠券
     * 对应的object_id为优惠券id
     */
    const OBJECT_TYPE_COUPON = 2;

    /**
     * 订阅对象：3 订单
     * 对应的object_id为订单id
     */
    const OBJECT_TYPE_ORDER = 3;

    /**
     * 订阅对象：4 售后订单
     * 对应的object_id为订单id
     */
    const OBJECT_TYPE_REFUND_ORDER = 4;

    /**
     * 订阅对象：5 发票
     * 对应的object_id为票据id
     */
    const OBJECT_TYPE_TICKET = 5;

    /**
     * 订阅对象：6 限时购
     * 对应的object_id为商品spu_id
     */
    const OBJECT_TYPE_FLASH_SALE = 6;

    /**
     * 订阅对象：7 助力抽奖
     * 对应的object_id为活动发起id
     */
    const OBJECT_TYPE_LUCK = 7;

    /**
     * 订阅对象：8 拼团
     * 对应的object_id为活动发起id
     */
    const OBJECT_TYPE_GROUP = 8;

    /**
     * 订阅对象：9 砍价
     * 对应的object_id为活动发起id
     */
    const OBJECT_TYPE_BARGAIN = 9;

    /**
     * 订阅对象：10 翻倍券
     * 对应的object_id为活动发起id
     */
    const OBJECT_TYPE_DOUBLE_COUPON = 10;

    /**
     * 订阅对象：11 好友助力
     * 对应的object_id为活动发起id
     */
    const OBJECT_TYPE_HELP = 11;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'subscribe_message_record';

    /**
     * 黑名单
     *
     * @var array
     */
    protected $guarded = [];

    public function template()
    {
        return $this->belongsTo(SubscribeMessageTouchTemplate::class, 'pri_tmpl_id', 'pri_tmpl_id');
    }
}
