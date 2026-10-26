<?php


namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\Work\Auth;


class ScrmGetTokenController extends BaseController
{
    public $requireLogin = false;

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $adminId = trim($this->request->post('admin_id'));

        if (empty($adminId)) {
            notice(100007);
        }

        $authInfo  = (new Account())->getToken(['admin_id' => $adminId, 'client' => self::CLIENT]);
        $adminInfo = $authInfo['enterprise_admin_info'];
        $hash      = $adminInfo['is_split_flow'] ? $adminInfo['hash'] : ($adminInfo['version'] == 1 ? 'public' : 'store');

        // 增加企业仓库类型参数 1.单仓 2.多仓
        $enterpriseSdk               = new EnterpriseConfig();
        $config                      = $enterpriseSdk->info(['enterprise_id' => $adminInfo['enterprise_id']])['data'];
        $configByKey                 = array_column($config, 'value', 'key');
        $adminInfo['warehouse_type'] = $configByKey['warehouse_type'] ?? 1;

        $domain = $authInfo['enterprise_admin_info']['domain'] ?? '';

        $jsApiConfig = (new Auth())->agentConfig([
                'domain'        => $domain,
                'appIdentifier' => 'crs_store_assistant'
            ]) ?? [];

        $data = array_merge([
            'hash'          => $hash,
            'data'          => $adminInfo['version'] == 1 ? $this->getAuthorityDataForV1($adminInfo, $adminId) : [],
            'total'         => 1,
            'version'       => $adminInfo['version'],
            'admin_list'    => [$adminInfo],
            // 增加前端js api config
            'js_api_config' => $jsApiConfig
        ], $authInfo);

        return success($data);
    }

}
