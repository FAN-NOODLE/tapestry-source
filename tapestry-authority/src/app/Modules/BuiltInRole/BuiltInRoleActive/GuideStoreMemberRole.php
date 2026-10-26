<?php

namespace App\Modules\BuiltInRole\BuiltInRoleActive;

use CrsSdk\Shop\StoreWorker;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminService;

class GuideStoreMemberRole extends CommonRole
{
    /**
     * 当前角色互斥的角色
     */
    protected $mutexRole = [
        'guide_enterprise_guide',
        'console_enterprise_guide',
        'guide_community_guide',
        'console_community_guide',
        'guide_enterprise_admin',
    ];

    /**
     * 店员不算在账号数量限制上
     *
     * @param $params
     * @param $defaultAdd
     * @throws \Exception
     */
    public function accountOperation($params, $defaultAdd = true)
    {
        $account         = (new EnterpriseAdminService())->createDefaultAccount($params['username'], $params['mobile']);

        // 获取当前的店员账号
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('account', 'like', "%{$account}")
            ->first();

        // 如果是添加店员, 创建账号并给与角色和数据权限
        if ($params['type'] === self::ONLINE) {
            if (empty($admin)) {
                $salt            = (new EnterpriseAdminService())->makeSalt();

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
                // 如果账号已存在, 检查是否角色互斥
                $this->adminRoleMutex($admin);
            }

            // 绑定角色
            $this->adminBindRole($admin->id);
            // 创建店员的数据权限
            $this->adminBindDataRange($admin->id, EnterpriseDataRange::TYPE_STORE, $params['store_id']);
            // 数据范围矫正
            $this->adminBindAvailableSystem($admin);
            // 账号启用
            $this->adminStatusEnable($admin);
        }

        // 离职
        $existsRange = null;
        if ($params['type'] === self::OFFLINE) {

            if (empty($admin)) {
                error(700024);
            }
            // 检查是否还存在其他店员或店长身份, 没有身份, 进行角色解绑, 然后去除可用系统
            $workers = (new StoreWorker())->lists([
                'is_page'       => false,
                'mobile'        => $params['mobile'],
                'enterprise_id' => $this->enterprise->id,
                'worker_status' => 1,
                'type'          => 1
            ])['list'];

            // 如果原店长的身份还存在, 去除掉, 拿剩余的身份
            $hasStoreMemberIdentity = false;
            $storeIds = array_column($workers, 'store_id');
            if (in_array($params['store_id'], $storeIds)) {
                foreach ($workers as $key => $storeWorker) {
                    if ($storeWorker['store_id'] == $params['store_id']) {
                        unset($workers[$key]);
                        continue;
                    }
                    // 存在其他门店店员
                    if ($storeWorker['is_boss'] != 1) {
                        $hasStoreMemberIdentity = true;
                    }
                }
            }

            // 没有其他身份了,解绑角色, 如果有, 最低也是门店店员, 所以不用解绑啥
            if (count($workers) == 0) {
                $this->unsetAvailableSystem($admin);
            }

            // 判断是否还有其他门店的店员 不存在才进行解除
            if ($hasStoreMemberIdentity === false) {
                $this->unsetRoleRelation($admin->id);
            }

            // 清除店员当前门店的权限
            $this->unsetDataRangeRelation($admin->id, EnterpriseDataRange::TYPE_STORE, $params['store_id']);
        }
    }
}