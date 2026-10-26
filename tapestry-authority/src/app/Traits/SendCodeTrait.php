<?php


namespace App\Traits;

use Exception;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Overtrue\EasySms\EasySms;

trait SendCodeTrait
{
    use RedisMutexLock;

    /**
     * @param string $mobile
     * @param string $client
     * @param string $event
     * @param int $code
     * @param int $ttl
     */
    public function sendCode($mobile, $client = '', $event = '', $code = 1234, $ttl= 900)
    {
        // 1分钟内同一手机不发送
        $lockKey = $this->mobileLockKey($mobile, $client, $event);
        $lock    = $this->lock($lockKey, 60 * 1000);

        if (!$lock) {
            notice(700009);
        }

        // 发送 code
        try {
            $config  = config('message')['sendCode'];
            $template = $config[$client][$event];
            if (empty($template)) {
                notice(700008);
            }
            $easySms = new EasySms($config);
            $easySms->send($mobile, [
                'template' => $template,
                'data'     => [
                    'code' => $code
                ],
            ]);
        } catch (Exception $e) {
            $this->release($lockKey);
            Log::info('发送验证码失败' . $mobile);
            notice(700008);
        }

        // 缓存 code
        $this->setCodeCache($mobile, $code, $client, $event, $ttl);
    }

    /**
     * @param $mobile
     * @param $code
     * @param $client
     * @param $event
     * @param $ttl
     */
    public function setCodeCache($mobile, $code, $client, $event, $ttl)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client, $event);
        // 验证码存入缓存，缓存15分钟
        Redis::connection()->setex($key, $ttl, $code);
    }

    /**
     * @param $mobile
     * @param $client
     * @param $event
     * @return string
     */
    public function getCodeCache($mobile, $client, $event)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client, $event);
        return Redis::connection()->get($key);
    }

    /**
     * @param $mobile
     * @param $client
     * @param $event
     * @return int
     */
    public function delCodeCache($mobile, $client, $event)
    {
        $key = $this->mobileCodeCacheKey($mobile, $client, $event);
        return Redis::connection()->del([$key]);
    }

    /**
     * @param $mobile
     * @param $client
     * @param $event
     * @return string
     */
    protected function mobileCodeCacheKey($mobile, $client, $event)
    {
        return 'micro:authority:easy-sms:' . $client . ':' . $event .':mobile:' . $mobile;
    }

    /**
     * @param $mobile
     * @param $client
     * @param $event
     * @return string
     */
    protected function mobileLockKey($mobile, $client, $event)
    {
        return 'micro:authority:easy-sms:' . $client . ':' .$event. ':mobile:lock:' . $mobile;
    }
}
