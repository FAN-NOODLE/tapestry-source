<?php
/**
 * Created by PhpStorm.
 * User: qihongchuan
 * Date: 2020-02-12
 * Time: 10:48
 */

namespace App\Http\Controllers\Service\V1\MiniProgram;

use App\Models\EnterpriseSdkConfig;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class SdkConfigController extends MiniBaseController
{

    /**
     * 小程序 - 数据收集sdk - 配置查询接口
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id', 0);
        $sdkConfig = EnterpriseSdkConfig::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('app_type', EnterpriseSdkConfig::TYPE_SHOP)
            ->first();
        $sdkConfig = empty($sdkConfig) ? null : $sdkConfig->toArray();

        return [
            'data' => $sdkConfig,
        ];
    }
}
