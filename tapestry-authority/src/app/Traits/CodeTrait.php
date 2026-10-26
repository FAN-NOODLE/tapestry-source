<?php


namespace App\Traits;

use App\Library\Message\FactoryMessage;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use Overtrue\EasySms\EasySms;
use App\Common\Constant;

trait CodeTrait
{
    private $failTimes = 5; //错误次数

    private $lockTime = 3600 * 24; //锁定时间

    use RedisMutexLock;

    public function sendCode($mobile, $client = 'console', $fixed = 2, $fixedCode = 1234, $ip = '', $template = '', $enterpriseId = '')
    {
        // 删除验证码过期已发送的验证码
        $this->delCodeCache($mobile, $client);

        // 1分钟内同一手机不发送
        $lockKey = $this->mobileLockKey($mobile, $client);
        $lock    = $this->lock($lockKey, 60 * 1000);

        if (!$this->isActionAllowed($ip, 60 * 60, 500)) {
            notice(700030);
        }
        $code = mt_rand(100000, 999999);
        if ($fixed === 1) {
            $code = $fixedCode;
        }

        // console、boss公用一个模板
        $tmpClient = $client;
        if ($client == "boss") {
            $tmpClient = "console";
        }

        $limit = $lockBool = false;
        // 发送 code
        try {
            $factoryMessageInstance = new FactoryMessage($enterpriseId);
            $messageObject = $factoryMessageInstance->getMessageHandlers();

            $account = $messageObject->getAccount();

            if(!in_array($mobile, array_keys(Constant::MESSAGE_MOBILE_WHITE_LIST))){
                if (($this->getSMSSendCount($mobile, $account) ?? 0) >= Constant::MESSAGE_LIMIT_COUNT) {
                    $limit = true;
                }
            }

            if (!$lock) {
                $lockBool = true;
            }

            if(!$limit && !$lockBool) {
                $config       = $factoryMessageInstance->getConfig();
                $templateCode = isset($config['template'][$tmpClient]) ? $config['template'][$tmpClient] : '';

                $messageData['subject']      = $mobile;
                $messageData['templateCode'] = $templateCode;
                $messageData['data']         = [
                    'code' => $code
                ];

                $messageObject->send($messageData);
                //统计发送次数
                Redis::connection()->incr($this->getSMSSendKey($mobile, $account));
                Redis::connection()->expire($this->getSMSSendKey($mobile, $account), (strtotime('23:59:59') - time()));
            }
        } catch (\Exception $e) {
            Log::error("短信验证码发送失败：", [$e->getMessage()]);
            $this->release($lockKey);
            Log::info('发送验证码失败' . $mobile, [
                'msg'  => $e->getMessage(),
                'line' => $e->getLine(),
                'file' => $e->getFile(),
                'code' => $e->getCode(),
            ]);
            notice(700008);
        }

        if($limit){
            notice(700033);
        }

        if($lockBool){
            notice(700009);
        }

        // 缓存 code
        $this->setCodeCache($mobile, $code, $client);
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

    protected function lockAccountCacheKey($account, $client)
    {
        return 'micro:authority:lock-login:' . $client . ':account:lock:' . $account;
    }

    protected function getSMSSendKey($mobile, $account)
    {
        return 'micro:authority:SMSSend:count:account:code:' . $mobile .'-'. $account;
    }

    protected function mobileLockKey($mobile, $client)
    {
        return 'micro:authority:easy-sms:' . $client . ':mobile:lock:' . $mobile;
    }

    public function getLockAccountCacheKey($account, $client)
    {
        $key = $this->lockAccountCacheKey($account, $client);
        return Redis::connection()->get($key);
    }

    public function getSMSSendCount($mobile, $account)
    {
        $key = $this->getSMSSendKey($mobile, $account);
        return Redis::connection()->get($key);
    }
}
