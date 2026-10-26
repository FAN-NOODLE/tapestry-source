<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/24 11:24
 */

namespace App\Services;


use Illuminate\Support\Facades\Log;
use Overtrue\EasySms\EasySms;

class SendSmsService extends BaseService
{
    public function sendCode($mobile, $data,$template = 'balance_status')
    {
        if (empty($mobile)){
            Log::info('发送短信收假后:'.json_encode([
                    'mobile' => $mobile,
                    'msg'            =>'手机号不能为空',
                ]));
        }
        switch ($template){
            case 'distribution_order_status':
                 $txt = '结算单变更';
                 break;
            case 'balance_status':
                 $txt = '余额变更';
                 break;
            default:
                $txt = '';
        }
        $template = config('easy-sms.template')[$template];
        // 发送 code
        try {
            $config  = config('easy-sms');
            $easySms = new EasySms($config);
            $easySms->send($mobile, [
                'template' => $template,
                'data'     => $data,
            ]);
        } catch (\Exception $e) {
            Log::error("{$mobile}{$txt}发送失败：", [$e->getExceptions()]);
            notice(1100023);
        }
        return true;
    }
}