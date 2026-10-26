<?php

namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseConfig;

class GetEnterpriseConfigController extends BaseController
{
    public $requireLogin = false;

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $enterpriseId = trim($this->request->post('enterprise_id'));
        if (empty($enterpriseId)) {
            notice($enterpriseId);
        }
        // 增加企业仓库类型参数 1.单仓 2.多仓
        $enterpriseSdk               = new EnterpriseConfig();
        $config                      = $enterpriseSdk->info(['enterprise_id' => $enterpriseId])['data'];
        $configByKey                 = array_column($config, 'value', 'key');
        return success($configByKey);
    }
}
