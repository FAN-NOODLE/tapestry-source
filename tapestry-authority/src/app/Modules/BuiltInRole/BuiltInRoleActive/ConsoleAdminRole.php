<?php

namespace App\Modules\BuiltInRole\BuiltInRoleActive;

use App\Models\EnterpriseAdmin;
use App\Models\RolePermissions;
use App\Services\EnterpriseAdminService;

class ConsoleAdminRole extends CommonRole
{
    public function accountOperation($params, $defaultAdd = true)
    {
        // 查找当前账号
        $location = strpos($this->enterprise->mobile, '-') ? strpos($this->enterprise->mobile, '-') + 1 : 0;
        $mobile = substr($this->enterprise->mobile, $location);
        $account = $mobile;
        if (isset($params['username'])) {
            $account = (new EnterpriseAdminService())->createDefaultAccount($params['username'], $mobile);
        }
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('account', 'like', "%{$account}")
            ->first();

        $superAdmin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('is_super', '=', EnterpriseAdmin::IS_SUPER)
            ->first();

        // 账号存在走更新
        if ($admin) {
            $admin->update(['admin_available_system' => $superAdmin->admin_available_system]);
        } else {
            // 账号不存在走新建
            $salt = (new EnterpriseAdminService())->makeSalt();

            $admin = EnterpriseAdmin::query()->create([
                'enterprise_id'          => $this->enterprise->id,
                'username'               => $this->enterprise->contact,
                'account'                => isset($params['username']) ? ("+86-" . $account) : $account,
                'password'               => md5(md5($account) . $salt),
                'salt'                   => $salt,
                'global_roaming'         => '+86',
                'admin_available_system' => $superAdmin->admin_available_system,
                'mobile'                 => $mobile,
                'is_super'               => EnterpriseAdmin::NO_SUPER,
                'range_type'             => EnterpriseAdmin::RANGE_TYPE_ALL,
                'operator_id'            => $superAdmin->id
            ]);
        }

        // 绑定角色
        $this->adminBindRole($admin->id);

        // 拿到超管所有权限
        $superPerList = RolePermissions::query()
            ->leftJoin('role', 'role.id', '=', 'role_permissions.role_id')
            ->where('role.mark', '=', 'super_admin')
            ->where('role.enterprise_id', '=', $this->enterprise->id)
            ->pluck('role_permissions.permissions_id')
            ->toArray();

        // 获取管理员所有权限
        $oldPerMissIons = RolePermissions::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('role_id', '=', $this->role->id)
            ->pluck('permissions_id')
            ->toArray();
        //需新增的权限
        $addPerMissions = array_diff($superPerList, $oldPerMissIons);
        //需删除的权限
        $delPerMissions = array_diff($oldPerMissIons, $superPerList);
        // 清除管理员所有权限
        RolePermissions::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('role_id', '=', $this->role->id)
            ->whereIn('permissions_id', $delPerMissions)
            ->forceDelete();
        $list = [];
        // 组装管理员所有权限
        foreach ($addPerMissions as $item) {
            $list[] = [
                'role_id'        => $this->role->id,
                'enterprise_id'  => $this->enterprise->id,
                'permissions_id' => $item,
            ];
        }
        // 创建权限关联
        if ($list) {
            RolePermissions::query()->insert($list);
        }
    }
}