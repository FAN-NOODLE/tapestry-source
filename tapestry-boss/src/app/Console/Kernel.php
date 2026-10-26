<?php

namespace App\Console;

use App\Console\Commands\AddPhonePrefix;
use App\Console\Commands\AutoDelOpenTemplate;
use App\Console\Commands\CleanSystemLog;
use App\Console\Commands\Document;
use App\Console\Commands\TaskCheckAppTokenIsInvalid;
use App\Console\Commands\UpdateAuditStatus;
use App\Console\Commands\SyncMiniCategory;
use App\Console\Commands\EnterpriseAccountSync;
use App\Console\Commands\UpdateRegisterStatus;
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
        'doc'                      => Document::class,
        'add_phone_prefix'         => AddPhonePrefix::class,
        'auto:update_audit_status' => UpdateAuditStatus::class,
        'auto:del_open_template'   => AutoDelOpenTemplate::class,
        'auto:clean_system_log'    => CleanSystemLog::class,
        'auto:sync_mini_category'  => SyncMiniCategory::class,
        'task:check-app-token-is-invalid' => TaskCheckAppTokenIsInvalid::class,
        'auto:enterprise_account_sync' => EnterpriseAccountSync::class,
        'auto:update_register_status' => UpdateRegisterStatus::class
    ];

    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule  $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        // 定时删除开放平台模版 每1小时执行一次 ,早上8点到晚上11点
        $schedule->command('auto:del_open_template')->cron('0 */1 * * *')->between('8:00', '23:00')->withoutOverlapping(30);

        // 清除日志
        $schedule->command('auto:clean_system_log')->dailyAt('2:00')->withoutOverlapping();

        // 定时更新refresh_token
        $schedule->command('task:check-app-token-is-invalid')->hourly()->withoutOverlapping();

        // 每分钟同步企业付费子账号
        $schedule->command('auto:enterprise_account_sync')->everyMinute()->withoutOverlapping();

        // 每十分钟查询创建小程序任务状态
        //$schedule->command('auto:update_register_status')->cron('*/10 * * * *')->withoutOverlapping();
    }
}
