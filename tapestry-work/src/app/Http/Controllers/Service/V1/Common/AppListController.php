<?php

namespace App\Http\Controllers\Service\V1\Common;

use App\Http\Controllers\Service\V1\BaseController;
use App\Exceptions\Error;
use App\Library\Uc\EnterprisePlugin;
use CrsSdk\Boss\Auth;
use Illuminate\Http\Request;

class AppListController extends BaseController
{
    /**
     * 获取企业商城小程序与绑定小程序列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $enterpriseId  = $request->post('enterprise_id', 0);   // 企业id
        $qyDomain      = $request->post('qy_domain', '');      // 企业domain
        $appIdentifier = $request->post('app_identifier', ''); // 应用标识

        // 初始化小程序列表
        $apps = [];

        // 商城小程序
        $authApp = (new Auth())->authDetail(['auth_type' => 1, 'enterprise_id' => $enterpriseId]);

        if ($authApp['data']) {
            $apps[] = [
                'app_id'   => $authApp['data']['app_id'],
                'app_name' => $authApp['data']['name'],
            ];
        }

        $apps[] = [
            'app_id'   => '',
            'app_name' => '其他',
        ];

        return success(['data' => $apps]);
    }
}
