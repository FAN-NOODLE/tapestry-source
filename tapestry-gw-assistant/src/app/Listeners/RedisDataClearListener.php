<?php

namespace App\Listeners;

use App\Events\OrderRedisDataEvent;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class RedisDataClearListener
{
    /**
     * Create the event listener.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Handle the event
     * @param OrderRedisDataEvent $event
     */
    public function handle(OrderRedisDataEvent $event)
    {
        // 订单信息更改，清除已存数据
        $orderKey = Redis::keys('order:list:' . $event->enterpriseId . ':' . $event->userId . '*');
        // 记录下删除的order记录
        Log::debug('redis_order:' . $event->enterpriseId . ':' . $event->userId . ':' . json_encode($orderKey));
        foreach ($orderKey as $item) {
            Redis::del($item);
        }
    }
}
