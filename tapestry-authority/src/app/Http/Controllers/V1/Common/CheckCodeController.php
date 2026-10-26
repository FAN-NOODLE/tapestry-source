<?php


namespace App\Http\Controllers\V1\Common;


use App\Http\Controllers\V1\BaseController;
use App\Library\Util\Logger;
use App\Traits\CodeTrait;

class CheckCodeController extends BaseController
{
    use CodeTrait;

    public function index()
    {
        $code   = $this->request->input('code', '');
        $mobile = $this->request->input('mobile', '');
        $client = $this->request->input('client', '');

        if (empty($code)) {
            notice(700007);
        }
        if (empty($mobile)) {
            notice(700010);
        }
        if (empty($client)) {
            notice(700013);
        }

        // 缓存中获取验证码
        $cacheCode = $this->getCodeCache($mobile, $client);
        if ($cacheCode != $code) {
            notice(700007);
        }

        // 删除验证码
       
        $this->delCodeCache($mobile, $client);

        return success();
    }
}
