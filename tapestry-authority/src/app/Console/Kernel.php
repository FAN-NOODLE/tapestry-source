<?php

namespace App\Console;

use App\Console\Commands\AccessTokenRefresh;
use App\Console\Commands\Document;
use App\Console\Commands\Message\AutoSendActivityMessage;
use App\Console\Commands\Message\AutoSyncTemplates;
use App\Console\Commands\Message\CsvTemplatesMini;
use App\Console\Commands\Message\GetJsonForTemplateKeyValues;
use App\Console\Commands\Message\SyncMessageScene;
use App\Console\Commands\Message\SyncMessageTids;
use App\Console\Commands\Message\SyncTemplateKeyRules;
use App\Console\Commands\Message\SyncTemplateKeyValues;
use App\Console\Commands\Message\SyncThirdPartAccountCategory;
use App\Console\Commands\Message\InitMessageSceneAndChannel;
use App\Console\Commands\Rbac\InitBuiltInRole;
use App\Console\Commands\Rbac\InitRole;
use App\Console\Commands\Rbac\PermissionRelease;
use App\Console\Commands\Rbac\SyncAdmin;
use App\Console\Commands\Rbac\SyncAdminPermission;
use App\Console\Commands\Message\SyncMessageNumbers;
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
        Document::class,
        AutoSyncTemplates::class,
        CsvTemplatesMini::class,
        SyncMessageTids::class,
        SyncThirdPartAccountCategory::class,
        SyncTemplateKeyValues::class,
        SyncTemplateKeyRules::class,
        SyncMessageScene::class,
        GetJsonForTemplateKeyValues::class,
        InitRole::class,
        InitBuiltInRole::class,
        PermissionRelease::class,
        SyncAdminPermission::class,
        SyncAdmin::class,
        SyncMessageNumbers::class,
        AutoSendActivityMessage::class,
        InitMessageSceneAndChannel::class,
        AccessTokenRefresh::class
    ];

    /**
     * Define the application's command schedule.
     *
     * @param \Illuminate\Console\Scheduling\Schedule $schedule
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        // 每天凌晨执行同步订阅消息
        $schedule->command('auto:sync_templates')->dailyAt('00:00');

        // 每分钟验证发送活动开始提醒订阅消息
        $schedule->command('auto:send_activity_message')->everyMinute()->withoutOverlapping();
    }
}
