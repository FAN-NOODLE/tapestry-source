<?php

namespace App\Console\Commands\Rbac;

use App\Models\Role;
use App\Models\RolePermissions;
use App\Modules\BuiltInRole\RoleStrategyConfig;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

/**
 * 同步管理员权限和超管保持一致
 *
 * @package App\Console\Commands\Rbac
 */
class SyncAdmin extends Command
{
    protected $signature = 'rbac:sync-admin';

    protected $description = '同步管理员权限和超管保持一致';

    public function handle()
    {
        $superAdmins = Role::query()->where('mark', '=', Role::ROLE['super_admin']['mark'])->get()->toArray();

        // 进度条
        $bar = $this->output->createProgressBar(count($superAdmins));
        $bar->start();

        foreach ($superAdmins as $superAdmin) {
            $admin = Role::query()
                ->where('enterprise_id', '=', $superAdmin['enterprise_id'])
                ->where('mark', '=', RoleStrategyConfig::ROLE['console_admin_role']['mark'])
                ->first();
            if ($admin !== null) {
                $adminArr = $admin->toArray();

                DB::transaction(function () use ($superAdmin, $adminArr) {
                    $thisTime = time();
                    $superAdminPermissions = RolePermissions::query()
                        ->where('role_id', '=', $superAdmin['id'])
                        ->pluck('permissions_id')
                        ->toArray();

                    $insertData = [];
                    foreach ($superAdminPermissions as $permissionID) {
                        $insertData[] = [
                            'role_id'        => $adminArr['id'],
                            'permissions_id' => $permissionID,
                            'enterprise_id'  => $superAdmin['enterprise_id'],
                            'created_at'     => $thisTime,
                            'updated_at'     => $thisTime
                        ];
                    }

                    RolePermissions::query()->where('role_id', '=', $adminArr['id'])->forceDelete();
                    RolePermissions::query()->insert($insertData);
                });
            } else {
                Log::warning('管理员不存在', [
                    'enterprise_id' => $superAdmin['enterprise_id']
                ]);
            }

            $bar->advance();
        }
        $bar->finish();
        $this->info(PHP_EOL);
    }
}
