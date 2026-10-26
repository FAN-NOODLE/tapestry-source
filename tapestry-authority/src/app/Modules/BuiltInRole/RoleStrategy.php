<?php

namespace App\Modules\BuiltInRole;

use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Models\Role;
use App\Models\Permissions;
use App\Models\RolePermissions;
use App\Models\AdminRoleRelation;
use App\Modules\BuiltInRole\BuiltInRoleActive\CommonRole;
use CrsSdk\Shop\StoreWorker;

class RoleStrategy
{
    /**
     * @var CommonRole
     */
    public $roleStrategy = null;

    public $prefix = 'App\Modules\BuiltInRole\BuiltInRoleActive\\';

    public function getRoleStrategy($enterpriseId, $adminId, $mark = '', $model = 'api')
    {
        if (empty($mark)) {
            error(702023);
        }

        $roleData = RoleStrategyConfig::getConfig($mark);

        if (class_exists($this->prefix . $roleData['namespace'])) {
            $fullNamespace = $this->prefix . $roleData['namespace'];
        } else {
            $fullNamespace = $this->prefix . 'CommonRole';
        }

        $this->roleStrategy = (new $fullNamespace($enterpriseId, $adminId, $roleData, $model));

        return $this;
    }

    /**
     * 初始化企业内置角色
     *
     * @param array $enterpriseIds
     * @param string $mode
     */
    public static function initFullRoles($enterpriseIds = [], $mode = 'script')
    {
        $prefix = 'App\Modules\BuiltInRole\BuiltInRoleActive\\';

        $query = EnterpriseAdmin::query()->select(['id', 'enterprise_id'])->where('is_super', '=', 1);

        if ($enterpriseIds) {
            $query = $query->whereIn('enterprise_id', $enterpriseIds);
        }

        $admins = $query->get()->toArray();

        foreach (RoleStrategyConfig::ROLE as $item) {

            if (class_exists($prefix . $item['namespace'])) {
                $fullNamespace = $prefix . $item['namespace'];
            } else {
                $fullNamespace = $prefix . 'CommonRole';
            }

            foreach ($admins as $admin) {
                echo '正在处理企业' . $admin['enterprise_id'] . '的内置角色' . PHP_EOL;
                (new $fullNamespace($admin['enterprise_id'], $admin['id'], $item, $mode));
            }
        }
    }

    /**
     * 门店数据处理 员工角色绑定
     *
     * @param array $ids
     * @param string $mode
     * @throws \Exception
     */
    public static function storeRoleRelation($ids = [], $mode = 'script')
    {

        $es = Enterprise::query()->whereIn('id', $ids)->get();

        foreach ($es as $e) {

            echo '正在处理企业' . $e->id . '的门店员工角色绑定' . PHP_EOL;

            // 获取超管数据
            $adminSuper = EnterpriseAdmin::query()->where('enterprise_id', '=', $e->id)->where('is_super', '=', 1)->first();
            if (empty($adminSuper)) {
                $adminSuper = [];
            } else {
                $adminSuper = $adminSuper->toArray();
            }

            $stores = (new StoreWorker())->lists([
                'enterprise_id' => $e->id,
                'type'          => 1,
                'is_page'       => false,
            ]);

            foreach ($stores['list'] as $store) {
                try {
                    // 店长
                    if ($store['is_boss'] === 1) {
                        $strategy = (new RoleStrategy())->getRoleStrategy($e->id, $adminSuper['id'] ?? 0, 'guide_store_manager', $mode);
                        $strategy->roleStrategy->accountOperation([
                            'type'          => 1,
                            'mobile'        => $store['mobile'],
                            'store_id'      => $store['store_id'],
                            'username'      => $store['name'],
                            'operator_id'   => $adminSuper['id'] ?? 0,
                            'enterprise_id' => $e->id,
                        ]);

                        $strategy = (new RoleStrategy())->getRoleStrategy($e->id, $adminSuper['id'] ?? 0, 'console_store_manager', $mode);
                        $strategy->roleStrategy->accountOperation([
                            'type'          => 1,
                            'mobile'        => $store['mobile'],
                            'store_id'      => $store['store_id'],
                            'username'      => $store['name'],
                            'operator_id'   => $adminSuper['id'] ?? 0,
                            'enterprise_id' => $e->id,
                        ]);
                    }

                    // 店员
                    if ($store['is_boss'] === 2) {
                        $strategy = (new RoleStrategy())->getRoleStrategy($e->id, $adminSuper['id'] ?? 0, 'guide_store_member', $mode);
                        $strategy->roleStrategy->accountOperation([
                            'type'          => 1,
                            'mobile'        => $store['mobile'],
                            'store_id'      => $store['store_id'],
                            'username'      => $store['name'],
                            'operator_id'   => $adminSuper['id'] ?? 0,
                            'enterprise_id' => $e->id,
                        ]);
                    }
                } catch (\Exception $e) {

                }
            }
        }
    }

    /**
     * 获取管理员智慧导购的权限
     *
     * @param $adminId
     * @param $type
     * @param array $roleMark
     * @return array
     */
    public static function getPermission($adminId, $type, $roleMark = [])
    {
        if (!in_array($type, [Permissions::SYSTEM_TYPE_GUIDE])) {
            error(700025);
        }

        $admin = EnterpriseAdmin::adminById($adminId);
        if (empty($admin)) {
            error(700024);
        }

        if ($roleMark) {
            $roleIds = Role::roleByMark($roleMark, $admin->enterprise_id);
        } else {
            $roleIds = AdminRoleRelation::roleIdsByAdminId($adminId);
        }

        $permissions = RolePermissions::pIdsByRoleIds($roleIds);

        return Permissions::filterPermissionId($permissions, $type);
    }

    /**
     * 获取管理员的角色
     *
     * @param $adminId
     * @param int $systemType
     * @return array
     */
    public static function getRole($adminId, $systemType = [0])
    {
        $roles = AdminRoleRelation::query()->where('admin_id', '=', $adminId)->pluck('role_id')->toArray();

        $roleIds = RolePermissions::query()
            ->leftJoin('permissions', 'role_permissions.permissions_id', 'permissions.id')
            ->whereIn('role_permissions.role_id', $roles)
            ->whereIn('permissions.system_type', $systemType)
            ->pluck('role_id')
            ->toArray();

        return Role::query()->select(['name as role_name', 'mark as role_mark'])->whereIn('id', $roleIds)->get()->toArray();
    }
}