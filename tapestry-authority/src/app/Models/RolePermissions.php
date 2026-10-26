<?php

namespace App\Models;

class RolePermissions extends Base
{
    /**
     * 表名称:角色权限表
     *
     * @var string
     */
    protected $table = 'role_permissions';

    protected $guarded = [];

    /**
     * 获取权限ID集合通过角色ID集合
     * @param $roleIds
     * @return array
     */
    public static function pIdsByRoleIds($roleIds)
    {
        return RolePermissions::query()
            ->whereIn('role_id', $roleIds)
            ->pluck('permissions_id')
            ->toArray();
    }
}
