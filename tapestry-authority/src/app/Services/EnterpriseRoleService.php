<?php

namespace App\Services;

use Exception;
use App\Models\EnterpriseRole;
use Illuminate\Support\Facades\Log;
use Illuminate\Contracts\Pagination\LengthAwarePaginator;

class EnterpriseRoleService extends BaseService
{
    /**
     * 角色权限类型
     */
    const ROLE_TYPE = [
        EnterpriseRole::ROLE_TYPE_POWER,
        EnterpriseRole::ROLE_TYPE_DATA,
    ];

    /**
     * 保存角色
     *
     * @param $data
     * @return int|void
     */
    public function saveRole($data)
    {
        try {
            $res = EnterpriseRole::query()->create($data);
        } catch (Exception $e) {
            Log::error('保存角色失败：' . $e->getMessage());
            error(700000, ['保存角色失败']);
        }

        return $res->id;
    }

    /**
     * 编辑角色
     * @param $where
     * @param $data
     * @return int|void
     */
    public function updateRole($where, $data)
    {
        try {
            return EnterpriseRole::query()->where($where)->update($data);
        } catch (Exception $e) {
            Log::error('编辑角色失败：' . $e->getMessage());
            error(700000, ['编辑角色失败']);
        }
    }

    /**
     * 删除角色
     * @param $roleId
     * @return int|void
     */
    public function deleteRole($roleId)
    {
        try {
            return EnterpriseRole::query()->where('id', $roleId)->delete();
        } catch (Exception $e) {
            Log::error('删除角色失败：' . $e->getMessage());
            error(700000, ['删除角色失败']);
        }
    }

    /**
     * 获取企业角色列表
     * @param $enterpriseId
     * @param array $roleIds
     * @return array
     */
    public function getRoleList($enterpriseId, $roleIds = [])
    {
        return EnterpriseRole::query()
            ->when(!empty($roleIds), function ($query) use ($roleIds) {
                return $query->whereIn('id', $roleIds);
            })
            ->where('enterprise_id', $enterpriseId)
            ->get()
            ->toArray();
    }

    /**
     * 根据角色ID获取角色信息
     *
     * @param $roleId
     * @return array|void
     */
    public function getRoleById($roleId)
    {
        try {
            $role = EnterpriseRole::query()->find($roleId);
        } catch (Exception $e) {
            Log::error('根据角色ID获取角色信息失败：' . $e->getMessage());
            error(700000, ['根据角色ID获取角色信息失败']);
        }

        return !empty($role) ? $role->toArray() : [];
    }

    /**
     * 获取角色详情信息
     *
     * @param $roleId
     * @return array|void
     */
    public function getRoleDetail($roleId)
    {
        try {
            $roleInfo = EnterpriseRole::query()
                ->where('id', $roleId)
                ->withCount("roleAdmins")
                ->with('enterpriseAdmin')
                ->first();
        } catch (Exception $e) {
            Log::error('获取角色详情信息失败：' . $e->getMessage());
            error(700000, ['获取角色详情信息失败']);
        }

        return !empty($roleInfo) ? $roleInfo->toArray() : [];
    }

    /**
     * 根据条件查询角色列表
     *
     * @param $where
     * @param $page
     * @param $size
     * @param $extParams
     *
     * @return LengthAwarePaginator|void
     */
    public function getRoleListByWhere($where, $page, $size, $extParams = [])
    {
        try {
            return EnterpriseRole::query()
                ->where($where)
                ->when(!empty($extParams['operator_id']), function ($query) use ($extParams) {
                    return $query->whereIn('operator_id', $extParams['operator_id']);
                })
                ->withCount("roleAdmins")
                ->with('enterpriseAdmin')
                ->when(empty($extParams['sort_type']), function ($query) {
                    $query->orderByDesc('created_at');
                })
                ->when(!empty($extParams['sort_type']), function ($query) use ($extParams) {
                    $query->orderBy('role_admins_count', $extParams['sort_type']);
                })
                ->paginate($size, ['*'], 'page', $page);
        } catch (Exception $e) {
            Log::error('根据条件查询角色列表失败：' . $e->getMessage());
            error(700000, ['根据条件查询角色列表失败']);
        }
    }

    /**
     * 根据条件查询角色基础信息
     *
     * @param $where
     * @return mixed
     */
    public function getRoleDetailByWhere($where)
    {
        try {
            return EnterpriseRole::query()->where($where)->first();
        } catch (Exception $e) {
            Log::error('根据条件查询角色基础信息失败：' . $e->getMessage());
            error(700000, ['根据条件查询角色基础信息失败']);
        }
    }

    /**
     * 根据企业ID获取所有的角色名
     * 1 可通过名字匹配
     * @param $enterpriseId
     * @param null $searchText
     * @return array|void
     */
    public function getRoleListByEnterpriseId($enterpriseId, $searchText = null)
    {
        try {
            $query = EnterpriseRole::query()
                ->select(['id', 'role_name'])
                ->where('enterprise_id', '=', $enterpriseId);

            if ($searchText) {
                $query->where('role_name', 'like', "%{$searchText}%");
            }
        } catch (Exception $e) {
            Log::error('根据企业ID获取所有的角色名失败：' . $e->getMessage());
            error(700000, ['根据企业ID获取所有的角色名失败']);
        }

        return $query->get()->toArray();
    }
}
