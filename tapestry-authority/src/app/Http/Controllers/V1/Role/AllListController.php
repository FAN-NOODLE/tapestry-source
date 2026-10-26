<?php

namespace App\Http\Controllers\V1\Role;

use App\Services\EnterpriseRoleService;
use App\Http\Controllers\V1\BaseController;

class AllListController extends BaseController
{
    public function index(EnterpriseRoleService $enterpriseRoleService)
    {
        $searchText   = $this->request->input('search_text');
        $enterpriseId = $this->request->input('enterprise_id');

        if (empty($enterpriseId)) {
            notice(700014);
        }

        $roleNameList = $enterpriseRoleService->getRoleListByEnterpriseId($enterpriseId, $searchText);

        return success([
            'list' => $roleNameList
        ]);
    }
}