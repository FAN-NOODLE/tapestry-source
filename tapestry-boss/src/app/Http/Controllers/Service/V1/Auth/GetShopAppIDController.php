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
use Illuminate\Http\Request;

class GetShopAppIDController extends BaseController
{
    /**
     * 获取企业商城小程序id
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');
        if (empty($enterpriseId)) {
            return error(600001);
        }

        $enterprise = Enterprise::query()->where('id', $enterpriseId)->first();

        if (empty($enterprise)) {
            // 企业信息不存在
            return error(600002);
        }

        $shopAuthInfo = EnterpriseThirdPartAccount::query()->where([
            'enterprise_id'  => $enterpriseId,
            'authorize_type' => EnterpriseThirdPartAccount::AUTHORIZE_TYPE_THIRD, // 第三方授权
            'status'         => EnterpriseThirdPartAccount::STATUS_ACTIVE, // 有效状态
            'platform'       => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT,// 小程序
            'type'           => EnterpriseThirdPartAccount::TYPE_SHOP // 商城小程序
        ])->first();

        $shopAppId = '';
        $releaseTag = '';
        if ($shopAuthInfo) {
            $shopAppId = $shopAuthInfo->app_id;
            $releaseTag = $shopAuthInfo->release_tag;
        }

        return success([
            'data' => [
                'shop_app_id' => $shopAppId,
                'release_tag' => $releaseTag
            ]
        ]);
    }
}
