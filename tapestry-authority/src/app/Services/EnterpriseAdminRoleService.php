<?php
/**
 * Created by PhpStorm.
 * User: wpf
 * Date: 2020/4/12
 * Time: 16:52
 */

namespace App\Services;

use App\Models\AdminRoleRelation;
use App\Models\EnterpriseAdminRole;
use App\Modules\Rbac\AdminRole;
use Exception;

class EnterpriseAdminRoleService extends BaseService
{
    /**
     * 插入账号角色关联表
     * @param $input
     * @return bool
     * @throws Exception
     */
    public function addEnterpriseAdminRole($input)
    {
        (new AdminRole($input['enterprise_admin_id']))->bind($input['role_ids']);

        $data = [];
        foreach ($input['role_ids'] as $role_id) {
            $data[] = [
                'enterprise_id'       => $input['enterprise_id'],
                'enterprise_admin_id' => $input['enterprise_admin_id'],
                'role_id'             => $role_id,
            ];
        }

        return EnterpriseAdminRole::query()->insert($data);
    }

    /**
     * 更新管理员账号关联角色数据
     * @param $input
     * @return bool
     * @throws Exception
     */
    public function editEnterpriseAdminRole($input, $isDel = true)
    {
        if (!isset($input['role_ids'])){
            return false;
        }

        // 传输的角色ID为空 则删除该账号下所有角色
        if (empty($input['role_ids'])) {
            $this->delAdminRoleByAdminId($input['id']);
        }

        $adminRoleList = $this->getAdminRoleByAdminId($input['id']);
        $roleIds = array_column($adminRoleList, 'role_id');

        // 取需要插入的id
        $inserRoleIds = array_diff($input['role_ids'], $roleIds);

        // 取需要删除的Id
        $delRoleIds = array_diff($roleIds, $input['role_ids']);

        if (!empty($inserRoleIds)) {
            $this->addEnterpriseAdminRole([
                'role_ids'            => $inserRoleIds,
                'enterprise_id'       => $input['enterprise_id'],
                'enterprise_admin_id' => $input['id']
            ]);
        }

        // isdel为true时才删除减量数据
        if (!empty($delRoleIds) && $isDel) {
            $this->delAdminRoleByAdminId($input['id'], $delRoleIds);
        }

        // 新角色管理系统
        $this->newRoleAdminDetail($input);

        return true;
    }

    /**
     * 根据管理员ID获取角色数据列表
     * @param $adminId
     * @return array
     */
    public function getAdminRoleByAdminId($adminId)
    {
        return EnterpriseAdminRole::query()
            ->where('enterprise_admin_id', $adminId)
            ->get()
            ->toArray();
    }

    /**
     * 删除管理员ID关联的角色
     * @param $adminId
     * @param array $roleIds
     * @return bool|mixed|null
     * @throws Exception
     */
    public function delAdminRoleByAdminId($adminId, $roleIds = [])
    {
        EnterpriseAdminRole::query()
            ->when(!empty($roleIds), function ($query) use ($roleIds) {
                return $query->whereIn('role_id', $roleIds);
            })
            ->where('enterprise_admin_id', $adminId)
            ->delete();

        AdminRoleRelation::query()
            ->when($roleIds, function ($query) use ($roleIds) {
                return $query->whereIn('role_id', $roleIds);
            })
            ->where('admin_id', '=', $adminId)
            ->forceDelete();

        return '';
    }

    /**
     * 删除角色与管理员的关系
     *
     * @param int|array $roleIds
     * @return bool|mixed|void|null
     */
    public function delAdminRoleByRoleId($roleIds)
    {
        try {
            if (is_array($roleIds)) {
                return EnterpriseAdminRole::query()->whereIn('role_id', $roleIds)->delete();
            } else {
                return EnterpriseAdminRole::query()->where('role_id', $roleIds)->delete();
            }
        } catch (Exception $e) {
            error(600000);
        }
    }

    /**
     * 获取账号关联角色列表
     * @param $enterpriseId
     * @param null $roleId
     * @return array
     */
    public function getAdminRoleList($enterpriseId, $roleId = null, $adminIds = [])
    {
        return EnterpriseAdminRole::query()
            ->leftJoin('enterprise_role', 'enterprise_admin_role.role_id', '=', 'enterprise_role.id')
            ->when($roleId != null && is_array($roleId), function ($query) use ($roleId) {
                return $query->whereIn('role_id', $roleId);
            })
            ->when($roleId != null && is_numeric($roleId), function ($query) use ($roleId) {
                return $query->where('role_id', $roleId);
            })
            ->when(!empty($adminIds), function ($query) use ($adminIds) {
                return $query->whereIn('enterprise_admin_id', $adminIds);
            })
            ->where('enterprise_admin_role.enterprise_id', $enterpriseId)
            ->get()
            ->toArray();
    }

    public function newRoleAdminDetail($input)
    {
        // 先清除 后绑定
        (new AdminRole($input['id']))->bind([]);
        // id为编辑时候参数 enterprise_admin_id为创建时候的参数
        (new AdminRole($input['id']))->bind($input['role_ids']);
    }
}
