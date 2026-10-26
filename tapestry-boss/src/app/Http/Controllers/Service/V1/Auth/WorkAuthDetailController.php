<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/27
 * Time: 下午8:46
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\EnterpriseWorkSetting;
use EasyWeChat\OpenPlatform\Application;
use Illuminate\Http\Request;

class WorkAuthDetailController extends BaseController
{
    /**
     * 企业微信小程序授权详情
     *
     * @param Request     $request
     * @param Application $openPlatform
     *
     * @return array
     */
    public function index(Request $request, Application $openPlatform)
    {
        $enterpriseId = $request->input('enterprise_id');
        $url = $request->input('back_url');
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

        $shopAuthInfo = EnterpriseThirdPartAccount::query()
            ->where([
                'enterprise_id'  => $enterpriseId,
                'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
                'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
                'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
                'type'           => EnterpriseThirdPartAccount::TYPE_QY_WX // 商城小程序
            ])->first();

        $authUrl = '';
        $workSetting = null;
        $host = config('app.env') == 'test' ? config('client.api_url')['boss'] : config('app.host');

        if (!$shopAuthInfo) {
            $query_params = [
                'auth_for'      => 'WORK',
                'name'          => $enterprise->name,
                'enterprise_id' => $enterpriseId,
                'back'          => $host . '/callback/authorize-access-token'
            ];

            // 替换企业开放平台配置
            changeOpenPlatformConfig($enterpriseId, $openPlatform, 0, true);

            $callbackUrl = $url . '?' . http_build_query($query_params);
            $authUrl = $openPlatform->getPreAuthorizationUrl($callbackUrl);
        } else {
            $shopAuthInfo = $shopAuthInfo->toArray();
            // 如果已授权，则查询微信设置信息
            $workSetting = EnterpriseWorkSetting::query()->where(['enterprise_id' => $enterpriseId])->first();
        }

        return success([
            'data' => [
                'enterprise_id'       => $enterpriseId,
                'installed'           => (bool)$shopAuthInfo,
                'name'                => $shopAuthInfo ? $shopAuthInfo['name'] : '',
                'logo'                => $shopAuthInfo ? $shopAuthInfo['logo'] : '',
                'auth_at'             => $shopAuthInfo ? date('Y-m-d H:i:s', $shopAuthInfo['created_at']) : '',
                'auth_url'            => $authUrl,
                'external_secret'     => $workSetting ? $workSetting->external_secret : '',
                'mini_program_secret' => $workSetting ? $workSetting->mini_program_secret : '',
                'corp_id'             => $workSetting ? $workSetting->corp_id : '',
                'pl_plugin_id'        => $workSetting ? $workSetting->pl_plugin_id : '',
                'is_open'             => $enterprise->crs_status
            ]
        ]);
    }
}
