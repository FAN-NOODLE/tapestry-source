<?php

namespace App\Console\Commands\Rbac;

use App\Modules\Rbac\Role;
use App\Models\EnterpriseAdmin;
use App\Models\EnterpriseConfig;
use App\Models\Role as RoleModel;
use App\Modules\Rbac\Permissions;
use Illuminate\Console\Command;
use App\Models\Permissions as PermissionsModel;

class InitRole extends Command
{
    protected $signature = 'rbac:init-role {enterprise_id}';

    protected $description = '初始化内置角色 {企业ID|all}';

    public function handle()
    {
        $enterpriseId = $this->argument('enterprise_id');

        // gw-console版
        $list = EnterpriseAdmin::query()
            ->select([
                'enterprise.id as enterprise_id',
                'enterprise_admin.id as admin_id',
                'enterprise.version'
            ])
            ->leftJoin('enterprise', 'enterprise_admin.enterprise_id', '=', 'enterprise.id')
            ->where('enterprise.version', '>', 1)
            ->where('enterprise_admin.is_super', '=', 1)
            ->when($enterpriseId !== 'all', function ($query) use ($enterpriseId) {
                return $query->where('enterprise.id', '=', $enterpriseId);
            })
            ->where('enterprise.deleted_at', '=', 0)
            ->get()
            ->toArray();

        $count = count($list);

        echo "待处理{$count}个企业" . PHP_EOL;

        // 过滤特殊权限
        $warehouseTypeRelation = EnterpriseConfig::query()
            ->whereIn('enterprise_id', array_column($list, 'enterprise_id'))
            ->where('key', '=', 'warehouse_type')
            ->pluck('value', 'enterprise_id')
            ->toArray();

        foreach ($list as $key => $value) {
            $num = $key + 1;

            echo "  {$num}.开始处理" . PHP_EOL;

            try {
                $role = new Role($value['admin_id'], $value['enterprise_id']);
                $permissions = new Permissions($value['version']);

                $filterWarehouse = PermissionsModel::filterWarehouse($warehouseTypeRelation[$value['enterprise_id']] ?? 0);

                // 给标配+选配权限
                $additional = array_column($permissions->getStandardConfig(), 'id');

                $role->init([
                    RoleModel::ROLE['super_admin']['mark'],
                    RoleModel::ROLE['store_manager']['mark'],
                    RoleModel::ROLE['store_member']['mark']
                ], $additional, $value['admin_id'], $filterWarehouse);
            } catch (\Exception $exception) {
                echo "    初始化异常: {$exception->getMessage()}" . PHP_EOL;
                echo "      file: {$exception->getFile()}" . PHP_EOL;
                echo "      line: {$exception->getLine()}" . PHP_EOL;
                continue;
            }

            echo "    初始化成功: {$value['enterprise_id']}" . PHP_EOL;
        }

        echo '全部初始化完成' . PHP_EOL;
    }
}
