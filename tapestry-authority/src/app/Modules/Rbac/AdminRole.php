<?php

namespace App\Modules\Rbac;

use App\Models\AdminRoleRelation as AdminRoleRelationModel;
use App\Models\EnterpriseAdmin;
use App\Models\Role as RoleModel;
use Exception;

/**
 * 管理员角色
 *
 * @package App\Modules\Rbac
 */
class AdminRole
{
    /**
     * 管理员ID
     *
     * @var int
     */
    private $adminId;

    public function __construct($adminId)
    {
        $this->adminId = $adminId;
    }

    /**
     * 获取角色列表
     *
     * @return array
     */
    public function roles()
    {
        return RoleModel::query()
            ->select(['role.*'])
            ->leftJoin('admin_role_relation', 'role.id', '=', 'admin_role_relation.role_id')
            ->where('admin_role_relation.admin_id', '=', $this->adminId)
            ->where('admin_role_relation.deleted_at', '=', 0)
            ->get()
            ->toArray();
    }

    /**
     * 管理员绑定角色
     *
     * @param array $roleIdList
     * @throws Exception
     */
    public function bind($roleIdList)
    {
        $relation = AdminRoleRelationModel::query()->where('admin_id', '=', $this->adminId)->pluck('role_id')->toArray();

        // 禁止对超管进行操作
        if (RoleModel::query()->whereIn('id', $relation)->where('mark', '=', RoleModel::ROLE['super_admin']['mark'])->exists()) {
            throw new Exception('超管禁止操作');
        }

        // TODO::暂时没验证账号必须存在一个角色
        $enterpriseId = EnterpriseAdmin::query()->where('id', '=', $this->adminId)->value('enterprise_id');
        AdminRoleRelationModel::query()->where('admin_id', '=', $this->adminId)->forceDelete();

        $relation = [];
        foreach ($roleIdList as $roleId) {
            $relation[] = [
                'admin_id'      => $this->adminId,
                'role_id'       => $roleId,
                'enterprise_id' => $enterpriseId
            ];
        }
        AdminRoleRelationModel::query()->insert($relation);
    }

    /**
     * 用管理员ID列表查角色相关信息 返回的key是管理员ID使用是不需要判断key是否存在,必然存在
     *
     * @param array $adminIdList 管理员ID列表
     * @return array
     */
    public function adminListRolesRelation($adminIdList)
    {
        $list = AdminRoleRelationModel::query()
            ->select([
                'admin_role_relation.admin_id',
                'role.id AS role_id',
                'role.pid AS role_pid',
                'role.name as role_name',
            ])
            ->leftJoin('role', 'admin_role_relation.role_id', '=', 'role.id')
            ->whereIn('admin_role_relation.admin_id', $adminIdList)
            ->get()
            ->toArray();

        $relationList = [];
        foreach ($adminIdList as $adminId) {
            // 这里保证传入的ID一定会有数据
            $relationList[$adminId] = [];
            foreach ($list as $key => $value) {
                if ($adminId === $value['admin_id']) {
                    $relationList[$adminId][] = [
                        'role_id'   => $value['role_id'],
                        'role_pid'  => $value['role_pid'],
                        'role_name' => $value['role_name'],
                    ];
                }
            }
        }

        return $relationList;
    }
}
