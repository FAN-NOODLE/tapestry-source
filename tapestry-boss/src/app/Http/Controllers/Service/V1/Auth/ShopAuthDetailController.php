<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/27
 * Time: 下午8:46
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Common\RedisKey;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseWorkPlugin;
use EasyWeChat\OpenPlatform\Application;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use App\Library\Json;


class ShopAuthDetailController extends BaseController
{
    /**
     * 商城小程序授权详情
     *
     * @param Request $request
     * @param Application $openPlatform
     *
     * @return mixed
     */
    public function index(Request $request, Application $openPlatform)
    {
        $enterpriseId = $request->input('enterprise_id');
        $url          = $request->input('back_url');
        $client       = $request->input('client', '');
        if (empty($enterpriseId)) {
            return error(600001);
        }

        $enterprise = Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();

        if (empty($enterprise)) {
            // 企业信息不存在
            return error(600002);
        }

        $node1 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 1, 'request_duration' => (int)(($node1 - START_TIME) * 1000)]);

        $redis        = redis();
        $authKey      = RedisKey::MINI_PROGRAM_AUTH_INFO_BY_ENTERPRISE_ID . $enterpriseId;
        $shopAuthInfo = Json::decode($redis->get($authKey));
        $node2        = microtime(true);
        if (empty($shopAuthInfo)) {
            $shopAuthInfo = EnterpriseThirdPartAccount::query()
                ->where([
                    'enterprise_id'  => $enterpriseId,
                    'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                    'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                    'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                    'type'           => EnterpriseThirdPartAccount::TYPE_SHOP // 商城小程序
                ])->first();
            if (!empty($shopAuthInfo)) {
                $shopAuthInfo = $shopAuthInfo->toArray();
                $redis->set($authKey, Json::encode($shopAuthInfo));
            }

            $node2 = microtime(true);
            Log::info('shop-auth-detail-time', ['node' => 2, 'request_duration' => (int)(($node2 - $node1) * 1000)]);
        }

        $node3 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 3, 'request_duration' => (int)(($node3 - $node2) * 1000)]);

        $host = config('app.env') == 'test' ? config('client.api_url')['boss'] : config('app.host');

        $node4 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 4, 'request_duration' => (int)(($node4 - $node3) * 1000)]);

        $query_params = [
            'auth_for'      => 'SHOP',
            'name'          => $enterprise->name,
            'enterprise_id' => $enterpriseId,
            'back'          => $host . '/callback/authorize-access-token',
        ];

        Log::info("wx_query_params", $query_params);

        $node5 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 5, 'request_duration' => (int)(($node5 - $node4) * 1000)]);

        // 替换企业开放平台配置
        changeOpenPlatformConfig($enterpriseId, $openPlatform, 0, true);
        $node6 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 6, 'request_duration' => (int)(($node6 - $node5) * 1000)]);

        $authUrl = '';
        if ($client == 'console') {
            $callbackUrl = $url . '?' . http_build_query($query_params);
            $authUrl     = $openPlatform->getPreAuthorizationUrl($callbackUrl);
        }

        $node7       = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 7, 'request_duration' => (int)(($node7 - $node6) * 1000)]);

        $shopSetting = null;
        if ($shopAuthInfo) {
            $workPluginKey = RedisKey::ENTERPRISE_WORK_PLUGIN_BY_ENTERPRISE_ID . $enterpriseId;
            $shopSetting   = Json::decode($redis->get($workPluginKey));
            if (empty($shopSetting)) {
                $shopSetting = EnterpriseWorkPlugin::query()->where('enterprise_id', $enterpriseId)->where('identifier', 'crs-shop')->first();
                if (!empty($shopSetting)) {
                    $redis->setex($workPluginKey, 3600, Json::encode($shopSetting->toArray()));
                }
            }
        }
        $node8 = microtime(true);
        Log::info('shop-auth-detail-time', ['node' => 8, 'request_duration' => (int)(($node8 - $node7) * 1000)]);

        return success([
            'data' => [
                'enterprise_id'       => $enterpriseId,
                'installed'           => (bool)$shopAuthInfo,
                'name'                => $shopAuthInfo ? $shopAuthInfo['name'] : '',
                'logo'                => $shopAuthInfo ? $shopAuthInfo['logo'] : '',
                'auth_at'             => $shopAuthInfo ? date('Y-m-d H:i:s', $shopAuthInfo['created_at']) : '',
                'auth_url'            => $authUrl,
                'is_open'             => $enterprise->store_status,
                'mini_program_secret' => $shopSetting ? $shopSetting['secret'] : '',
                'corp_id'             => $enterprise ? $enterprise->corp_id : '',
                'pl_plugin_id'        => $shopSetting ? $shopSetting['pl_plugin_id'] : '',
                'agent_id'            => $shopSetting ? $shopSetting['agent_id'] : '',
                'app_secret'          => $shopAuthInfo ? $shopAuthInfo['app_secret'] : '',
                'release_branch'      => $shopAuthInfo['release_branch'] ?? '',
                'release_tag'         => $shopAuthInfo['release_tag'] ?? '',
                'ext_build_json'      => $shopAuthInfo ? json_decode($shopAuthInfo['ext_build_json'], true) : ''
            ]
        ]);
    }
}
