<?php

namespace App\Modules\BuiltInRole\BuiltInRoleActive;

use App\Models\Role;
use App\Models\Enterprise;
use App\Models\EnterpriseAdmin;
use App\Models\RolePermissions;
use App\Models\AdminRoleRelation;
use App\Models\EnterpriseDataRange;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Builder;
use App\Modules\BuiltInRole\RoleStrategyConfig;
use App\Modules\BuiltInRole\Interfaces\RoleInterface;

class CommonRole implements RoleInterface
{
    /*
     * 1. 社群导购和企业导购的数据范围仅为商城
     * 2. 店长店员的数据范围为区域或门店
     * */

    // 添加
    const ONLINE = 1;

    // 移除
    const OFFLINE = 2;

    // 配置数据
    public $data;

    // 角色
    public $role;

    // 添加当前人的管理员
    protected $adminId;

    // 互斥角色
    protected $mutexRole = [];

    /**
     * @var Enterprise|null
     */
    protected $enterprise;

    public function __construct($enterpriseId, $adminId, $roleData, $mode = 'script')
    {
        // 验证企业是否存在
        $this->enterprise = Enterprise::query()->where('id', '=', $enterpriseId)->first();
        if (empty($this->enterprise) && $mode == 'script') {
            return [];
        }
        if (empty($this->enterprise)) {
            error(701005);
        }

        $this->data = $roleData;
        $this->adminId = $adminId;

        // 验证角色是否存在
        $role = Role::query()->where('enterprise_id', '=', $enterpriseId)->where('mark', '=', $this->data['mark'])->first();

        $this->role = $role;
        if ($role === null) {
            $this->initRole();
        }
    }

    /**
     * 初始化角色
     */
    public function initRole()
    {
        // 初始化创建角色
        $this->role = Role::query()->create([
            'mark'          => $this->data['mark'],
            'name'          => $this->data['name'],
            'admin_id'      => $this->adminId,
            'enterprise_id' => $this->enterprise->id,
        ]);

        // 初始化角色权限关联
        $rolePermissionRelation = [];
        foreach ($this->data['permission'] as $pId) {
            $rolePermissionRelation[] = [
                'role_id'        => $this->role->id,
                'enterprise_id'  => $this->enterprise->id,
                'permissions_id' => $pId,
            ];
        }

        // 插入关联数据
        if ($rolePermissionRelation) {
            RolePermissions::query()->insert($rolePermissionRelation);
        }
    }

    /**
     * 获取角色
     *
     * @param $adminId
     * @param $roleId
     * @return AdminRoleRelation|Builder|Model|object|null
     */
    public function getRole($adminId, $roleId)
    {
        return AdminRoleRelation::query()->where('admin_id', '=', $adminId)->where('role_id', '=', $roleId)->first();
    }

    /**
     * 管理员绑定角色
     *
     * @param $adminId
     */
    public function adminBindRole($adminId)
    {
        $role = $this->getRole($adminId, $this->role->id);
        if (empty($role)) {
            AdminRoleRelation::query()->create([
                'role_id'       => $this->role->id,
                'admin_id'      => $adminId,
                'enterprise_id' => $this->enterprise->id,
            ]);
        }
    }

    /**
     * 解绑角色
     *
     * @param $adminId
     */
    public function unsetRoleRelation($adminId)
    {
        $role = $this->getRole($adminId, $this->role->id);
        if ($role) {
            $role->forceDelete();
        }
    }

    /**
     * 获取数据范围
     *
     * @param $adminId
     * @param $type
     * @param $objectId
     * @return EnterpriseDataRange|Builder|Model|object|null
     */
    public function getDataRange($adminId, $type, $objectId)
    {
        $query = EnterpriseDataRange::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('enterprise_admin_id', '=', $adminId)
            ->where('type', '=', $type);

        // 门店有值, 商城仅有一个值
        if ($objectId != 0) {
            $query->where('object_id', '=', $objectId);
        }

        return $query->first();
    }

    /**
     * 数据范围绑定
     *
     * @param $adminId
     * @param $type
     * @param int $objectId
     */
    public function adminBindDataRange($adminId, $type, $objectId = 0)
    {
        $range = $this->getDataRange($adminId, $type, $objectId);
        if (empty($range)) {
            EnterpriseDataRange::query()->create([
                'type'                => $type,
                'object_id'           => $objectId,
                'enterprise_id'       => $this->enterprise->id,
                'enterprise_admin_id' => $adminId,
            ]);
        }
    }

