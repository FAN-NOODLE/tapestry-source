<?php

namespace App\Modules\Rbac;

use App\Models\Permissions as PermissionsModel;
use App\Models\Role as RoleModel;
use App\Models\RolePermissions as RolePermissionsModel;
use App\Modules\BuiltInRole\RoleStrategyConfig;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

/**
 * 角色权限
 *
 * @package App\Modules\Rbac
 */
class RolePermissions extends Permissions
{
    /**
     * @var array
     */
    private $roleIdList = [];

    /**
     * @var int
     */
    private $enterpriseId;

    /**
     * 角色权限关系 加载后暂存起来方便后面直接用 减少查询
     *
     * @var array
     */
    private $rolePermissionsRelation = [];

    private $fields = [
        'role_id',
        'permissions_id'
    ];

    /**
     * RolePermissions constructor.
     *
     * @param array $roleIdList
     * @param int $version
     * @param int $enterpriseId
     * @throws Exception
     */
    public function __construct($roleIdList, $version, $enterpriseId)
    {
        parent::__construct($version);

        $this->roleIdList = $roleIdList;
        $this->enterpriseId = $enterpriseId;
    }

    /**
     * 获取当前操作人角色的所有权限列表
     *
     * @return array
     */
    public function lists()
    {
        $permissionsIdList = RolePermissionsModel::query()
            ->whereIn('role_id', $this->roleIdList)
            ->pluck('permissions_id')
            ->unique()
            ->values()
            ->toArray();

        return $this->all([
            ['id', 'in', $permissionsIdList]
        ]);
    }

    /**
     * 当前角色拥有的迭代号列表
     *
     * @return array
     */
    public function versionTagList()
    {
        $list = RolePermissionsModel::query()
            ->leftJoin('permissions', 'role_permissions.permissions_id', '=', 'permissions.id')
            ->whereIn('role_permissions.role_id', $this->roleIdList)
            ->where('role_permissions.enterprise_id', '=', $this->enterpriseId)
            ->where('permissions.deleted_at', '=', 0)
            ->pluck('permissions.version_tag')
            ->toArray();

        return array_values(array_unique($list));
    }

    /**
     * 获取一个角色的所有权限ID列表
     *
     * @param int $roleId
     * @return array
     */
    public function getList($roleId)
    {
        $list = RolePermissionsModel::query()
            ->select($this->fields)
            ->where('role_id', $roleId)
            ->pluck('permissions_id')
            ->toArray();
        $this->rolePermissionsRelation[$roleId] = $list;

        Log::info('预加载角色权限', [
            '角色ID' => $roleId
        ]);

        return $list;
    }

    /**
     * 角色权限树形结构 只能上级查下级的
     *
     * @param int $roleId 要查看的角色ID
     * @return array
     * @throws Exception
     */
    public function tree($roleId)
    {
        // 查出当前操作人角色和被操作人对应的权限列表 如果预加载了就不查询了
        if (isset($this->rolePermissionsRelation[$roleId])) {
            $roleIdList = $this->roleIdList;
        } else {
            $roleIdList[] = $roleId;
        }

        $permissionsList = RolePermissionsModel::query()
            ->select($this->fields)
            ->whereIn('role_id', $roleIdList)
            ->when(isset($this->rolePermissionsRelation[$roleId]), function ($query) {
                return $query->pluck('permissions_id');
            }, function ($query) {
                return $query->get();
            })
            ->toArray();

        if (isset($this->rolePermissionsRelation[$roleId])) {
            $rolePermissionsIdList = $this->rolePermissionsRelation[$roleId];
        } else {
            // 被操作人的权限ID列表
            $rolePermissionsIdList = [];
            foreach ($permissionsList as $key => $value) {
                if ($value['role_id'] === $roleId) {
                    $rolePermissionsIdList[] = $value['permissions_id'];
                    unset($permissionsList[$key]);
                }
            }
        }

        // 操作人的权限ID列表 默认操作人角色就是被操作角色
        $permissionsIdList = $rolePermissionsIdList;

        // 判断不是超管看自己
        if (count($this->roleIdList) > 1 && $this->roleIdList[0] !== $roleId) {
            $permissionsIdList = array_column($permissionsList, 'permissions_id');
        }

        // 取出当前操作人所有权限,看被操作的角色有哪些权限 has为1反之0
        $validAll = $this->all([
            ['id', 'in', $permissionsIdList]
        ]);
        foreach ($validAll as $key => $value) {
            if (in_array($value['id'], $rolePermissionsIdList)) {
                $validAll[$key]['has'] = 1;
            } else {
                $validAll[$key]['has'] = 0;
            }
        }

        return $this->getChildrenTreeFilterEmptyMenu($validAll, $roleIdList);
    }

    /**
     * 查看自己的权限
     *
     * @return array
     * @throws Exception
     */
    public function myTree()
    {
        $permissionsIdList = RolePermissionsModel::query()
            ->select($this->fields)
            ->whereIn('role_id', $this->roleIdList)
            ->pluck('permissions_id')
            ->unique()
            ->values()
            ->toArray();

        $validAll = $this->all([
            ['id', 'in', $permissionsIdList]
        ]);

        return $this->getChildrenTreeFilterEmptyMenu($validAll, $this->roleIdList);
    }

    public function deleteByRoleId($roleId)
    {
        RolePermissionsModel::query()->where('role_id', '=', $roleId)->forceDelete();
    }

