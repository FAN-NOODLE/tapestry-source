<?php

namespace App\Providers;

use App\Events\OrderServiceNotice;
use App\Events\StoreOrderStatus;
use App\Listeners\PromoteOrderNoticeListener;
use App\Listeners\ServiceOrderNoticeListener;
use Laravel\Lumen\Providers\EventServiceProvider as ServiceProvider;

class EventServiceProvider extends ServiceProvider
{
    public function boot()
    {
        parent::boot();
    }

    /**
     * The event listener mappings for the application.
     *
     * @var array
     */
    protected $listen = [
        // 订单通知处理
        StoreOrderStatus::class => [
            PromoteOrderNoticeListener::class,
        ],

        // 导购接单短信通知
        OrderServiceNotice::class => [
            ServiceOrderNoticeListener::class
        ]
    ];
}
