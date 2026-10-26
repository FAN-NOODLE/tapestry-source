<?php

namespace App\Console\Commands;

use App\Models\EnterpriseLogSystem;
use App\Services\EnterpriseLogService;
use Illuminate\Console\Command;

class CleanSystemLog extends Command
{
    protected $signature = 'auto:clean_system_log';

    protected $description = '清理日志表记录';

    public function handle()
    {
        $startTime = microtime(true);

        // 今天0点前30天的日志
        $timestamp = strtotime(date('Y-m-d')) - 86400 * 30;

        // 获取n天前的最后一条记录ID
        $id = EnterpriseLogSystem::query()->where('created_at', '<', $timestamp)->orderBy('id', 'DESC')->value('id');

        if ($id) {
            $num = EnterpriseLogSystem::query()->where('id', '<', $id)->forceDelete();
            EnterpriseLogService::info('清理日志表记录', [
                'runtime' => round((microtime(true) - $startTime) * 1000, 2).'ms',
                'num'     => $num
            ]);
        }

        return true;
    }
}