<?php

namespace App\Console\Commands\Rbac;

use App\Library\CodeTimer;
use App\Models\Permissions;
use App\Modules\Rbac\Role;
use App\Services\PermissionReleaseService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

/**
 * 权限发布
 *
 * @package App\Console\Commands\Rbac
 */
class PermissionRelease extends Command
{
    protected $signature = 'rbac:permission-release';

    protected $description = '权限发布';

    public function handle()
    {
        $this->info('完成以下5个步骤将发布权限');

        $operator = $this->ask("1.告诉我你是谁?");
        $systemType = $this->ask('2.发布系统类型' . PHP_EOL . $this->systemTypeDesc());
        $versionTag = $this->ask('3.发布迭代号 全部: all');
        $enterpriseID = $this->ask('4.发布企业ID 全部: 0');

        if ($this->confirm('5.离发布就差一步了, 是否确认发布?' . PHP_EOL . "  -系统类型: {$systemType}, 迭代号: {$versionTag}, 发布企业: {$enterpriseID}")) {
            $codeTimer = new CodeTimer();
            $codeTimer->open();

            Log::info("{$this->description}执行", [
                '操作人'  => $operator,
                '系统类型' => $systemType,
                '迭代号'  => $versionTag,
                '企业'   => $enterpriseID
            ]);

            $this->info("发布超管权限, 即将开始执行");

            $prs = new PermissionReleaseService($systemType, $versionTag, $enterpriseID);
            $enterpriseAdminRelationList = $prs->enterpriseAdminRelationList();

            // 进度条
            $bar = $this->output->createProgressBar(count($enterpriseAdminRelationList));
            $bar->start();

            foreach ($enterpriseAdminRelationList as $enterpriseAdmin) {
                try {
                    $role = new Role($enterpriseAdmin['admin_id'], $enterpriseAdmin['enterprise_id']);
                    $res = $role->permissions->superAdminRelease($systemType, $versionTag);
                    if ($res === '没有权限需要处理') {
                        // 如果进到这里说明找不到权限数据，无需再遍历执行了over
                        $this->info(PHP_EOL . "{$res}");
                        break;
                    }
                } catch (\Exception $exception) {
                    Log::error('发布超管权限异常', [
                        'message' => $exception->getMessage(),
                        'file'    => $exception->getFile(),
                        'line'    => $exception->getLine(),
                        'args'    => [
                            'system_type' => $systemType,
                            'version_tag' => $versionTag
                        ]
                    ]);
                }
                $bar->advance();
            }

            $bar->finish();
            $this->info(PHP_EOL);
            $codeTimer->close();

            $this->info(PHP_EOL . "发布超管权限结束, 耗时:{$codeTimer->getRunTime()}ms");
        } else {
            $this->info("已取消发布");
        }
    }

    private function systemTypeDesc()
    {
        $systemTypeDesc = '';
        foreach (Permissions::CONSOLE_SYSTEM_TYPE_RELATION as $type => $desc) {
            $systemTypeDesc .= "  -{$desc}: {$type}" . PHP_EOL;
        }
        $systemTypeDesc .= "  -全部: 0" . PHP_EOL;

        return $systemTypeDesc;
    }
}
