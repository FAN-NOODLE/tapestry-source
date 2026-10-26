<?php

namespace App\Modules\Rbac;

use App\Common\CacheKey;
use App\Common\Constant;
use App\Library\CodeTimer;
use App\Library\Util\Logger;
use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Traits\RecursionTrait;
use App\Models\Role as RoleModel;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use App\Modules\BuiltInRole\RoleStrategy;
use App\Models\Permissions as PermissionsModel;
use App\Models\RolePermissions as RolePermissionsModel;
use App\Models\AdminRoleRelation as AdminRoleRelationModel;

/**
 * 角色
 *
 * @package App\Modules\Rbac
 */
class Role
{
    use RecursionTrait;

    /**
     * @var RolePermissions
     */
    public $permissions;

    /**
     * @var AdminRole
     */
    public $adminRole;

    /**
     * 当前操作人角色ID列表
     *
     * @var array
     */
    private $roleIdList = [];

    /**
     * 管理员ID 操作人
     *
     * @var int
     */
    private $adminId;

    /**
     * 企业ID
     *
     * @var int
     */
    private $enterpriseId;

    /**
     * 企业模型
     *
     * @var Enterprise
     */
    private $enterprise;

    /**
     * 企业超管角色模型
     *
     * @var RoleModel
     */
    private $superAdminRole;

    /**
     * 被操作的角色模型
     *
     * @var null|RoleModel
     */
    private $operationRole;

    /**
     * @var \Illuminate\Redis\Connections\Connection
     */
    private $redis;

    /**
     * @var CodeTimer
     */
    private $codeTimer;

    public function __construct($adminId, $enterpriseId)
    {
        $enterprise = Enterprise::query()->where('id', '=', $enterpriseId)->first();
        if ($enterprise === null) {
            throw new \Exception('企业数据不存在', Constant::NOT_FOUND_CODE);
        }

        $this->adminRole    = new AdminRole($adminId);
        $this->roleIdList   = array_column($this->adminRole->roles(), 'id');
        $this->permissions  = new RolePermissions($this->roleIdList, $enterprise->getAttribute('version'), $enterprise->getAttribute('id'));
        $this->adminId      = $adminId;
        $this->enterprise   = $enterprise;
        $this->enterpriseId = $enterprise->getAttribute('id');
        $this->redis        = redis();
        $this->codeTimer    = new CodeTimer();
    }

    /**
     * 角色详情
     *
     * @param $roleIdOrMark
     * @param int $versionConfig -1只查基础信息
     * @return array
     * @throws \Exception
     */
    public function detail($roleIdOrMark, $versionConfig = -1)
    {
        $this->checkAuth($roleIdOrMark);
        $roleArr = $this->operationRole->toArray();

        switch ($versionConfig) {
            case 0:
                $key  = sprintf(CacheKey::ROLE, $this->enterpriseId, $this->operationRole->getAttribute('id'));
                $data = $this->redis->get($key);
                if ($data === null) {
                    $roleArr['permissions_tree']    = $this->permissions->tree($this->operationRole->getAttribute('id'));
                    $roleArr['permissions_id_list'] = $this->permissions->getList($this->operationRole->getAttribute('id'));
                    $this->redis->setex($key, 3600, json_encode($roleArr, JSON_UNESCAPED_UNICODE + JSON_UNESCAPED_SLASHES));
                } else {
                    $roleArr = json_decode($data, true);
                }
                break;
            case PermissionsModel::VERSION_STANDARD_CONFIG:
                $roleArr['permissions'] = $this->permissions->getStandardConfig();
                break;
            case PermissionsModel::VERSION_SELECT_CONFIG:
                $roleArr['permissions'] = $this->permissions->additionalConfigList();
                break;
            case PermissionsModel::VERSION_NO_CONFIG:
                $roleArr['permissions'] = $this->permissions->noConfigList();
                break;
            default:
                return $roleArr;
        }

        // 如果不是要全部权限遍历检查下哪些有哪些没有
        if ($versionConfig !== 0) {
            $rolePermissionsIdList = RolePermissionsModel::query()
                ->where('role_id', '=', $this->operationRole->getAttribute('id'))
                ->pluck('permissions_id')
                ->toArray();
            foreach ($roleArr['permissions'] as $key => $value) {
                if (in_array($value['id'], $rolePermissionsIdList)) {
                    $roleArr['permissions'][$key]['has'] = 1;
                } else {
                    $roleArr['permissions'][$key]['has'] = 0;
                }
            }
        }

        return $roleArr;
    }

