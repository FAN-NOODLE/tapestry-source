<?php


namespace App\Micros\Message;

use App\Traits\SendCodeTrait;

class SendCode extends AbstractMessage
{
    use SendCodeTrait;

    public function send($params)
    {
        $mobile = $params['mobile'];
        $client = $params['client'];
        $event  = $params['event'];
        $code   = $params['code'];

        $this->sendCode($mobile, $client, $event, $code);
    }
}
