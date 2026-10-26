<?php
/**
 * Created by PhpStorm.
 * User: admin
 * Date: 2018/12/28
 * Time: 11:38
 */

namespace App\Http\Controllers\Service\V2\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;
use App\Models\EnterpriseAdmin;
use Illuminate\Http\Request;

class UcCallbackController extends BaseController
{
    /**
     * token
     */
    const TOKEN = 'REDACTED_FOR_SOURCE_SHARING';

    /**
     * key
     */
    const KEY = 'wxqyEncodingAeskey';

    /**
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');
        $adminId      = $request->input('admin_id', 0);

        $enterpriseInfo = \App\Models\Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();
        // 该企业未授权CRS小程序
        if (empty($enterpriseInfo) || empty($enterpriseInfo['corp_id']) || empty($enterpriseInfo['domain'])) {
            return successV2([
                'uc_callback_url' => '',
                'token'           => '',
                'key'             => '',
                'message'         => '该企业未授权CRS小程序'
            ]);
        }
        $adminInfo = EnterpriseAdmin::query()->find($adminId);
        $ucEnterprise = new Enterprise();
        return successV2([
            'uc_callback_url' => vsprintf(config('app.uc_customer_callback_url'), [$enterpriseInfo['corp_id']]),
            'token'           => $ucEnterprise->entSettingGet($enterpriseInfo['domain'], self::TOKEN)['data'][self::TOKEN] ?? '',
            'key'             => $ucEnterprise->entSettingGet($enterpriseInfo['domain'], self::KEY)['data'][self::KEY] ?? '',
            'is_super'        => $adminInfo->is_super ?? 2
        ]);
    }

    public function setToken(Request $request)
    {
        $enterpriseId = $request->input('enterprise_id');
        $wxqyToken    = $request->input('wxqy_token', '');

        $enterpriseInfo = \App\Models\Enterprise::query()
            ->where('id', $enterpriseId)
            ->first();

        // 该企业未授权CRS小程序
        if (empty($enterpriseInfo) || empty($enterpriseInfo['corp_id']) || empty($enterpriseInfo['domain'])) {
            return successV2([
                'message' => '该企业未授权CRS小程序'
            ]);
        }
        $ucEnterprise = new Enterprise();
        $ucEnterprise->entSettingToken($enterpriseInfo['domain'], $wxqyToken);
        return successV2([
            'token' => $ucEnterprise->entSettingGet($enterpriseInfo['domain'], self::TOKEN)['data'][self::TOKEN] ?? '',
        ]);
    }
}