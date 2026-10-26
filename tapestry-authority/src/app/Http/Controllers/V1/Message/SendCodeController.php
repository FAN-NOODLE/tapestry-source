<?php


namespace App\Http\Controllers\V1\Message;


use App\Http\Controllers\V1\BaseController;
use App\Micros\Message\Message;

class SendCodeController extends BaseController
{
    public function index()
    {
        // 需要验证参数
        $mobile = $this->request->input('mobile', '');
        $client = $this->request->input('client', '');
        $event  = $this->request->input('event', '');
        $length = $this->request->input('len', '');

        $codeMin = str_pad(1, $length, 0);
        $codeMax = str_pad(9, $length, 9);

        $code = mt_rand($codeMin, $codeMax);

        $message   = new Message();
        $className = $message->getClassName("SendCode");

        if ($className) {
            $message->send([
                'mobile' => $mobile,
                'client' => $client,
                'event'  => $event,
                'code'   => $code
            ]);
        }

        return success();
    }
}
