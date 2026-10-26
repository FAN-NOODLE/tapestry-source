<?php

namespace App\Http\Controllers\App\V1\Reserve;

use App\Http\Controllers\Controller;
use App\Micros\Common;
use App\Services\Wechat\WechatMiniLogin;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use Exception;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;
use Illuminate\Support\Facades\Log;

class WechatMiniLoginController extends Controller
{
    // 企业名称
    const STORE_NAME_KEY = 'store_name';

    // 企业配置的redis-key
    const REDIS_ENTERPRISE_CONFIG_KEY = "CRM:ENTERPRISE_CONFIG_KEY";

    public function index(Request $request)
    {
        $params = $request->all();

        if (empty($params['code'])) {//code的不能为空
            error(652012);
        }
        if (empty($params['enterprise_id'])) {
            error(652008);
        }

//        $name = $this->getConfigInfo($params['enterprise_id']);
//        if ($name == "coach") {
//            if (env("APP_ENV") == "prod") {
//                $appid  = "wx449e93ebf040f98c";
//                $secret = "REDACTED_FOR_SOURCE_SHARING";
//            } else {
//                $appid  = "wx3618ede4dbe080fa";
//                $secret = "REDACTED_FOR_SOURCE_SHARING";
//            }
//        } else {

        try {
            $pluginLists = (new EnterpriseWorkPlugin())->lists(['enterprise_id' => $params['enterprise_id']]);
        } catch (Exception $e) {
            error($e->getMessage());
        }
        if (empty($pluginLists)) {
            error(652012);
        }
        $appid  = "";
        $secret = "";
        foreach ($pluginLists as $item) {
            if ($item['identifier'] == "crs_store_assistant_self") {
                $appid  = $item['mini_app_id'];
                $secret = $item['mini_app_secret'];
            }
        }
//        }


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
        Common::microScrm()->setApi("micro-scrm/v1/customer/set-info")->post($tempToken);
        return success($data);
    }

    protected function getConfigInfo($enterpriseId)
    {
        $redisKey   = $enterpriseId . ":" . self::REDIS_ENTERPRISE_CONFIG_KEY;
        $configInfo = getOrSetData($redisKey);
        if (!$configInfo) {
            $configList                  = (new EnterpriseConfig())->info(['enterprise_id' => $enterpriseId, "key" => "store_name"]);
            $configInfo                  = current($configList['data']) ?? [];
            $configInfo['enterprise_id'] = $enterpriseId;
            getOrSetData($redisKey, $configInfo, 86400);
        }
        return $configInfo['name'];
    }


}