    public function deleteByPermissionsId($permissionsIdList)
    {
        if ($permissionsIdList !== []) {
            RolePermissionsModel::query()
                ->where('enterprise_id', '=', $this->enterpriseId)
                ->whereIn('permissions_id', $permissionsIdList)
                ->forceDelete();
        }
    }

    /**
     * 超管权限发布
     *
     * @param int $systemType 系统类型
     * @param string $versionTag 迭代号
     * @return string 注意此返回值调用者判断结果了
     */
    public function superAdminRelease(int $systemType, string $versionTag): string
    {
        $permissionList = PermissionsModel::query()->when($systemType !== 0, function ($query) use ($systemType) {
                return $query->where('system_type', '=', $systemType);
            })->when($versionTag !== 'all', function ($query) use ($versionTag) {
                return $query->where('version_tag', '=', $versionTag);
            })->get()->toArray();
        if ($permissionList === []) {
            return '没有权限需要处理';
        }

        $superAdminRoleID = $this->roleIdList[0];
        $permissionIDList = array_column($permissionList, 'id');

        // 找出当前要发布权限中的选配权限及选配权限下级
        $additionalPermissionsIDList = array_column($this->filterPermissions($permissionList, PermissionsModel::VERSION_SELECT_CONFIG), 'id');
        foreach ($additionalPermissionsIDList as $id) {
            $additionalPermissionsIDList = array_merge($additionalPermissionsIDList, $this->getChildrenIDList($id));
        }

        // 去差拿到要删除的权限ID列表，这里要删掉部分旧的再插入，例如权限删除了 达到清理角色权限表数据的目的
        $deletePermissionsIDList = array_diff($permissionIDList, $additionalPermissionsIDList);

        // 找出内置管理员,给这个角色也发布上权限
        $adminRole = RoleModel::query()
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->where('mark', '=', RoleStrategyConfig::ROLE['console_admin_role']['mark'])
            ->first();
        $adminRoleID = 0;
        if ($adminRole !== null) {
            $adminRoleID = $adminRole->getAttribute('id');
        }
        try {
            DB::beginTransaction();

            RolePermissionsModel::query()
                ->whereIn('role_id', [$superAdminRoleID, $adminRoleID])
                ->whereIn('permissions_id', $deletePermissionsIDList)
                ->delete();

            $insertData = [];
            $adminRoleInsertData = [];
            $thisTime = time();
            foreach ($this->filterPermissions($permissionList, PermissionsModel::VERSION_STANDARD_CONFIG) as $permission) {
                $insertData[] = [
                    'role_id'        => $superAdminRoleID,
                    'permissions_id' => $permission['id'],
                    'enterprise_id'  => $this->enterpriseId,
                    'created_at'     => $thisTime,
                    'updated_at'     => $thisTime
                ];
                if ($adminRoleID !== 0) {
                    $adminRoleInsertData[] = [
                        'role_id'        => $adminRoleID,
                        'permissions_id' => $permission['id'],
                        'enterprise_id'  => $this->enterpriseId,
                        'created_at'     => $thisTime,
                        'updated_at'     => $thisTime
                    ];
                }
            }
            RolePermissionsModel::query()->insert($insertData);
            if ($adminRoleID !== 0) {
                RolePermissionsModel::query()->insert($adminRoleInsertData);
            }

            DB::commit();
        } catch (Exception $exception) {
            DB::rollBack();

            Log::error('发布超管角色权限异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            return "发布超管角色权限异常：{$exception->getMessage()}";
        }

        return '发布超管角色权限完成';
    }

    /**
     * 根据本本配置过滤权限
     *
     * @param array $permissions
     * @param int $config
     * @return array
     */
    private function filterPermissions(array $permissions, int $config): array
    {
        $data = [];
        foreach ($permissions as $permission) {
            if ($permission[$this->getVersionField()] === $config) {
                $data[] = $permission;
            }
        }

        return $data;
    }

    /**
     * 获取指定元素下的所有子级树过滤空菜单
     *
     * @param array $data
     * @param array $roleIdList
     * @param int $id
     * @param int $level
     * @return array 树形结构
     * @throws Exception
     */
    private function getChildrenTreeFilterEmptyMenu($data, $roleIdList, $id = 0, $level = 0)
    {
        $tree = [];
        $emptyMenuIdList = [];
        foreach ($data as $key => $value) {
            if ($value['pid'] === $id) {
                $value['level'] = $level + 1;
                $value['children'] = $this->getChildrenTreeFilterEmptyMenu($data, $roleIdList, $value['id'], $value['level']);

                // 删除空菜单为,暂时没有更好的做法了,先在查看角色的时候删除掉
                if ($value['type'] === PermissionsModel::TYPE_MENU && $value['children'] === []) {
                    $emptyMenuIdList[] = $value['id'];
                    continue;
                }

                $tree[] = $value;
                unset($data[$key]);
            }
        }
        if ($emptyMenuIdList !== []) {
            RolePermissionsModel::query()
                ->where('enterprise_id', '=', $this->enterpriseId)
                ->whereIn('role_id', $roleIdList)
                ->whereIn('permissions_id', $emptyMenuIdList)
                ->delete();

            Log::info('清理角色空菜单权限', [
                'empty_menu_id_list' => $emptyMenuIdList,
                'role_id_list'       => $roleIdList
            ]);
        }

        return $tree;
    }

    public function setRoleIdList($roleIdList)
    {
        $this->roleIdList = $roleIdList;
    }
}
