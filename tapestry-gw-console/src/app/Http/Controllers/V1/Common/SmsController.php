<?php


namespace App\Http\Controllers\V1\Common;

use CrsSdk\MicroAuthority\Common;
use App\Http\Controllers\V1\BaseController;

class SmsController extends BaseController
{
    protected $requireLogin = false;

    /**
     * 发送短信验证码
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function sendCode()
    {
        $input = $this->request->all();

        $data = (new Common())->sendSmsCode($input);

        return success($data);
    }
}
