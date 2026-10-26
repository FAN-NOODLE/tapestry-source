<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/15
 * Time: 4:05 PM
 */

namespace App\Models;

class EnterpriseApp extends Base
{
    /**
     * 表名称
     *
     * @var string
     */
    protected $table = 'enterprise_app';

    public $app = [
        'coupon' => [
            'title'   => '礼劵（默认开启）',
            'content' => '劵动客心，营销必备'
        ],

        'novice' => [
            'title'   => '新人优惠',
            'content' => '新客专享，提升转化'
        ],

        'sec_kill' => [
            'title'   => '限时购',
            'content' => '商品限时打折促销'
        ],

        'group' => [
            'title'   => '拼团',
            'content' => '拼单成团，裂变式转化'
        ],

        'cheap' => [
            'title'   => '超低购',
            'content' => '超低价格，让你购'
        ],

        'activity' => [
            'title'   => '优惠',
            'content' => '满折、满减、满赠'
        ],
        
        'collection' => [
            'title'   => '收藏有礼',
            'content' => '收藏小程序，赠送优惠券'
        ]
    ];
}
