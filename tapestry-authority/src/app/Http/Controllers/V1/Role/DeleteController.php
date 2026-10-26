<?php


namespace App\Http\Controllers\V1\Role;


use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseRoleService;

class DeleteController extends BaseController
{
    public function index(
        EnterpriseRoleService $enterpriseRoleService,
        EnterpriseAdminRoleService $enterpriseAdminRoleService
    )
    {
        $roleId = $this->request->input('role_id');
        if (empty($roleId)) {
            notice(702001);
        }

        // 查询角色是否存在
        $role = $enterpriseRoleService->getRoleById($roleId);
        if (empty($role)) {
            notice(702021);
        }

        // 店长角色和店员角色不能删除
        if ($role['identity_type'] != 1) {
            notice(702022);
        }

        // 删除角色
        $enterpriseRoleService->deleteRole($roleId);
        // 删除员工与角色的关联关系
        $enterpriseAdminRoleService->delAdminRoleByRoleId($roleId);

        return success([
            'name' => $role['role_name']
        ]);
    }
}
