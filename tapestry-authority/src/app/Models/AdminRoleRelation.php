<?php

namespace App\Models;

class AdminRoleRelation extends Base
{
    /**
     * 表名称:管理员角色关系表
     *
     * @var string
     */
    protected $table = 'admin_role_relation';

    protected $guarded = [];

    /**
     * 获取角色IDs通过管理员ID
     * @param $adminId
     * @return array
     */
    public static function roleIdsByAdminId($adminId)
    {
        return AdminRoleRelation::query()
            ->where('admin_id', '=', $adminId)
            ->pluck('role_id')
            ->toArray();
    }

    /**
     * 获取角色关联的管理员数量
     * @param $roleId
     * @return int
     */
    public static function adminRelationNums($roleId)
    {
        return self::query()
            ->leftJoin('enterprise_admin', 'enterprise_admin.id', 'admin_role_relation.admin_id')
            ->where("enterprise_admin.deleted_at", 0)
            ->where('admin_role_relation.role_id', '=', $roleId)
            ->count();
    }

    /**
     * 获取角色关联的管理员ID集合
     * @param $roleId
     * @return array
     */
    public static function adminRelationIds($roleId)
    {
        if (!is_array($roleId)) {
            $roleId = [$roleId];
        }

        return self::query()->whereIn('role_id', $roleId)->pluck('admin_id')->toArray();
    }

    /**
     * 根据账号ID查询管理员拥有的角色
     * @param $adminId
     * @param $eId
     * @return array
     */
    public static function adminRoles($adminId, $eId)
    {
        return self::query()
            ->select(['role.name', 'role.mark', 'role.id as role_id'])
            ->leftJoin('role', 'admin_role_relation.role_id', 'role.id')
            ->where('admin_role_relation.enterprise_id', '=', $eId)
            ->where('admin_role_relation.admin_id', '=', $adminId)
            ->get()
            ->toArray();
    }
}
