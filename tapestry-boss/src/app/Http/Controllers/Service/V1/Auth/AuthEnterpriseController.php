<?php
/**
 * Created by PhpStorm.
 * User: daijun
 * Date: 2018/11/27
 * Time: 下午8:46
 */

namespace App\Http\Controllers\Service\V1\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;

class AuthEnterpriseController extends BaseController
{
    /**
     * 小程序通过app_id获取企业ID
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $appId = $request->input('app_id', '');
        if (empty($appId)) {
            return error(201103);
        }

        $authInfo = EnterpriseThirdPartAccount::query()
            ->where([
                'app_id'   => $appId,
                'status'   => EnterpriseThirdPartAccount::STATUS_ACTIVE,
                'platform' => EnterpriseThirdPartAccount::MINI_PROGRAM_ACCOUNT
            ])->first(['enterprise_id']);

        // 返回企业ID
        return success([
            'data' => [
                'enterprise_id' => $authInfo->enterprise_id ?? ''
            ]
        ]);
    }
}
