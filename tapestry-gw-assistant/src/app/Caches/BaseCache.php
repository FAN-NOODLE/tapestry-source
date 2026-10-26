<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 10:23 AM
 */

namespace App\Caches;

class BaseCache
{
    protected $cacheOn = true;

    const PREFIX = 'crs:shop:activity';

    const TTL = 3600;

    /**
     * App\Caches\Cache
     */
    protected $driver;

    public function __construct(Cache $cacheDriver = null)
    {
        // 缓存开关
        $this->cacheOn = env('ACTIVITY_CACHE_ON', true);

        if (empty($cacheDriver)) {
            $this->driver = new RedisCache();
        } else {
            $this->driver = $cacheDriver;
        }
    }

    final public function key($params, $prefix)
    {
        return $this->driver->key($params, $prefix);
    }

    final public function set($key, $value, $ttl = 3600)
    {
        $this->driver->set($key, $value, $ttl);
    }

    final public function get($key)
    {
        return $this->driver->get($key);
    }

    final public function setArray($key, $array, $ttl = 3600)
    {
        $this->driver->setArray($key, $array, $ttl);
    }

    final public function getArray($key)
    {
        return $this->driver->getArray($key);
    }

    final public function updateField($key, $field, $value)
    {
        $this->driver->updateField($key, $field, $value);
    }

    final public function delete($keys)
    {
        $this->driver->delete($keys);
    }
}
