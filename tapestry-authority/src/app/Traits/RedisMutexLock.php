<?php

namespace App\Traits;

use Illuminate\Support\Facades\Redis;
use Illuminate\Redis\Connections\Connection;

trait RedisMutexLock
{
    /**
     * 获取 Redis 连接
     *
     * @return Connection
     */
    public static function getRedis()
    {
        return Redis::connection();
    }

    /**
     * 获得锁,如果锁被占用,阻塞,直到获得锁或者超时
     * -- 1、如果 $timeout 参数为 0,则立即返回锁
     * -- 2、建议 timeout 设置为 0,避免 redis 因为阻塞导致性能下降。请根据实际需求进行设置
     *
     * @param string $key 缓存KEY
     * @param int $lockSecond 锁定时间。单位(毫秒)。默认锁5000毫秒
     * @param int $timeout 取锁超时时间。单位(毫秒)。默认0毫秒超时, 如果为0,如果当前锁被占用,则立即返回失败。如果大于0,则反复尝试获取锁直到达到该超时时间
     * @return bool 成功:true 失败:false
     */
    public function lock(string $key, $lockSecond = 5000, $timeout = 200): bool
    {
        $redis = self::getRedis();
        $start = self::millisecondTime();

        do {
            // [1] 锁的 KEY 不存在时设置其值, 并把过期时间设置为指定的时间。锁的值并不重要。重要的是利用 Redis 的特性。
            $acquired = $redis->set("{$key}:lock", LOG_ID, 'PX', $lockSecond, 'NX');

            if ($acquired) {
                break;
            }
            if ($timeout === 0) {
                break;
            }

            // 2ms重试一次
            usleep(2000);

            // 当前时间 小于 锁过期时间, 说明已经超时
        } while (self::millisecondTime() < ($start + $timeout));

        return $acquired ? true : false;
    }

    /**
     * @param string $key 被加锁的KEY
     * @return int
     */
    public function release($key)
    {
        $redis = self::getRedis();

        return $redis->del(["{$key}:lock"]);
    }

    public function isActionAllowed($ip, $period, $maxCount)
    {
        $key   = 'limit:' . $ip;
        $nowTs = (int)(microtime(true) * 1000);
        $redis = Redis::connection();
        //$redis->multi();
        $redis->zadd($key, [$nowTs => $nowTs]);
        $redis->zremrangebyscore($key, 0, $nowTs - $period * 1000);
        $count = $redis->zcard($key);
        $redis->expire($key, $period + 1);
        //$redis->exec();
        return $count <= $maxCount;
    }

    /**
     * 当前毫秒时间
     *
     * @return false|float
     */
    protected static function millisecondTime()
    {
        return round(microtime(true) * 1000, 2);
    }
}