    /**
     * 初始化内置角色
     *
     * @param array $markList
     * @param array $permissionsIdList
     * @param int $adminId
     * @param array $filterPermissionsIdList
     * @return array
     * @throws \Exception
     */
    public function init($markList, $permissionsIdList, $adminId = 0, $filterPermissionsIdList = [])
    {
        if (in_array(RoleModel::ROLE['super_admin']['mark'], $markList) === true) {
            if ($adminId === 0) {
                throw new \Exception('初始化超管必须要管理员ID');
            }

            // 只能存在一个超管
            if ($this->getSuperAdminRole() !== null) {
                throw new \Exception('已存在超管角色');
            }
        }

        // 初始化超管角色,找出标配的权限跟传入的权限进行合并 初始化3个角色 管理员 店长 店员(权限为空)
        $roleList = [];
        // 标记当前初始化的角色是否有超管
        $hasSuperAdmin = false;
        foreach ($markList as $mark) {

            // 内置角色不走这了
            if (in_array($mark, ['area', 'store_manager', 'store_member'])) {
                continue;
            }

            // 内置角色的基础信息 下面switch进行每个角色选配功能的合并
            if (!isset(RoleModel::ROLE[$mark])) {
                throw new \Exception('内置角色信息不存在');
            }
            $role = RoleModel::ROLE[$mark];

            // 角色名称不能重复
            if ($this->checkNameExists($role['name'])) {
                throw new \Exception("{$role['name']}角色已存在");
            }

            switch ($mark) {
                case RoleModel::ROLE['super_admin']['mark']:
                    $hasSuperAdmin          = true;
                    $validPermissionsIdList = $this->buildSuperAdminPermissionsIdList($permissionsIdList);
                    break;
                default:
                    throw new \Exception('内置角色不存在');
            }

            // 过滤特殊权限以及依赖权限
            $filterPermissionsIdList = array_merge($filterPermissionsIdList, $this->filterDependList($validPermissionsIdList));
            $validPermissionsIdList  = $this->filterPermissions($validPermissionsIdList, $filterPermissionsIdList);

            sort($validPermissionsIdList);
            $roleList[] = [
                'name'                => $role['name'],
                'mark'                => $role['mark'],
                'permissions_id_list' => $validPermissionsIdList
            ];
        }

        // 新增角色
        $roles = [];
        $time  = time();
        try {
            DB::beginTransaction();
            foreach ($roleList as $value) {
                $role = RoleModel::query()->create([
                    'pid'           => 0,
                    'admin_id'      => $adminId,
                    'mark'          => $value['mark'],
                    'name'          => $value['name'],
                    'weight'        => 0,
                    'enterprise_id' => $this->enterpriseId
                ]);

                $rolePermissions = [];
                foreach ($value['permissions_id_list'] as $id) {
                    $rolePermissions[] = [
                        'role_id'        => $role->getAttribute('id'),
                        'permissions_id' => $id,
                        'enterprise_id'  => $this->enterpriseId,
                        'created_at'     => $time,
                        'updated_at'     => $time
                    ];
                }
                // 更新角色权限信息
                foreach (array_chunk($rolePermissions, 100) as $rolePermission) {
                    RolePermissionsModel::query()->insert($rolePermission);
                }

                // 加一条管理员角色关系记录 如果要初始化多个含超管的角色,仅创建超管角色与账号的关系
                if ($adminId > 0) {
                    if ($hasSuperAdmin) {
                        if ($value['mark'] === RoleModel::ROLE['super_admin']['mark']) {
                            AdminRoleRelationModel::query()->create([
                                'admin_id'      => $adminId,
                                'role_id'       => $role->getAttribute('id'),
                                'enterprise_id' => $this->enterpriseId
                            ]);
                        }
                    } else {
                        AdminRoleRelationModel::query()->create([
                            'admin_id'      => $adminId,
                            'role_id'       => $role->getAttribute('id'),
                            'enterprise_id' => $this->enterpriseId
                        ]);
                    }
                }
                $roles[$value['mark']] = $role->toArray();
            }

            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();

            Log::error('初始化内置角色异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            throw new \Exception("初始化内置角色异常: {$exception->getMessage()}");
        }

        return $roles;
    }

    /**
     * 新增角色
     *
     * @param array $roleData 角色信息
     * @param array $permissionsIdList 权限ID列表
     * @param int $adminId 管理员ID
     * @param int $rolePid 父级角色ID
     * @return array
     * @throws \Exception
     */
    public function add($roleData, $permissionsIdList, $adminId, $rolePid = RoleModel::SUPER_ADMIN_ROLE_PID)
    {
        $validPermissionsIdList = $this->getValidPermissionsIdList($rolePid, $permissionsIdList);

        if ($validPermissionsIdList === []) {
            throw new \Exception('角色无任何权限, 不能新增');
        }

        // 添加普通管理员名称不能重复
        if ($this->checkNameExists($roleData['name'])) {
            throw new \Exception('角色已存在');
        }
        $validPermissionsIdList = array_unique($validPermissionsIdList);
        sort($validPermissionsIdList);

        // 新增角色
        $time = time();
        try {
            DB::beginTransaction();
            $role = RoleModel::query()->create([
                'pid'           => $rolePid,
                'admin_id'      => $adminId,
                'mark'          => '',
                'name'          => $roleData['name'],
                'weight'        => $roleData['weight'] ?? 0,
                'enterprise_id' => $this->enterpriseId
            ]);

            $rolePermissions = [];
            foreach ($validPermissionsIdList as $id) {
                $rolePermissions[] = [
                    'role_id'        => $role->getAttribute('id'),
                    'permissions_id' => $id,
                    'enterprise_id'  => $this->enterpriseId,
                    'created_at'     => $time,
                    'updated_at'     => $time
                ];
            }

            // 更新角色权限信息
            foreach (array_chunk($rolePermissions, 100) as $rolePermission) {
                RolePermissionsModel::query()->insert($rolePermission);
            }

            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();

            Log::error('添加角色异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            throw new \Exception("新增角色异常: {$exception->getMessage()}");
        }
        $this->operationRole = $role;

        return $this->operationRole->toArray();
    }

    /**
     * 编辑角色
     *
     * @param int $roleId 角色ID
     * @param array $roleData
     * @param $permissionsIdList
     * @param array $filterPermissionsIdList
     * @return array
     * @throws \Exception
     */
    public function edit($roleId, $roleData, $permissionsIdList, $filterPermissionsIdList = [])
    {
        Log::info('编辑角色', [
            'args'       => func_get_args(),
            'enterprise' => $this->enterprise->toArray(),
            'version'    => $this->permissions->getVersion()
        ]);
        $this->codeTimer->open();

        $this->permissions->setRoleIdList([$roleId]);
        $role = RoleModel::query()->where('id', '=', $roleId)->first();
        if ($role === null) {
            throw new \Exception('角色不存在', Constant::NOT_FOUND_CODE);
        }

        // 编辑超管角色
        if ($role->getAttribute('mark') === RoleModel::ROLE['super_admin']['mark']) {
            $this->operationRole = $role;

            // 找出当前角色拥有的version_tag列表
            $hasVersionTagList      = $this->permissions->versionTagList();
            $validPermissionsIdList = $this->buildSuperAdminPermissionsIdList($permissionsIdList, $hasVersionTagList, 'edit');

            // 查出所有的可选功能，在前端提交之外的都删除
            $additionalConfigList = $this->permissions->additionalConfigList();
            $permissionsIdDiff    = array_diff(array_column($additionalConfigList, 'id'), $permissionsIdList);

            // 过滤依赖权限
            $filterPermissionsIdList = array_merge($filterPermissionsIdList, $this->filterDependList($validPermissionsIdList));

            // 清理该企业下所有角色不应该有的可选功能
            if ($permissionsIdDiff !== []) {

                $all = $this->permissions->all();

                $clearPermissionsIdList = $permissionsIdDiff;
                foreach ($permissionsIdDiff as $id) {
                    $clearPermissionsIdList = array_merge($clearPermissionsIdList, array_column($this->getChildren($all, $id), 'id'));
                }
            }
        } else {
            // 编辑普通管理员名称不能重复
            if ($this->checkNameExists($roleData['name'], $roleId)) {
                throw new \Exception('角色已存在');
            }

            $this->checkAuth($roleId);
            $validPermissionsIdList = $this->getValidPermissionsIdList($this->operationRole->getAttribute('pid'), $permissionsIdList);
        }

        // 过滤特殊权限
        $validPermissionsIdList = $this->filterPermissions($validPermissionsIdList, $filterPermissionsIdList);

        // 普通角色必须要有权限
        if ($validPermissionsIdList === [] && in_array($role->getAttribute('mark'), [
                RoleModel::ROLE['area']['mark'],
                RoleModel::ROLE['store_member']['mark']
            ]) === false) {
            throw new \Exception('角色无任何权限, 不能编辑');
        }
        $validPermissionsIdList = array_unique($validPermissionsIdList);
        sort($validPermissionsIdList);

        // 更新角色信息
        $time            = time();
        $rolePermissions = [];
        foreach ($validPermissionsIdList as $id) {
            $rolePermissions[] = [
                'role_id'        => $this->operationRole->getAttribute('id'),
                'permissions_id' => $id,
                'enterprise_id'  => $this->enterpriseId,
                'created_at'     => $time,
                'updated_at'     => $time
            ];
        }

        try {
            DB::beginTransaction();

            $this->operationRole->update([
                'name'   => $roleData['name'] ?? $this->operationRole->getAttribute('name'),
                'weight' => $roleData['weight'] ?? $this->operationRole->getAttribute('weight'),
            ]);

            // 更新角色权限信息
            $oldPermissionIds = RolePermissionsModel::query()
                ->where('role_id', '=', $this->operationRole->getAttribute('id'))
                ->pluck('permissions_id')
                ->toArray();
            $newPermissionIds = array_column($rolePermissions, 'permissions_id');
            //需新增的权限
            $addPerMissionIds = array_diff($newPermissionIds, $oldPermissionIds);
            //需删除的权限
            $delPerMissionIds = array_diff($oldPermissionIds, $newPermissionIds);
            // 删除需要删除的权限
            RolePermissionsModel::query()
                ->where('role_id', '=', $this->operationRole->getAttribute('id'))
                ->whereIn('permissions_id', $delPerMissionIds)
                ->forceDelete();
            //获取需要新增的权限
            $addPerMissions = [];
            foreach ($rolePermissions as $permissions) {
                if (in_array($permissions['permissions_id'], $addPerMissionIds)) {
                    $addPerMissions[] = $permissions;
                }
            }

            foreach (array_chunk($addPerMissions, 100) as $rolePermission) {
                RolePermissionsModel::query()->insert($rolePermission);
            }

            // boss修改了选配功能 删掉此功能及下级
            if (isset($clearPermissionsIdList) && $clearPermissionsIdList !== []) {
                $this->permissions->deleteByPermissionsId($clearPermissionsIdList);
            }

            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();

            Log::error('编辑角色异常', [
                'message' => $exception->getMessage(),
                'file'    => $exception->getFile(),
                'line'    => $exception->getLine(),
                'args'    => func_get_args()
            ]);

            throw new \Exception("编辑角色异常: {$exception->getMessage()}");
        }

        $this->codeTimer->close();
        Log::info('编辑角色完成', [
            'runtime' => $this->codeTimer->getRunTime()
        ]);

        $pattern = sprintf(CacheKey::ROLE, $this->enterpriseId, '*');
        $keys    = $this->redis->keys($pattern);
        if ($keys !== []) {
            $this->redis->del($keys);
        }

        return $this->operationRole->toArray();
    }

    /**
     * 删除角色
     *
     * @param int $roleId 角色ID
     * @return bool|mixed|null
     * @throws \Exception
     */
    public function delete($roleId)
    {
        $this->checkAuth($roleId);

        if ($this->operationRole->getAttribute('mark') === RoleModel::ROLE['super_admin']['name']) {
            throw new \Exception('超管不能被删除');
        }

        $relationIdList = AdminRoleRelationModel::query()->where('role_id', '=', $roleId)->pluck('id')->toArray();
        if ($relationIdList !== []) {
            throw new \Exception('该角色下存在管理员, 不能删除');
        }

        $this->operationRole->delete();
        $this->permissions->deleteByRoleId($roleId);
    }

    /**
     * 角色列表
     *
     * @param int $page
     * @param int $size
     * @param array $input
     * @return array
     */
    public function lists($page, $size = 10, $input = [])
    {
        /*$allRole = RoleModel::query()->where('enterprise_id', $this->enterpriseId)->get()->toArray();

        // 找出当前角色可见的下级角色
        $lookRoleIdList = [];
        foreach ($this->roleIdList as $roleId) {
            $lookRoleIdList = array_merge($lookRoleIdList, array_column($this->getChildren($allRole, $roleId), 'id'));
        }*/
        $orm = RoleModel::query()
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->where('mark', '!=', RoleModel::ROLE['super_admin']['mark'])
            ->when(!empty($input['role_name']), function (
                $query
            ) use ($input) {
                return $query->where('name', 'like', "%{$input['role_name']}%");
            })
            ->when(!empty($input['add_account']), function ($query) use ($input) {
//                $excludeRoleMark = ['guide_enterprise_guide','guide_enterprise_admin','console_enterprise_admin'];
                $excludeRoleMark = ['console_enterprise_admin'];

                return $query->whereNotIn('mark', $excludeRoleMark);
            })
            ->when(!empty($input['operator']), function ($query) use ($input) {
                return $query->whereIn('admin_id', $input['operator']);
            });

        $total = $orm->count('id');

        $list  = $orm->withCount(['adminRelation as admin_relation_count' => function ($query) {
            $query->leftJoin('enterprise_admin', 'enterprise_admin.id', 'admin_role_relation.admin_id')->where("enterprise_admin.deleted_at", 0);
        }])
            ->orderBy('admin_relation_count', $input['sort_type'])
            ->orderByDesc("created_at")
            ->orderBy('weight', 'DESC')
            ->orderBy('id', 'DESC')
            ->forPage($page, $size)
            ->get()
            ->toArray();

        return [$total, $list];
    }

    /**
     * 查询当前企业超管角色
     *
     * @return RoleModel|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function getSuperAdminRole()
    {
        return RoleModel::query()
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->where('mark', '=', RoleModel::ROLE['super_admin']['mark'])
            ->first();
    }

    /**
     * 检查角色名称是否存在
     *
     * @param string $roleName 角色名称
     * @param int $filterRoleId 过滤的角色ID
     * @return bool
     */
    public function checkNameExists($roleName, $filterRoleId = 0)
    {
        return RoleModel::query()
            ->where('enterprise_id', '=', $this->enterpriseId)
            ->where('name', '=', $roleName)
            ->when($filterRoleId > 0, function ($query) use ($filterRoleId) {
                return $query->where('id', '!=', $filterRoleId);
            })
            ->exists();
    }

    /**
     * 检查操作权限
     *
     * @param $roleIdOrMark
     * @throws \Exception
     */
    private function checkAuth($roleIdOrMark)
    {
        // 查一下被操作的角色是否存在
        $role = RoleModel::query()->when(is_int($roleIdOrMark), function ($query) use ($roleIdOrMark) {
            return $query->where('id', '=', $roleIdOrMark);
        })->when(is_string($roleIdOrMark), function ($query) use ($roleIdOrMark) {
            return $query->where('mark', '=', $roleIdOrMark)->where('enterprise_id', '=', $this->enterpriseId);
        })->first();
        if ($role === null) {
            throw new \Exception('角色不存在', Constant::NOT_FOUND_CODE);
        }

        /* v3.0.30 因为现在角色没有层级先去掉判断
        // 被操作角色如果不是超管则继续验证
        if ($role->getAttribute('mark') !== RoleModel::ROLE['super_admin']['mark']) {
            $roleList = RoleModel::query()->where('enterprise_id', '=', $this->enterpriseId)->get()->toArray();

            // 多角色遍历一下 看被操作角色是否是当前角色的下级
            foreach ($this->roleIdList as $key => $value) {
                $children = $this->getChildren($roleList, $value);
                if (in_array($role->getAttribute('id'), array_column($children, 'id')) === false) {
                    throw new \Exception('该角色不是操作角色下级, 不能操作');
                }
            }
        }*/

        $this->operationRole = $role;
    }

    /**
     * 获取有效权限
     *
     * @param int $rolePid 父级角色ID
     * @param array $permissionsIdList 权限ID列表
     * @return array
     */
    private function getValidPermissionsIdList($rolePid, $permissionsIdList)
    {
        return $permissionsIdList;
        $validPermissionsIdList = $permissionsIdList;

        // 如果不是超管角色则验证权限 传进来的权限列表是否在当前操作角色拥有的权限范围内
        if (in_array(0, $this->roleIdList) === false || $rolePid !== RoleModel::SUPER_ADMIN_ROLE_PID) {
            // 查出当前操作人角色对应的权限ID列表
            $roleListPermissionsIdList = array_column($this->permissions->lists(), 'id');

            // 有效权限ID 取交集过滤掉父级以外的权限
            $validPermissionsIdList = array_intersect($roleListPermissionsIdList, $permissionsIdList);
        }

        return $validPermissionsIdList;
    }

    /**
     * 生成超管的权限ID列表
     *
     * @param $selectPermissionsIdList
     * @param array $hasVersionTagList
     * @param string $operation
     * @return array
     */
    private function buildSuperAdminPermissionsIdList($selectPermissionsIdList, $hasVersionTagList = [], $operation = 'init')
    {
        Log::info('生成超管的权限ID列表', [
            'args'    => $selectPermissionsIdList,
            'version' => $this->permissions->getVersion()
        ]);

        // 如果是init操作说明初始化超管，此时不需要没有全量的权限， 反之为编辑
        $standardConfig = $this->permissions->getStandardConfig($hasVersionTagList, $operation === 'init');

        $all = $this->permissions->all();

        // 找出勾选的可选配置所有上下级
        $parent   = [];
        $children = [];
        foreach ($selectPermissionsIdList as $id) {
            $parent   = array_merge($parent, $this->getParentFilterNoConfig($all, $id));
            $children = array_merge($children, $this->getChildren($all, $id));
        }
        // 选配权限跟标配权限合并
        $selectIdList = array_column(array_merge($parent, $children), 'id');

        // 放入选配本身
        $selectIdList         = array_merge($selectPermissionsIdList, $selectIdList);
        $standardConfigIdList = array_column($standardConfig, 'id');
        $allIdList            = array_merge($standardConfigIdList, $selectIdList);

        // 添加一个内置角色
        (new RoleStrategy())->getRoleStrategy($this->enterpriseId, $this->adminId, 'guide_enterprise_admin');

        return array_values(array_unique($allIdList));
    }

    /**
     * 过滤的依赖权限ID列表
     *
     * @param array $validPermissionsIdList
     * @return array
     */
    private function filterDependList($validPermissionsIdList)
    {
        $filterPermissionsIdList = [];
        foreach ($this->permissions->dependRelation() as $dependId => $functionIds) {
            if ($functionIds === []) {
                continue;
            }
            // 判断依赖功能是否在当前要添加的权限,不在的话直接移除所有的被依赖权限
            if (in_array($dependId, $validPermissionsIdList) === false) {
                $filterPermissionsIdList = array_merge($filterPermissionsIdList, $functionIds);
                continue;
            }
        }

        return $filterPermissionsIdList;
    }

    /**
     * 过滤掉指定权限
     *
     * @param array $permissionsIdList 有效权限
     * @param array $filterPermissionsIdList 过滤权限
     * @return array
     */
    private function filterPermissions($permissionsIdList, $filterPermissionsIdList)
    {
        if ($filterPermissionsIdList !== []) {
            $permissionsIdList = array_filter($permissionsIdList, function ($val) use ($filterPermissionsIdList) {
                // 返回true则不过滤
                if (in_array($val, $filterPermissionsIdList)) {
                    return false;
                }
                return true;
            });
            $permissionsIdList = array_values($permissionsIdList);
        }

        return $permissionsIdList;
    }

    private function getParentFilterNoConfig($data, $pid)
    {
        static $parent = [];

        foreach ($data as $key => $value) {
            if ($value['id'] === $pid && $value[$this->permissions->getVersionField()] !== PermissionsModel::VERSION_NO_CONFIG) {
                $parent[] = $value;

                $this->getParent($data, $value['pid']);
            }
        }

        return $parent;
    }
}
