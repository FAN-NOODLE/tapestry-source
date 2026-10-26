<?php


namespace App\Models;


class EnterpriseActivityMessage extends Base
{

    const PAGE_TYPE_GOODS = 1; // 商品
    const PAGE_TYPE_ACTIVITY_PAGE = 2; // 活动页
    const PAGE_TYPE_LIVE_LIST = 3; // 直播列表
    const PAGE_TYPE_LIVE_DETAIL = 4; // 直播详情
    const PAGE_TYPE_ACTIVITY = 5; // 活动
    const PAGE_TYPE_URL = 6; // 地址


    const SCENE_TYPE = 28;

    const OBJECT_TYPE = 16;

    const SEND_STATUS_NO = 0;
    const SEND_STATUS_YES = 1;

    const ENABLE_STATUS_ON = 0;
    const ENABLE_STATUS_OFF = 1;
    // 页面地址
    const PAGE_URL = [
        self::PAGE_TYPE_GOODS         => [
            'url'     => '/goods/goods-detail/goods-detail',
            'query'   => '?id=%spu_id&rule_id=%spu_id&sku_id=%sku_id',
            'replace' => ['spu_id', 'sku_id'],
        ],
        self::PAGE_TYPE_ACTIVITY_PAGE => [
            'url'     => '/other/activity/activity',
            'query'   => '?activity_id=%id',
            'replace' => ['id']
        ],
        self::PAGE_TYPE_LIVE_LIST     => [
            'url'     => '/live/live-list/live-list',
            'query'   => '',
            'replace' => []
        ],
        self::PAGE_TYPE_LIVE_DETAIL   => [
            'url'     => '/live/live-detail/live-detail',
            'query'   => '?live_id=%id',
            'replace' => ['id']
        ],
        self::PAGE_TYPE_ACTIVITY      => [
            'default'   => [ // 默认
                'url'     => '/goods/goods-detail/goods-detail',
                'query'   => '?id=%spu_ids&rule_id=%id&type=%activity_type',
                'replace' => ['id', 'spu_ids', 'activity_type']
            ],
            'one-price' => [ // 一口价
                'url'     => '/one-price/one-price/one-price',
                'query'   => '?id=%id',
                'replace' => ['id']
            ]
        ],
        self::PAGE_TYPE_URL           => [
            'url'     => '',
            'query'   => '',
            'replace' => []
        ],
    ];
    protected $table = 'enterprise_activity_message';

    protected $guarded = [];
}