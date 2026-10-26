<?php


namespace App\Traits;

use App\Common\Constant;
use App\Library\Message\FactoryMessage;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

trait SendMessageTrait
{
    use RedisMutexLock;

    public function sendMessage($mobile, $content, $enterpriseId, $client = 'scrm', $template = '', $ip = '')
    {
        // 1分钟内同一手机不发送
        $lockKey = $this->mobileLockKey($mobile, $client);
        $lock    = $this->lock($lockKey, 60 * 1000);

        if (!$this->isActionAllowed($ip, 60 * 60, 500)) {
            notice(700037);
        }

        $limit = $lockBool = false;
        // 发送信息
        try {
            $factoryMessageInstance = new FactoryMessage($enterpriseId);
            $messageObject = $factoryMessageInstance->getMessageHandlers($factoryMessageInstance->getType());
            $account = $messageObject->getAccount();

            if(!in_array($mobile, array_keys(Constant::MESSAGE_MOBILE_WHITE_LIST))) {
                if (($this->getSMSSendCount($mobile, $account) ?? 0) >= Constant::MESSAGE_LIMIT_COUNT) {
                    $limit = true;
                }
            }

            if (!$lock) {
                $lockBool = true;
            }

            if(!$limit && !$lockBool){
                $config = $factoryMessageInstance->getConfig();
                $template = isset($config['template'][$template]) ? $config['template'][$template] : '';

                $messageData['subject'] = $mobile;
                $messageData['templateCode'] = $template;
                $messageData['data'] = $content;
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
            notice(700039);
        }

        if($limit){
            notice(700038);
        }

        if($lockBool){
            notice(700009);
        }

    }

    protected function getSMSSendKey($mobile, $account)
    {
        return 'micro:authority:SMSSend:count:account:message:' . $mobile . '-'. $account;
    }

    protected function mobileLockKey($mobile, $client)
    {
        return 'micro:authority:easy-sms:' . $client . ':mobile:lock:' . $mobile;
    }

    public function getSMSSendCount($mobile, $account)
    {
        $key = $this->getSMSSendKey($mobile, $account);
        return Redis::connection()->get($key);
    }
}
