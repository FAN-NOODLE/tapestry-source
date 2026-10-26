<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2020-03-19
 * Time: 16:56
 */

namespace App\Micros\Authority;

use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class AccountMicro
{
    protected $params;
    protected $response;

    public function __construct($params)
    {
        $this->params = $params;
        Log::info('micro-authority-v1-request', $params);
    }

    public function decryptToken()
    {
        $token = $this->params['token'];
        $key   = 'gw:console:token:' . $token;
        $redis = Redis::connection();
        if ($redis->exists($key)) {
            return json_decode($redis->get($key), true);
        }

        $tokenArr = (new Account())->decryptToken($this->params);

        $ttl = min(1, $tokenArr['expire_time']);
        $redis->setex($key, $ttl, json_encode($tokenArr));

        return $tokenArr;
    }

    public function dataRange()
    {
        $adminId = $this->params['admin_id'];
        $key     = 'gw:console:data-range:' . $adminId;

        $redis = Redis::connection();
        if ($redis->exists($key)) {
            return json_decode($redis->get($key), true);
        }

        $dataRangeList = (new Account())->dataRange($this->params);

        // 缓存3分钟
//        $ttl = 180;
//        $redis->setex($key, $ttl, json_encode($dataRangeList));

        return $dataRangeList;
    }

    public function __destruct()
    {
        Log::info('micro-authority-v1-response', $this->response ? (array)$this->response : []);
    }
}
