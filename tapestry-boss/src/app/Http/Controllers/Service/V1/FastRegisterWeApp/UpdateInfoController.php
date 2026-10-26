<?php


namespace App\Http\Controllers\Service\V1\FastRegisterWeApp;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use App\Models\RegisterProgramAgent;
use App\Services\OpenPlatform\EnterpriseService;
use App\Services\OpenPlatform\OpenPlatformService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;


class UpdateInfoController extends BaseController
{


    public function index(Request $request, EnterpriseService $enterpriseService, OpenPlatformService $openPlatformService)
    {

        $params = $request->only([
            "app_id",
            'enterprise_id',
        ]);

        if (empty($params['enterprise_id'])) return error(600001);

//      $info = (new RegisterProgramAgent())->getInfoByEnterpriseId($params['enterprise_id']);
        $info = EnterpriseThirdPartAccount::query()->where('enterprise_id', $params['enterprise_id'])->first();
        if (empty($info)) return success();

//      if ($info['status'] != RegisterProgramAgent::STATUS_APPLY_SUCCESS) {
//          return notice(100000, ['小程序未申请成功，请稍后再试']);
//      }
        try {
            $authorizer = $openPlatformService->getAuthorizerConfigByAppId($info['app_id'], $params['enterprise_id']);

            // logo http转https
            $logo = !empty($authorizer['authorizer_info']['head_img']) ? str_replace('http:', 'https:', $authorizer['authorizer_info']['head_img']) : '';

            $thirdPartConfig = [
                'enterprise_id' => $params['enterprise_id'],
                'app_id'        => $authorizer['authorization_info']['authorizer_appid'],
                'name'          => $authorizer['authorizer_info']['nick_name'],
                'platform'      => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,
                'verify_status' => $authorizer['authorizer_info']['verify_type_info']['id'] == 0 ? EnterpriseThirdPartAccount::VERIFIED : EnterpriseThirdPartAccount::UNVERIFIED,
                'refresh_token' => $authorizer['authorization_info']['authorizer_refresh_token'],
                'logo'          => $logo,
            ];
            $enterpriseService->updateOrCreateThirdPartConfig($thirdPartConfig);
        } catch (\Exception $e) {
            Log::info($params['enterprise_id'] . ' 更新小程序信息失败 ' . $e->getMessage());
            return notice(100000, ['更新失败']);
        }

        return success();
    }

}
