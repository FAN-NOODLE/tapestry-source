<?php

namespace App\Traits;

use Illuminate\Redis\Connections\Connection;
use Illuminate\Support\Facades\Redis;

trait ImportTrait
{
    /**
     * @var Connection
     */
    public $redis;
    public $successType = 1;
    public $errorType = 2;
    public $ttl = 7200;

    public function getRedis()
    {
        if (!($this->redis instanceof Connection)) {
            $this->redis = Redis::connection();
        }

        return $this->redis;
    }

    public function getKey($flag = '')
    {
        $key = md5(uniqid());

        if (!empty($key)) {
            $key .= '_' . $flag;
        }

        return $key;
    }

    public function pushData(array $row, $key, $type)
    {
        $dataKey = $this->getDataKey($key, $type);

        $redis = $this->getRedis();
        $redis->lpush($dataKey, [json_encode($row)]);
        $this->redis->expire($dataKey, $this->ttl);
    }

    public function getDataTotal($key, $type)
    {
        $dataKey = $this->getDataKey($key, $type);
        $redis   = $this->getRedis();
        return $redis->llen($dataKey);
    }

    public function getDataList($key, $type = 1, $page = 0, $size = 0)
    {
        $start   = max(0, ($page - 1) * $size);
        $stop    = $start + $size - 1;
        $dataKey = $this->getDataKey($key, $type);
        $redis   = $this->getRedis();
        $list    = $redis->lrange($dataKey, $start, $stop);

        $formatList = [];

        foreach ($list as $k => $v) {
            $tmp          = json_decode($v, true);
            $tmp['index'] = $start + $k;

            $formatList[] = $tmp;
        }

        return $formatList;
    }

    public function updateDataByIndex($key, $type, $index, $row)
    {
        $redis   = $this->getRedis();
        $dataKey = $this->getDataKey($key, $type);

        $redis->lset($dataKey, $index, json_encode($row));
    }

    protected function getDataKey($key, $type)
    {
        return $key . ':' . $type;
    }
}
