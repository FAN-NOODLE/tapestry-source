<?php

namespace App\Console\Commands\Rbac;

use App\Models\Enterprise;
use Illuminate\Console\Command;
use App\Modules\BuiltInRole\RoleStrategy;

class InitBuiltInRole extends Command
{
    protected $signature = 'rbac:init-built-role';

    protected $description = '初始化内置角色 {1初始化内置角色|2处理门店员工数据|3都处理} {企业ID|all}';

    /**
     * @throws \Exception
     */
    public function handle()
    {
        $type = $this->argument('type');
        $enterpriseId = $this->argument('enterprise_id');

        if (empty($type) || !in_array($type, [1, 2, 3])) {
            echo '类型错误';
            die;
        }
        if (empty($enterpriseId)) {
            echo '企业ID必填';
            die;
        }

        if ($enterpriseId == 'all') {
            $ids = Enterprise::query()->where('status', '=', 1)->where('is_split_flow', '=', 0)->pluck('id')->toArray();
        } else {
            $e = Enterprise::query()->where('id', '=', $enterpriseId)->first();
            if (empty($e)) {
                echo '企业ID错误';
                die;
            }
            $ids = [$enterpriseId];
        }

        if ($type == 1) {
            RoleStrategy::initFullRoles($ids);
        }

        if ($type == 2) {
            RoleStrategy::storeRoleRelation($ids);
        }
        if ($type == 3) {
            RoleStrategy::initFullRoles($ids);
            RoleStrategy::storeRoleRelation($ids);
        }
    }
}
