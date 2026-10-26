<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/25
 * Time: 6:17 PM
 */

namespace App\Http\Controllers\Callback\V1\Customer;

use App\Http\Controllers\Callback\V1\BaseController;
use App\Library\Qy_callback\WXBizMsgCrypt;
use App\Library\Uc\Enterprise;
use Illuminate\Http\Request;


class OtherController extends BaseController
{
    /**
     * 多余的外部联系人回调
     *
     * @param Request $request
     * @return array|false
     * @throws \Exception
     */
    public function index(Request $request,$id)
    {
        $params = $request->all();
        if (empty($params['msg_signature']) ||empty($params['timestamp']) ||empty($params['nonce'])||empty($params['echostr'])) {
            return '';
        }
        $enterprise = (new \CrsSdk\Boss\Enterprise())->detail(['enterprise_id' => $id]);
        if (empty($enterprise['data'])) {
            return '';
        }
        $ucEnterprise = new Enterprise();
        $token = $ucEnterprise->entSettingGet($enterprise['data']['domain'],'wxqyToken')['wxqyToken'];
        $encodingAesKey = $ucEnterprise->entSettingGet($enterprise['data']['domain'],'wxqyEncodingAeskey')['wxqyEncodingAeskey'];
        $corpId = $enterprise['data']['corp_id'];
        if (empty($token) ||empty($encodingAesKey) ||empty($corpId)) {
            return '';
        }
        $wxcpt = new WXBizMsgCrypt($token, $encodingAesKey, $corpId);
        $sEchoStr = "";
        $errCode = $wxcpt->VerifyURL($params['msg_signature'], $params['timestamp'], $params['nonce'],$params['echostr'],$sEchoStr);
        if ($errCode == 0) {
             echo $sEchoStr;
        } else {
            print("ERR: " . $errCode . "\n\n");
        }
    }
}

