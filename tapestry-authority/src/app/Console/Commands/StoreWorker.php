<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class StoreWorker extends Command
{
    protected $signature = 'store_worker';

    protected $description = '';

    public function handle()
    {
        $workers = \Illuminate\Support\Facades\DB::table('store_worker')
            ->select(['id', 'enterprise_id', 'mobile', 'is_boss'])
            ->get();

        foreach ($workers as $worker) {

            $admin = \App\Models\EnterpriseAdmin::query()
                ->where('account', '=', '+86-' . $worker->mobile)
                ->where('enterprise_id', '=', $worker->enterprise_id)
                ->first();

            if ($admin) {

                // 检查管理员在不, 在就更新可用系统
                $oSystem = json_decode($admin->admin_available_system, true);
                if ($oSystem) {
                    $system = array_values(array_unique(array_merge($oSystem, [2])));
                } else {
                    if ($worker->is_boss == 1) {
                        $system = [1, 2];
                    } else {
                        $system = [2];
                    }
                }

                // 可用系统更新
                $admin->update(['admin_available_system' => json_encode($system)]);

                // 获取店长和店员的角色
                $roles = \App\Models\Role::query()
                    ->where('enterprise_id', '=', $admin->enterprise_id)
                    ->whereIn('mark', ['store_manager', 'store_member'])
                    ->pluck('id', 'mark')
                    ->toArray();
                // 是店长
                if ($worker->is_boss == 1) {
                    $role = \App\Models\AdminRoleRelation::query()
                        ->where('admin_id', '=', $admin->id)
                        ->where('role_id', '=', $roles['store_manager'])
                        ->first();
                    // 没绑定过 绑定下
                    if (empty($role)) {
                        \App\Models\AdminRoleRelation::query()->create([
                            'enterprise_id' => $admin->enterprise_id,
                            'admin_id'      => $admin->id,
                            'role_id'       => $roles['store_manager']
                        ]);
                    }
                } else {
                    // 是店员
                    $role = \App\Models\AdminRoleRelation::query()
                        ->where('admin_id', '=', $admin->id)
                        ->where('role_id', '=', $roles['store_member'])
                        ->first();
                    // 没绑定过 绑定下
                    if (empty($role)) {
                        \App\Models\AdminRoleRelation::query()->create([
                            'enterprise_id' => $admin->enterprise_id,
                            'admin_id'      => $admin->id,
                            'role_id'       => $roles['store_member']
                        ]);
                    }
                }

                // 绑定数据权限
                $dataRange = \App\Models\EnterpriseDataRange::query()
                    ->where('enterprise_id', '=', $admin->enterprise_id)
                    ->where('enterprise_admin_id', '=', $admin->id)
                    ->where('type', '=', 2)
                    ->first();
                if (empty($dataRange)) {
                    \App\Models\EnterpriseDataRange::query()->create([
                        'enterprise_id'       => $admin->enterprise_id,
                        'enterprise_admin_id' => $admin->id,
                        'type'                => 2,
                        'object_id'           => $worker->id
                    ]);
                }
            }
        }
    }
}
