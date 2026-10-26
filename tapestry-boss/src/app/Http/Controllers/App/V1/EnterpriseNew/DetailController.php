<?php

namespace App\Http\Controllers\App\V1\EnterpriseNew;

use App\Models\Industry;
use Illuminate\Http\Request;
use CrsSdk\MicroAuthority\Role;
use App\Services\EnterpriseConfigService;
use App\Http\Requests\V1\EnterpriseNewRequest;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseNew\EnterpriseNewService;

class DetailController extends BaseController
{
    /**
     * 获取企业详情
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $input = $request->only([
            'enterprise_id'
        ]);

        EnterpriseNewRequest::detail($input);

        $detail = (new EnterpriseNewService())->getEnterpriseDetail($input['enterprise_id'], false);
        
        // 获取行业名称
        $name = (new Industry())->idKeyName([$detail['industry_id']]);
        // 获取仓库类型
        $config = EnterpriseConfigService::getWarehouseConfig($input['enterprise_id']);
        if (empty($config['value'])) {
            error(201113);
        }

        // 0全部 1标配 2选配
        $tree = (new Role())->detail([
            'role_mark'      => 'super_admin',
            'enterprise_id'  => $input['enterprise_id'],
            'version_config' => 2,
        ]);

        $detail['industry']         = $name[$detail['industry_id']];
        $detail['warehouse_type']   = $config['value'];
        $detail['permission_list']  = $tree['permissions'];
        $detail['available_system'] = json_decode($detail['available_system'], true);

        return success(['data' => $detail]);
    }
}