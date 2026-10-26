<?php
/**
 * 获取某导购的管理后台权限id集合
 */
namespace App\Http\Controllers\V1\Account;

use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseAdminRole;
use App\Http\Controllers\V1\BaseController;
use App\Models\RolePermissions;

class AccountPermissionsController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'admin_id',
            'mobile',
            'enterprise_id'
        ]);
        if (empty($input['enterprise_id'])) {
            error(700014);
        }

        $adminId = $input['admin_id'] ?? 0;
        if (!empty($input['mobile'])) {
            $adminId = EnterpriseAdmin::query()->where(['enterprise_id' => $input['enterprise_id']])
                ->where('mobile', $input['mobile'])
                ->value('id');
        }
        if (empty($adminId)) {
            error(700024);
        }
        //获取账号角色id集合
        $roleIds = EnterpriseAdminRole::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->where('enterprise_admin_id', $adminId)
            ->pluck('role_id')
            ->toArray();
        $permissions = []; //初始化权限id集合
        if (!empty($roleIds)) {
            //获取权限id集合
            $permissions = RolePermissions::query()
                ->where('enterprise_id', $input['enterprise_id'])
                ->whereIn('role_id', $roleIds)
                ->pluck('permissions_id')
                ->toArray();
        }
        return success(['permissions' => array_values(array_unique($permissions))]);
    }
}
