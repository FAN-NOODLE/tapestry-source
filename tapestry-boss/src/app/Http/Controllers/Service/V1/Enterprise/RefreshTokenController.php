<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use App\Services\EnterpriseNew\EnterpriseNewService;
use Illuminate\Http\JsonResponse;
use Overtrue\LaravelWeChat\Facade as EasyWeChat;

class RefreshTokenController extends BaseController
{
    /**
     * 更新企业的RefreshToken
     *
     * @return JsonResponse
     */
    public function updateMiniProgramToken(): JsonResponse
    {
        $enterpriseID = $this->request->input('enterprise_id');

        $account = EnterpriseThirdPartAccount::query()
            ->where('enterprise_id', '=', $enterpriseID)
            ->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP)
            ->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE)
            ->where('verify_status', '=', EnterpriseThirdPartAccount::VERIFIED)
            ->first();
        if (empty($account)) {
            error(300502);
        }

        $enterprise = Enterprise::query()
            ->where('id', '=', $enterpriseID)
            ->where('status', '=', Enterprise::STATUS_ENABLE)
            ->first();
        if (empty($enterprise)) {
            error(201108);
        }

        $appsList = $this->authorizersMiniProgram($enterprise->open_platform_id);
        if (empty($appsList)) {
            error(300517);
        }

        $refreshToken = $this->getAccountRefreshToken($appsList, $account->app_id);
        if (empty($refreshToken)) {
            error(300518);
        }

        $account->update([
            'updated_at'    => time(),
            'refresh_token' => $refreshToken,
        ]);

        (new EnterpriseNewService())->clearAuthInfoCache($enterprise->id);

        return success([
            'status' => 'update_success'
        ]);
    }

    /**
     * 获取开放平台绑定的所有小程序信息
     *
     * @param $openPlatformId
     * @return array
     */
    protected function authorizersMiniProgram($openPlatformId): array
    {
        $offset = 0;
        $appListRes = [];
        $openPlatform = EasyWeChat::openPlatform();

        do {
            changeOpenPlatformConfig(0, $openPlatform, $openPlatformId);

            $appsList = $openPlatform->getAuthorizers($offset);

            if (!empty($appsList['list'])) {
                $offset = $offset + 500;
                $appListRes = array_merge($appListRes, $appsList['list']);
            }
        } while ($offset <= $appsList['total_count']);

        return $appListRes;
    }

    /**
     * 获取小程序的RefreshToken
     *
     * @param $appsList
     * @param $appId
     * @return string
     */
    protected function getAccountRefreshToken($appsList, $appId): string
    {
        foreach ($appsList as $item) {
            if ($item['authorizer_appid'] == $appId) {
                return $item['refresh_token'];
            }
        }

        return '';
    }
}
