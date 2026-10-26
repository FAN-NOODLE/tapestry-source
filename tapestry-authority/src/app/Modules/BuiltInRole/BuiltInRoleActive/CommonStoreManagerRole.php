<?php

namespace App\Modules\BuiltInRole\BuiltInRoleActive;

use CrsSdk\Shop\StoreWorker;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseDataRange;
use App\Services\EnterpriseAdminService;

class CommonStoreManagerRole extends CommonRole
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

    public function accountOperation($params, $defaultAdd = true)
    {
        parent::accountOperation($params, $defaultAdd);

        // 更换店长,去掉上一位店长的当前门店数据范围
        $beforeAdmin = null;
        if (!empty($params['before_mobile'])) {
            $beforeAccount = (new EnterpriseAdminService())->createDefaultAccount($params['username'], $params['before_mobile']);

            $beforeAdmin = EnterpriseAdmin::query()
                ->where('account', 'like', "%{$beforeAccount}")
                ->where('enterprise_id', '=', $this->enterprise->id)
                ->first();
        }
        // 存在上一位店长
        if ($beforeAdmin) {
            // 检查是否还存在其他电员或店长身份, 没有身份, 进行角色解绑, 然后去除可用系统
            $workers = (new StoreWorker())->lists([
                'enterprise_id' => $this->enterprise->id,
                'mobile'        => $params['before_mobile'],
                'worker_status' => 1,
                'is_page'       => false,
                'type'          => 1
            ])['list'];

            // 如果原店长的身份还存在, 去除掉, 拿剩余的身份
            $storeIds = array_column($workers, 'store_id');
            if (in_array($params['store_id'], $storeIds)) {
                foreach ($workers as $key => $storeWorker) {
                    if ($storeWorker['store_id'] == $params['store_id']) {
                        unset($workers[$key]);
                    }
                }
            }

            // 如果这个手机号有多个身份, 检查还有其他店长权限吗, 没有接去除店长相关权限
            $hasStoreMemberIdentity = false;
            $hasStoreManagerIdentity = false;
            if (count($workers) > 0) {
                foreach ($workers as $worker) {
                    if ($worker['is_boss'] === 1) {
                        $hasStoreManagerIdentity = true;
                    } else {
                        $hasStoreMemberIdentity = true;
                    }
                }
            } else {
                // 没有其他身份了, 解绑角色, 去除可用系统
                $this->unsetRoleRelation($beforeAdmin->id);
                $this->unsetAvailableSystem($beforeAdmin);
            }

            // 如果还有店长权限, 那就不解决角色, 如果只剩下门店权限, 那要解绑角色, 可用系统只有智慧导购小程序
            if ($hasStoreManagerIdentity === false && $hasStoreMemberIdentity === true) {
                $this->unsetRoleRelation($beforeAdmin->id);
                $this->unsetAvailableSystem($beforeAdmin, [2]);
            }

            // 去掉他当前门店的数据权限
            //            $this->unsetDataRangeRelation($beforeAdmin->id, EnterpriseDataRange::TYPE_STORE, $params['store_id']);
        }

        $account = (new EnterpriseAdminService())->createDefaultAccount($params['username'], $params['mobile']);

        // 查找当前店长账号
        $admin = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $this->enterprise->id)
            ->where('account', 'like', "%{$account}")
            ->first();

        $salt = (new EnterpriseAdminService())->makeSalt();
        if (empty($admin)) {
            // 不存在创建
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

            // 原先是其他店店员, 现在是店长, 要生成密码
            if (empty($admin->password)) {
                $admin->update(['password' => md5(md5('888888') . $salt)]);
            }
        }
        $this->adminBindRole($admin->id);
        $this->adminBindDataRange($admin->id, EnterpriseDataRange::TYPE_STORE, $params['store_id']);
        $this->adminBindAvailableSystem($admin);
        $this->adminStatusEnable($admin);
    }
}