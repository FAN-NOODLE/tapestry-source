<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Models\AdminRoleRelation;
use App\Models\EnterpriseAdmin;

class AdminRbacController extends BaseController
{
    public function adminRoles()
    {
        $this->validate($this->request, [
            'mobile'        => 'required',
            'enterprise_id' => 'required',
        ]);

        $params = $this->request->all();

        $admin = EnterpriseAdmin::adminByAccount($params['mobile'], $params['enterprise_id']);

        if (empty($admin)) {
            $admin = EnterpriseAdmin::adminByMobile($params['mobile'], $params['enterprise_id']);
            if (empty($admin)) {
                error(700012);
            }
        }

        $roles = AdminRoleRelation::adminRoles($admin->id, $params['enterprise_id']);

        return success($roles);
    }
}