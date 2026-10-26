<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 10:23 AM
 */

namespace App\Caches;

use Illuminate\Support\Facades\Redis;

class RedisCache
{
    const PREFIX = 'crs:shop:activity';

    const TTL = 3600;

    /**
     * @var \Illuminate\Redis\Connections\Connection
     */
    protected $redis;

    public function __construct()
    {
        $this->cacheOn = env('CACHE_ON', true);
        $this->redis   = Redis::connection();
    }

    public function key($params, $prefix)
    {
        $ext = '';

        if (!empty($params)) {

            $result = [];
            array_walk_recursive($params, function ($value, $key) use (&$result) {
                if (!is_numeric($key)) {
                    $value = $key . '_' . $value;
                }
                array_push($result, $value);
            });

            $ext .= implode(':', $result);

        }
        return rtrim(rtrim($prefix, ':') . ':' . $ext, ':');

    }

    public function setex($key, $value, $ttl = 3600)
    {
        $this->redis->setex($key, $ttl, $value);
    }

    public function setArray($key, $array, $ttl = 3600)
    {
        $this->redis->hmset($key, $array);
        $this->redis->expire($key, $ttl);
    }

    public function get($key)
    {
        $value = $this->redis->get($key);

        return is_numeric($value) ? (int)$value : $value;
    }

    public function getArray($key)
    {
        $array = $this->redis->hgetall($key);

        return $array;
    }

    public function decrBy($key, $num)
    {
        return $this->redis->decrby($key, $num);
    }

    public function incrBy($key, $num)
    {

        return $this->redis->incrby($key, $num);
    }

    public function del($keys)
    {
        if (is_string($keys)) {
            $keys = [$keys];
        }

        return $this->redis->del($keys);
    }

    public function exists($key)
    {
        return $this->redis->exists($key);
    }
}
