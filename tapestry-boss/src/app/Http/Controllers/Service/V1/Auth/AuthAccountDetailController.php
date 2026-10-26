<?php

namespace App\Http\Controllers\Service\V1\Auth;

use App\Common\RedisKey;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Json;
use App\Models\Enterprise;
use App\Models\EnterpriseOpenPlatform;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class AuthAccountDetailController extends BaseController
{
    /**
     * 获取小程序授权详情数据
     *
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $enterpriseId   = $request->input('enterprise_id', 0);
        $enterpriseHash = $request->input('enterprise_hash', '');
        $auth_type      = $request->input('auth_type', 0);

        if (empty($enterpriseId) && empty($enterpriseHash)) {
            error(600001);
        }

        if (!in_array($auth_type, [EnterpriseThirdPartAccount::TYPE_SHOP, EnterpriseThirdPartAccount::TYPE_QY_WX])) {
            error(600003);
        }

        if (!empty($enterpriseId)) {
            $enterprise = Enterprise::query()
                ->where('id', '=', $enterpriseId)
                ->first();
        } else {
            $enterprise = Enterprise::query()
                ->where('hash', '=', $enterpriseHash)
                ->first();
        }


        $redis = redis();
        $key   = sprintf(RedisKey::AUTH_APP_INFO, $enterprise->id);

        $data = $redis->get($key);
        if ($data === null) {

            // 查询授权详情
            $authInfo = EnterpriseThirdPartAccount::query()
                ->where([
                    'enterprise_id'  => $enterprise->id,
                    'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                    'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                    'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                    'type'           => $auth_type // 小程序类型
                ])->first();

            $openPlatformConfig = EnterpriseOpenPlatform::query()
                ->where('id', $enterprise->open_platform_id)
                ->first(['app_id', 'secret', 'token', 'aes_key'])
                ->toArray();

            // 初始化待返数据
            $data = [];
            if ($authInfo) {
                $data = [
                    'enterprise_id'        => $authInfo->enterprise_id,
                    'refresh_token'        => $authInfo->refresh_token,
                    'app_id'               => $authInfo->app_id,
                    'name'                 => $authInfo->name,
                    'main_qrcode'          => $authInfo->main_qrcode,
                    'authorize_type'       => $authInfo->authorize_type,
                    'wx_category_ids'      => $authInfo->wx_category_ids,
                    'platform'             => $authInfo->platform,
                    'type'                 => $authInfo->type,
                    'verify_status'        => $authInfo->verify_status,
                    'status'               => $authInfo->status,
                    'logo'                 => $authInfo->logo,
                    'version'              => $enterprise->version,
                    'hash'                 => $enterprise->hash,
                    'origin_id'            => $authInfo->origin_id,
                    'open_platform_config' => $openPlatformConfig
                ];
                $redis->setex($key, 3600, Json::encode($data));
            }
        } else {
            $data = Json::decode($data);
        }

        return success(['data' => $data]);
    }
}
