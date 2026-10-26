<?php

namespace App\Console;

use App\Console\Commands\AutoCleanOperationLog;
use App\Console\Commands\Document;
use App\Console\Commands\ExportCustomer;
use App\Console\Commands\TaskStandard;
use App\Console\Commands\UpdateSourceName;
use Illuminate\Console\Scheduling\Schedule;
use Laravel\Lumen\Console\Kernel as ConsoleKernel;

class Kernel extends ConsoleKernel
{
    /**
     * The Artisan commands provided by your application.
     *
     * @var array
     */
    protected $commands = [
        'doc' => Document::class,
        'task:task_standard_number' => TaskStandard::class,
        'auto:clean_operation_log'  => AutoCleanOperationLog::class,
        'update_source_name'        => UpdateSourceName::class,
        'export:customer-list' => ExportCustomer::class,
    ];

    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        // 每小时定时更新统计表
        $schedule->command('task:task_standard_number')->hourly();
        // 每日01点00分执行删除日志操作
        $schedule->command('auto:clean_operation_log')->dailyAt('01:00');
    }
}
