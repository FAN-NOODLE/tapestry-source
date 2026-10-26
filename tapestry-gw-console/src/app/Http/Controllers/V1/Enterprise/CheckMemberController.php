<?php

namespace App\Http\Controllers\V1\Enterprise;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePlugin;
use CrsSdk\Boss\EnterpriseWorkPlugin;
use Illuminate\Support\Facades\Log;

class CheckMemberController extends BaseController
{
    /**
     * 检查企业是否可见员工
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $data = [
            'status' => 'ok'
        ];

        $pluginList = (new EnterprisePlugin())->lists([
            'qy_domain' => $this->enterprise['domain']
        ]);
        $pluginIdList = array_column($pluginList, 'plIdentifier');

        Log::debug('uc插件列表', $pluginIdList);

        // 如果uc没有安装通讯录就判断是否按照自建应用
        if (!in_array('contact', $pluginIdList)) {
            $data['status'] = 'err';

            $usePluginList = (new EnterpriseWorkPlugin())->lists([
                'enterprise_id' => $this->enterprise['id']
            ]);
            $usePluginIdList = array_column($usePluginList, 'identifier');

            Log::debug('使用的插件列表', $usePluginIdList);

            foreach ($usePluginList as $key => $value) {
                if ($value['agent_id'] !== '' && $value['secret']) {
                    $data['status'] = 'ok';
                    break;
                }
            }
        }

        return success($data);
    }
}
