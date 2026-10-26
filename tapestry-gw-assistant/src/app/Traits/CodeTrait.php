<?php


namespace App\Traits;

use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Overtrue\EasySms\EasySms;

trait CodeTrait
{
    use RedisMutexLock;

    public function sendCode($mobile, $client = 'store', $fixed = 2, $fixedCode = 1234,$tll = 900,$template = 'store' )
    {
        // 1分钟内同一手机不发送
        $lockKey = $this->mobileLockKey($mobile, $client);
        $lock    = $this->lock($lockKey, 60 * 1000);

        if (!$lock) {
            notice(700009);
        }

        $code = mt_rand(1000, 9999);
        if ($fixed === 1) {
            $code = $fixedCode;
        }

       $template = config('easy-sms.template')[$template];

        // 发送 code
        try {
            $config  = config('easy-sms');
            $easySms = new EasySms($config);
            $easySms->send($mobile, [
                'template' => $template,
                'data'     => [
                    'code' => $code
                ],
            ]);
        } catch (\Exception $e) {
            Log::error("{$mobile}短信验证码发送失败：", [$e->getExceptions()]);
            $this->release($lockKey);
            notice(700008);
        }
        // 缓存 code
        $this->setCodeCache($mobile, $code, $client,$tll);
    }

    public function setCodeCache($mobile, $code, $client, $ttl = 900)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client);
        // 验证码存入缓存，缓存15分钟
        Redis::connection()->setex($key, $ttl, $code);
    }

    public function getCodeCache($mobile, $client)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client);
        return Redis::connection()->get($key);
    }

    public function delCodeCache($mobile, $client)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client);
        return Redis::connection()->del([$key]);
    }

    protected function mobileCodeCacheKey($mobile, $client)
    {
        return 'micro:authority:easy-sms:' . $client . ':mobile:' . $mobile;
    }

    protected function mobileLockKey($mobile, $client)
    {
        return 'micro:authority:easy-sms:' . $client . ':mobile:lock:' . $mobile;
    }
}
