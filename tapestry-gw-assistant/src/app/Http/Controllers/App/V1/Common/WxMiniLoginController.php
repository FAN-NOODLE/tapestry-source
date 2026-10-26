<?php

namespace App\Http\Controllers\App\V1\Common;

use App\Http\Controllers\Controller;
use App\Services\Wechat\WechatMiniLogin;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class WxMiniLoginController extends Controller
{
    public function index(Request $request)
    {
        $params = $request->all();

        if (empty($params['code'])) {//code的不能为空
            error(652012);
        }
        if (empty($params['enterprise_id'])) {
            error(652008);
        }

        try {
            $pluginLists = (new EnterpriseWorkPlugin())->lists(['enterprise_id' => $params['enterprise_id']]);
        } catch (Exception $e) {
            error($e->getMessage());
        }
        if (empty($pluginLists)) {
            error(652012);
        }
        foreach ($pluginLists as $item) {
            if ($item['identifier'] == "crs_store_assistant_self") {
                $appid  = $item['mini_app_id'];
                $secret = $item['mini_app_secret'];
            }
        }


        //获取微信openid
        $result = (new WechatMiniLogin($params['code'], $appid, $secret))->get();

        $openId     = $result['openid'];
        $sessionKey = $result['session_key'];

        if (!is_array($result)) {
            Log::error($result);
            error(657005);
        }

        // 授权返回结果错误
        if (empty($result['openid']) || empty($result['session_key'])) {
            Log::error('微信授权数据错误:' . var_export($result, true));
            error(657006);
        }

        $tempToken     = [
            'expire_time'   => time() + 86400,
            'enterprise_id' => $params['enterprise_id'],
            'open_id'       => $result['openid'],
            'union_id'      => $result['unionid']
        ];
        $data['token'] = Crypt::encrypt(json_encode($tempToken));
        return success($data);
    }
}