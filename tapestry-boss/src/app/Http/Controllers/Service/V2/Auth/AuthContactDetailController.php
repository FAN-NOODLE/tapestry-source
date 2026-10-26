<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 12:11 PM
 */

namespace App\Http\Controllers\Service\V2\Auth;

use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Http\Request;
use App\Models\Enterprise;

class AuthContactDetailController extends BaseController
{
    /**
     * 关闭crs服务
     */
    const CRS_CLOSE = 2;

    /**
     * 外部联系人secret
     *
     * @param Request $request
     *
     * @return array|mixed
     * @author 陈朔  chenshuo@vchangyi.com
     * @date   2018/12/3 1:22 PM
     */
    public function index(Request $request)
    {
        $enterpriseId = $request->post('enterprise_id');
        $enterprise = Enterprise::query()->where(['id' => $enterpriseId])->first();
        if (empty($enterprise)) {
            // 企业信息不存在
            return error(600002);
        }

        $queryParams = [
            'enterprise_id' => $enterpriseId
        ];

        $callbackUrl = config('app.host') . '/callback/contact-install?' . http_build_query($queryParams);
        $installUrl = config('app.uc_install_contact_url') . '&callbackUrl=' . $callbackUrl;
        return successV2([
                'auth_url'        => $installUrl,
                'enterprise_id'   => $enterpriseId,
                'is_open'         => $enterprise->crs_status == self::CRS_CLOSE ? 0 : $enterprise->crs_status,
                'installed'       => $enterprise->contact_status == self::CRS_CLOSE ? false : true,
                'corp_id'         => $enterprise->corp_id ?: '',
                'external_secret' => $enterprise ? $enterprise->external_secret : '',
            ]);
    }
}
