<?php

namespace App\Http\Controllers\V1\Common;

use App\Http\Controllers\V1\BaseController;
use App\Services\SendCodeService;
use App\Traits\RedisMutexLock;
use Illuminate\Support\Facades\Log;
use Overtrue\EasySms\EasySms;

class SmsController extends BaseController
{
    use RedisMutexLock;

    /**
     * 发送短信验证码
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function sendCode()
    {
        $this->validate($this->request, [
            'mobile' => 'required|string',
            'event'  => 'required|string',
            'rule'   => 'string',
            'length' => 'integer|min:4|max:6'
        ]);
        $input = $this->request->all();

        if (isset($input['rule']) && $input['rule'] !== '') {
            (new SendCodeService())->sendCodeRule($input['rule'], $input['mobile']);
        }

        // 验证码长度 默认4位
        if (!isset($input['length'])) {
            $input['length'] = 4;
        }

        // 读取该事件下对应的短信模版
        $template = config('easy-sms.event_template_relation')[$input['event']] ?? null;
        if ($template === null) {
            notice(100000, ['事件不存在']);
        }

        // 1分钟内同一手机不发送
        $lockKey = "micro:authority:sms:{$input['event']}:mobile:lock:{$input['mobile']}";
        $lock = $this->lock($lockKey, 60 * 1000);
        if (!$lock) {
            notice(700009);
        }

        // 生成验证码
        $code = '';
        for ($i = 0; $i < $input['length']; $i++) {
            $code .= mt_rand(0, 9);
        }

        // 发送 code
        try {
            $config = config('easy-sms');
            $easySms = new EasySms($config);
            $easySms->send($input['mobile'], [
                'template' => $template,
                'data'     => [
                    'code' => $code
                ],
            ]);
        } catch (\Exception $exception) {
            $this->release($lockKey);

            Log::error('发送验证码失败', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'input'   => $input
            ]);

            error(700008);
        }

        // 验证码存入缓存，缓存15分钟
        $key = $this->cacheKey($input['event'], $input['mobile']);
        self::getRedis()->setex($key, 900, $code);

        return success();
    }

    /**
     * 验证短信验证码
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function checkCode()
    {
        $this->validate($this->request, [
            'mobile' => 'required|string',
            'event'  => 'required|string',
            'code'   => 'required|string'
        ]);
        $input = $this->request->all();

        $key = $this->cacheKey($input['event'], $input['mobile']);

        // 缓存中获取验证码
        $cacheCode = self::getRedis()->get($key);
        if ($cacheCode != $input['code']) {
            notice(700007);
        }

        // 删除验证码
        // self::getRedis()->del([$key]);

        return success();
    }

    private function cacheKey($event, $mobile)
    {
        return "micro:authority:sms:{$event}:mobile:{$mobile}";
    }
}
