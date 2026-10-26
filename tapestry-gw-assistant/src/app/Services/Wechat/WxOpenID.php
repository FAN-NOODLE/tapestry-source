<?php
/**
 * @author  zhaopengfei<zhaopengfei@vchangyi.com>
 * @date    2020/11/18 18:29
 */

namespace App\Services\Wechat;


use App\Exceptions\Error;
use App\Services\Decrypt;
use CrsSdk\Boss\ReleaseEnterpriseWork;
use think\Exception;

class WxOpenID
{
    protected $code;
    protected $wxAppID;
    protected $wxAppSecret;
    protected $wxLoginUrl = "https://api.weixin.qq.com/sns/jscode2session?"."appid=%s&secret=%s&js_code=%s&grant_type=authorization_code";

    function __construct($code,$identifier)
    {
        $this->code = $code;
        //更具env的状态返回对应的值
        [$appId,$secret] = Decrypt::getEnvMiniProgram($identifier);
        $this->wxAppID = $appId;
        $this->wxAppSecret = $secret;

        $this->wxLoginUrl = sprintf($this->wxLoginUrl, $this->wxAppID, $this->wxAppSecret, $this->code);
    }

    //依据code 从微信获得openid等信息
    public function get()
    {
        $result = curl_get_https($this->wxLoginUrl);
        $wxResult = json_decode($result, true);
        if (empty($wxResult)) {
            throw new Error(110002, '获取session_key 及 openID异常，微信内部错误');
        } else {
            $loginFail = array_key_exists('errcode', $wxResult);

            if ($loginFail) {
                $errmsg = $wxResult['errmsg'];
                $errcode = $wxResult['errcode'];
                throw new Error('110002',$errcode.'-'.$errmsg);
            } else {
                return $wxResult;
            }
        }
    }
}