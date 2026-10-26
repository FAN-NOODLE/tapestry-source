<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/19 10:04
 */

namespace App\Services\Wechat;


use App\Exceptions\Error;
use App\Services\Decrypt;
use App\Services\Wechat\Decrypt\WXBizDataCrypt;
use Illuminate\Support\Facades\Log;

class WxDecrypt
{

    protected $wxAppID;
    protected $sessionKey;

    function __construct($sessionKey, $identifier)
    {
        $this->sessionKey = $sessionKey;
        //更具env的状态返回对应的值
        [$appId,$secret] = Decrypt::getEnvMiniProgram($identifier);
        $this->wxAppID = $appId;
    }

    public function index($input)
    {
        $pc = new WXBizDataCrypt($this->wxAppID, $this->sessionKey);
        $data = $pc->decryptData($input['encryptedData'], $input['iv']);

        if ($data['errorCode'] == 0) {
            return $data['result'];
        } else {
            Log::info('手机号授权失败:' . $data['errorCode'] . '|appId:' . $this->wxAppID . '|sessionKey:' . $this->sessionKey);
            throw new Error(110002, '用户授权获得手机号失败-'.$data['errorCode']);
        }
    }
}