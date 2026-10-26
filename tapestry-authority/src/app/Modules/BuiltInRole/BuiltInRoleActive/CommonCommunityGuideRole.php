<?php

namespace App\Modules\BuiltInRole\BuiltInRoleActive;

use App\Common\Constant;
use App\Models\EnterpriseAdmin;
use App\Models\AdminRoleRelation;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;

class CommonCommunityGuideRole extends CommonRole
{
    /**
     * 当前角色互斥的角色
     */
    protected $mutexRole = [
        'guide_store_member',
        'guide_store_manager',
        'console_store_manager',
        'guide_enterprise_guide',
        'console_enterprise_guide',
        'guide_enterprise_admin',
    ];

    public function accountOperation($params, $defaultAdd = true)
    {
        if ($params['type'] === parent::OFFLINE) {
            $defaultAdd = false;
        }

        parent::accountOperation($params, $defaultAdd);

        $account = (new EnterpriseAdminService())->createDefaultAccount($params['username'], $params['mobile']);

        // 查找当前账号
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('account', 'like', "%{$account}")
            ->first();

        // 添加导购
        if ($params['type'] === parent::ONLINE) {
            // 不存在创建
            if (empty($admin)) {
                $salt = (new EnterpriseAdminService())->makeSalt();

                $admin = EnterpriseAdmin::query()->create([
                    'enterprise_id'          => $this->enterprise->id,
                    'username'               => $params['username'],
                    'account'                => $account,
                    'password'               => md5(md5($account) . $salt),
                    'salt'                   => $salt,
                    'global_roaming'         => '+86',
                    'admin_available_system' => json_encode($this->data['role_available_system']),
                    'mobile'                 => $params['mobile'],
                    'is_super'               => EnterpriseAdmin::NO_SUPER,
                    'range_type'             => EnterpriseAdmin::RANGE_TYPE_ONLY,
                    'operator_id'            => $params['operator_id']
                ]);
            } else {
                // 导购升级企业管理 去掉角色互斥校验
                if (!isset($params['noCheck'])){
                    // 如果账号已存在, 检查是否角色互斥
                    $this->adminRoleMutex($admin);
                }
            }

            // 重新绑定角色和数据范围 - 存在以前加过,后移除,账号在,角色和数据范围被清除
            $this->adminBindRole($admin->id);

            // 企业管理角色
            if (!empty($params['range'])) {
                // 插入账号范围数据
                $dataRange = [
                    'range'               => $params['range'],
                    'enterprise_admin_id' => $admin->id,
                    'enterprise_id'       => $admin->enterprise_id
                ];
                (new EnterpriseDataRangeService())->addEnterpriseDataRange($dataRange);
            } else {
                $this->adminBindDataRange($admin->id, EnterpriseDataRange::TYPE_SHOP);
            }
            $this->adminBindAvailableSystem($admin);
            $this->adminStatusEnable($admin);
        }

        // 移除导购
        if ($params['type'] === parent::OFFLINE) {
            // 判断账号存在不, 不存在抛异常
            if (empty($admin)) {
                error(700012);
            }

            // 解除角色和数据范围的绑定
            $this->unsetRoleRelation($admin->id);

            // 企业管理角色
            if ($this->data['mark'] == Constant::GUIDE_ENTERPRISE_ADMIN) {
                EnterpriseDataRange::query()->where('enterprise_admin_id', $admin->id)->delete();
            } else {
                $this->unsetDataRangeRelation($admin->id, EnterpriseDataRange::TYPE_SHOP);
            }
            $this->unsetAvailableSystem($admin);
        }
    }

    public function unsetDataRangeRelation($adminId, $type, $objectId = 0)
    {
        $range = $this->getDataRange($adminId, $type, $objectId);

        // 如果账号关联了同样拥有商城数据权限的角色, 不进行清除
        // 获取当前绑定的所有角色
        $adminAllRole = AdminRoleRelation::query()
            ->select(['role.id', 'role.name', 'role.mark', 'admin_role_relation.admin_id'])
            ->leftJoin('role', 'admin_role_relation.role_id', 'role.id')
            ->where('admin_role_relation.admin_id', '=', $adminId)
            ->whereIn('role.mark', ['guide_community_guide', 'guide_enterprise_guide'])
            ->get()
            ->toArray();

        foreach ($adminAllRole as $key => $item) {
            if ($item['mark'] === $this->data['mark']) {
                unset($adminAllRole[$key]);
            }
        }
        // 拥有相同的商城数据类型的角色
        if (count($adminAllRole) > 0) {
            return true;
        }

        if ($range) {
            $range->forceDelete();
        }
    }
}