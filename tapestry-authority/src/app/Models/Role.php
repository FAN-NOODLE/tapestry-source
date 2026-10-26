<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Collection;

class Role extends Base
{
    /**
     * 超管角色ID
     */
    const SUPER_ADMIN_ROLE_PID = 0;

    /**
     * 内置角色
     */
    const ROLE = [
        'super_admin'   => [
            'name'                => '超管',
            'mark'                => 'super_admin',
            'permissions_id_list' => [],
            'key_list'            => []
        ],
        'area'          => [
            'name'                => '区域',
            'mark'                => 'area',
            'permissions_id_list' => [],
            'key_list'            => []
        ],
        'store_manager' => [
            'name'                => '店长',
            'mark'                => 'store_manager',
            'permissions_id_list' => [
                111,
                115,
                116,
                112,
                114,
                113,
                121,
                125,
                128,
                133,
                134,
                130,
                129,
                131,
                132,
                127,
                126,
                117,
                118,
                119,
                120,
                122,
                123,
                124
            ],
            'key_list'            => [
                'myapp_store_store-manager',
                'myapp_store_store-manager_goods-decorate-edit',
                'myapp_store_store-manager_commodity-management',
                'myapp_store_store-manager_store-detail',
                'myapp_store_store-manager_store-decorate',
                'myapp_store_store-manager_performance-detail',
                'myapp_store_store-manager_performance-list',
                'myapp_store_store-manager_clerk-list',
                'myapp_store_store-manager_store-edit',
                'myapp_store_store-manager_store-list',
                'myapp_store',
                'myapp_order-manager_order',
                'myapp_order-manager_order_detail',
                'myapp_order-manager_order_list',
                'myapp_order-manager',
                'myapp_order-manager_after-sale-manage',
                'myapp_order-manager_after-sale-manage_list'
            ]
        ],
        'store_member'  => [
            'name'                => '店员',
            'mark'                => 'store_member',
            'permissions_id_list' => [],
            'key_list'            => []
        ],
    ];

    /**
     * 表名称:角色表
     *
     * @var string
     */
    protected $table = 'role';

    protected $guarded = [];

    /**
     * ID集合中的内置角色
     * @param $ids
     * @return Role[]|Builder[]|Collection
     */
    public static function markRoles($ids)
    {
        return self::query()
            ->whereIn('id', $ids)
            ->where('mark', '!=', '')
            ->get();
    }

    /**
     * 通过标识获取角色IDs
     * @param $marks
     * @param $eId
     * @return array
     */
    public static function roleByMark($marks, $eId)
    {
        return self::query()
            ->where('enterprise_id', '=', $eId)
            ->whereIn('mark', $marks)
            ->pluck('id')
            ->toArray();
    }

    /**
     * 角色互斥
     * @param $roleIds
     * @return bool
     */
    public static function roleMutexLimit($roleIds)
    {
        $roles = Role::markRoles($roleIds);

        // 没有内置角色
        if (empty($roles)) {
            return true;
        }

        // 多个内置角色
        if (count($roles) > 1) {
            error(702025);
        }

        foreach ($roles as $role) {
            if ($role['mark'] !== 'guide_enterprise_admin') {
                error(702026);
            }
        }

        return true;
    }

    public function adminRelation()
    {
        return $this->hasMany(AdminRoleRelation::class, 'role_id', 'id');
    }
}