    /**
     * 数据范围解绑
     *
     * @param $adminId
     * @param $type
     * @param int $objectId
     */
    public function unsetDataRangeRelation($adminId, $type, $objectId = 0)
    {
        $range = $this->getDataRange($adminId, $type, $objectId);
        if ($range) {
            $range->forceDelete();
        }
    }

    /**
     * 账号操作 - 添加或更新
     *
     * @param $params
     * @param $defaultAdd
     */
    public function accountOperation($params, $defaultAdd = true)
    {
//        if ($defaultAdd) {
//            $this->accountLimitValidate();
//        }
    }

    /**
     * 账号数量限制
     */
    public function accountLimitValidate()
    {
        $adminCount = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('admin_available_system', '!=', json_encode([2]))
            ->count();
        if ($adminCount >= $this->enterprise->account_limit) {
            error(701004, [$this->enterprise->account_limit]);
        }
    }

    public function adminRoleMutex($admin)
    {
        // 获取当前已有角色
        $roles = AdminRoleRelation::adminRoles($admin->id, $admin->enterprise_id);

        // 判断已有角色是否在当前角色的互斥角色中不
        $marks = array_unique(array_values(array_column($roles, 'mark')));

        // 获取已有角色在互斥中的标识
        $existsMutex = array_intersect($marks, $this->mutexRole);

        // 互斥抛异常
        if (array_values($existsMutex)) {
            error(702024);
        }
    }

    /**
     * 数据范围绑定
     *
     * @param $admin EnterpriseAdmin|Builder|Model|object|null
     */
    public function adminBindAvailableSystem($admin)
    {
        $system = json_decode($admin->admin_available_system, true);
        $fullSystem = array_merge($system, $this->data['role_available_system']);
        $uniqueSystem = array_sort(array_unique(array_values($fullSystem)));
        $admin->update([
            'admin_available_system' => json_encode(array_values($uniqueSystem))
        ]);
    }

    /**
     * 数据范围解绑
     *
     * @param $admin EnterpriseAdmin|Builder|Model|object|null
     * @param array $mustSystem
     */
    public function unsetAvailableSystem($admin, $mustSystem = [])
    {
        // 获取当前绑定的所有角色
        $adminAllRole = AdminRoleRelation::query()
            ->select(['role.id', 'role.name', 'role.mark', 'admin_role_relation.admin_id'])
            ->leftJoin('role', 'admin_role_relation.role_id', 'role.id')
            ->where('admin_role_relation.admin_id', '=', $admin->id)
            ->where('admin_role_relation.deleted_at', '=', 0)
            ->get()
            ->toArray();

        $hasMark = [];
        $notHasMark = [];
        // 区分开有标识和没标识, 并且去除当前绑定的角色
        foreach ($adminAllRole as $role) {
            if ($role['mark'] === $this->data['mark']) {
                continue;
            }

            if (empty($role['mark'])) {
                // 无标识的角色
                $notHasMark[] = $role;
            } else {
                // 获取有标识的角色
                $hasMark[] = $role['mark'];
            }
        }

        // 账号存在无标识的角色, 说明是超管建立的角色, 必须给管理后台的可用系统
        $systems = [];
        if ($notHasMark) {
            $systems = [1];
        }

        // 获取当前账号拥有的所有角色(除了当前角色), 对应的可用系统
        if ($hasMark) {
            foreach (RoleStrategyConfig::ROLE as $everyRole) {
                if (in_array($everyRole['mark'], $hasMark)) {
                    $systems = array_merge($systems, $everyRole['role_available_system']);
                }
            }
        }

        // 像店长, 一个账号可以是多门店店长, 去除当前门店, 可能还在其他门店为店长, 在调用外部进行判断
        if ($mustSystem) {
            $systems = array_merge($systems, $mustSystem);
        }

        // 实际用户可用系统
        $systems = array_sort(array_unique($systems));

        $admin->update([
            'admin_available_system' => json_encode(array_values($systems))
        ]);
    }

    public function adminStatusEnable($admin)
    {
        // 管理员状态启用
        $admin->update([
            'status' => 1
        ]);
    }

    public function getPermission()
    {
    }
}