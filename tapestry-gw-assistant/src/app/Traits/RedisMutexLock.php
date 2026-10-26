<?php

namespace App\Traits;

use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use \Illuminate\Redis\Connections\Connection;

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
     * 获得锁,如果锁被占用,阻塞,直到获得锁或者超时。
     * -- 1、如果 $timeout 参数为 0,则立即返回锁。
     * -- 2、建议 timeout 设置为 0,避免 redis 因为阻塞导致性能下降。请根据实际需求进行设置。
     *
     * @param string $key 缓存KEY
     * @param int $timeout 取锁超时时间。单位(毫秒)。默认2000毫秒超时, 如果为0,如果当前锁被占用,则立即返回失败。如果大于0,则反复尝试获取锁直到达到该超时时间。
     * @param int $lockSecond 锁定时间。单位(毫秒)。默认锁1500毫秒
     * @param int $sleep 取锁间隔时间。单位(微秒)。默认100微秒 ,当锁为占用状态时, 每 100 微秒一次取锁。
     * @return bool 成功:true 失败:false
     */
    public function lock($key, $timeout = 2, $lockSecond = 1500, $sleep = 100)
    {
        if (strlen($key) === 0) {

            Log::error('分布式加锁异常', [
                'key'        => $key,
                'sleep'      => $sleep,
                'timeout'    => $timeout,
                'lockSecond' => $lockSecond,
            ]);
            error(100000);
        }

        $redis = self::getRedis();
        $start = self::getMicroTime();

        do {
            // [1] 锁的 KEY 不存在时设置其值, 并把过期时间设置为指定的时间。锁的值并不重要。重要的是利用 Redis 的特性。
            $acquired = $redis->set("{$key}:lock", 1, 'PX', $lockSecond, 'NX');

            if ($acquired) {
                break;
            }
            if ($timeout === 0) {
                break;
            }
            usleep($sleep);

            // 当前时间 小于 超时时间, 说明已经超时
        } while (!is_numeric($timeout) || (self::getMicroTime()) < ($start + ($timeout * 1000)));

        return $acquired ? true : false;
    }

    /**
     * 释放锁
     *
     * @param mixed $key 被加锁的KEY
     * @return int 大于0,说明删除成功
     */
    public function release($key)
    {
        if (strlen($key) === 0) {
            Log::error('分布式释放锁异常', [
                'key' => $key,
            ]);
            error(100000);
        }
        $redis = self::getRedis();
        return $redis->del(["{$key}:lock"]);
    }

    /**
     * 获取当前微秒。
     *
     * @return string
     */
    protected static function getMicroTime()
    {
        return bcmul(microtime(true), 1000000);
    }
}