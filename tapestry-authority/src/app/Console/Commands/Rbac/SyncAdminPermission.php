<?php

namespace App\Console\Commands\Rbac;

use App\Models\Role;
use Illuminate\Console\Command;
use App\Models\RolePermissions;

class SyncAdminPermission extends Command
{
    protected $signature = 'rbac:sync-admin-permission {enterprise_id}';

    protected $description = '重置管理员内置角色权限 {企业ID|all}';

    /**
     * @throws \Exception
     */
    public function handle()
    {
        $enterpriseId = $this->argument('enterprise_id');

        if (empty($enterpriseId)) {
            echo '企业ID必填';
            die;
        }

        if ($enterpriseId == 'all') {
            $roles = Role::query()
                ->where('mark', '=', 'console_admin_role')
                ->orderBy('enterprise_id')
                ->pluck('id', 'enterprise_id')
                ->toArray();
        } else {
            $roles = Role::query()
                ->where('mark', '=', 'console_admin_role')
                ->where('enterprise_id', '=', $enterpriseId)
                ->orderBy('enterprise_id')
                ->pluck('id', 'enterprise_id')
                ->toArray();
        }

        if (empty($roles)) {
            echo '未查询到内置管理员角色';
            die;
        }

        $time = time();
        foreach ($roles as $eid => $roleId) {

            echo '正在处理企业:id' . $eid . PHP_EOL;

            $super = Role::query()->where('mark', '=', 'super_admin')->where('enterprise_id', '=', $eid)->first();
            if ($super) {
                $ps = RolePermissions::query()->where('role_id', '=', $super->id)->pluck('permissions_id')->toArray();
                if ($ps) {
                    $data = [];
                    foreach ($ps as $p) {
                        $data[] = [
                            'enterprise_id'  => $eid,
                            'permissions_id' => $p,
                            'role_id'        => $roleId,
                            'created_at'     => $time,
                            'updated_at'     => $time
                        ];
                    }
                    if ($data) {
                        RolePermissions::query()
                            ->where('enterprise_id', '=', $eid)
                            ->where('role_id', '=', $roleId)
                            ->forceDelete();
                        RolePermissions::query()->insert($data);
                    }
                }
            }
        }
    }
}
