<?php

namespace App\Http\Controllers\V1\Auth;

use CrsSdk\Boss\Auth;
use CrsSdk\Work\EnterprisePlugin;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Support\Facades\Log;

class DaoGouConfigController extends BaseController
{
    public function daoGouConfig()
    {
        $params                  = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;

        // 获取自建应用信息
        $originConfig = (new Auth())->getDaoGouConfig($params);

        // 更新自建应用信息用
        $params['qy_domain']       = $this->qyDomain;
        $params['pl_plugin_id']    = empty($originConfig['data']['pl_plugin_id']) ? '' : $originConfig['data']['pl_plugin_id'];
        $params['origin_secret']   = empty($originConfig['data']['secret']) ? '' : $originConfig['data']['secret'];
        $params['origin_agent_id'] = empty($originConfig['data']['agent_id']) ? '' : $originConfig['data']['agent_id'];

        // 安装或更新自建应用
        $result = (new EnterprisePlugin())->daoGouInstall($params);

        Log::info('安装智慧导购UC返回数据', ['data' => $result]);

        $params = array_merge($params, [
            'pl_plugin_id' => $result['data']['mini_program_install_result']['plPluginid'] ?? ''
        ]);

        // 更新存储的自建应用信息
        (new Auth())->daoGouInstall($params);

        return success();
    }

    /**
     * 获取导购配置
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function getDaoGouConfig()
    {
        $result = (new Auth())->getDaoGouConfig(['enterprise_id' => $this->enterpriseId]);

        return success($result['data']);
    }
}