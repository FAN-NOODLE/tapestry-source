<?php

namespace App\Console;

use App\Console\Commands\AddQyWechatOrder;
use App\Console\Commands\AutoAddGuide;
use App\Console\Commands\AutoAllocateUserToGuide;
use App\Console\Commands\AutoCIIEExpire;
use App\Console\Commands\AutoCIIEExport;
use App\Console\Commands\AutoCIIESendMessage;
use App\Console\Commands\AutoCloseGuideAllocation;
use App\Console\Commands\AutoContactMember;
use App\Console\Commands\AutoCreateMarketingPlanPreviewPoster;
use App\Console\Commands\AutoCustomerRelationBindSucceedMessageNew;
use App\Console\Commands\AutoDeleteCoachItBySftp;
use App\Console\Commands\AutoExternalMissing;
use App\Console\Commands\AutoGetStoreBySftp;
use App\Console\Commands\AutoGetSyncContactStatus;
use App\Console\Commands\AutoNoticeMarketingPlanTask;

//use App\Console\Commands\AutoPullCDPUserActData;
use App\Console\Commands\AutoQueryCashCouponInfo;
use App\Console\Commands\AutoRefreshWelcomeSpeechMediaId;
use App\Console\Commands\AutoSourceMediaId;
use App\Console\Commands\AutoSourceMediaIdNew;
use App\Console\Commands\AutoStaffData;
use App\Console\Commands\AutoSyncContact;
use App\Console\Commands\AutoSyncContactCode;
use App\Console\Commands\AutoUpdateMarketingPlanPoster;
use App\Console\Commands\AutoUpdateMarketingPlanTaskResource;
use App\Console\Commands\AutoUpdateUcByStore;
use App\Console\Commands\BatchAddMemberToGuide;
use App\Console\Commands\CountOrderFromShop;
use App\Console\Commands\CustomerGroupRelation;
use App\Console\Commands\CustomerTransferResult;
use App\Console\Commands\DailySyncUserBehaviorData;
use App\Console\Commands\Document;
use App\Console\Commands\ExternalPayOrderResult;

//use App\Console\Commands\FixDBCommand;
use App\Console\Commands\ForbidLeaveAdmin;
use App\Console\Commands\HandleGuidePoster;
use App\Console\Commands\InitializeMemberBehaviorData;
use App\Console\Commands\InitTomorrowCustomerStatistic;
use App\Console\Commands\MigrateHistoryExternal;
use App\Console\Commands\SaveUserBehaviorData;
use App\Console\Commands\StoreBillFilesByEmail;
use App\Console\Commands\SyncDailyBehaviorStatisticData;
use App\Console\Commands\SyncDataLakeInfo;
use App\Console\Commands\SynchronousUserAndOrder;
use App\Console\Commands\SyncMemberUserRelation;
use App\Console\Commands\SyncStoreDepartmentByCoach;

//use App\Console\Commands\Temp\MigrateHistoryBindingCustomerStatistic;
use App\Console\Commands\SyncUCMember;
use App\Console\Commands\Temp\MigrateHistoryCustomerStatistic;
use App\Console\Commands\Temp\MigrateHistoryMemberCustomerStatistic;
use App\Console\Commands\Temp\TempGenerateWorkerPoster;
use App\Console\Commands\Temp\TempMemberRelation;
use App\Console\Commands\TestCommand;
use App\Console\Commands\TestCrmSdk;
use App\Console\Commands\UnbindStaffShare;
use App\Console\Commands\UpdateGuideQr;
use App\Console\Commands\UpdateMemberDataStatistic;
use App\Console\Commands\UpdateMemberExtendRoleField;
use App\Console\Commands\UpdateMemberGuide;
use App\Console\Commands\UpdateOrderRelationDpId;
use App\Console\Commands\UpdateQyWechatOrder;
use App\Console\Commands\UpdateShoppingGuidePoster;
use App\Console\Commands\UpdateStaffRight;
use App\Console\Commands\UpdateStoreManagerByPosition;
use App\Console\Commands\UpdateStoreWorkerInfo;
use App\Console\Commands\UpdateTodayCumulativeCustomerStatistic;
use App\Console\Commands\UpdateUcGuideQr;
use App\Console\Commands\WithdrawalAutoAudit;
use App\Console\Commands\YilouDataOrderTransfer;
use App\Console\Commands\YilouDataSowTransfer;
use App\Console\Commands\YilouDataTransfer;
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
        'doc'                                        => Document::class,
        'auto:allocate_user_to_guide'                => AutoAllocateUserToGuide::class,
        'auto:close_guide_allocation'                => AutoCloseGuideAllocation::class,
        //        'auto:pull_cdp_user_act_data'                => AutoPullCDPUserActData::class,  // 客服小程序同步旧用户和订单数据
        'synchronous-user-and-order'                 => SynchronousUserAndOrder::class,  // 更新员工扩展表字段role
        'update-member-extend-role-field'            => UpdateMemberExtendRoleField::class,
        'auto:refresh_welcome_speech_media_id'       => AutoRefreshWelcomeSpeechMediaId::class,
        //        'fix_db_command'                             => FixDBCommand::class,
        'auto:query_cash_coupon_info'                => AutoQueryCashCouponInfo::class,
        'customer_group_relation'                    => CustomerGroupRelation::class,
        'update-order-relation-dp-id'                => UpdateOrderRelationDpId::class,
        'update-staff-right'                         => UpdateStaffRight::class,
        'auto:sync_contact'                          => AutoSyncContact::class,
        'auto:sync_contact_code'                     => AutoSyncContactCode::class,
        'auto:sync_staff_data'                       => AutoStaffData::class,
        'auto:get_sync_contact_status'               => AutoGetSyncContactStatus::class,
        'update-member-guide {type}'                 => UpdateMemberGuide::class,
        'update-guide-qr {type}'                     => UpdateGuideQr::class,
        'update-uc-guide-qr {type}'                  => UpdateUcGuideQr::class,
        'auto:source_media_id'                       => AutoSourceMediaId::class,
        'auto:source_media_id_new'                   => AutoSourceMediaIdNew::class,
        'count_order_from_shop {--enterprise=}'      => CountOrderFromShop::class,
        'auto:sync_guide_data'                       => AutoContactMember::class,
        'temp:member-customer-relation'              => TempMemberRelation::class,
        'auto:withdrawal-auto-audit'                 => WithdrawalAutoAudit::class,
        'auto:customer-relation-bind-succeed-new'    => AutoCustomerRelationBindSucceedMessageNew::class,
        'handle_guide_poster'                        => HandleGuidePoster::class,
        'update_shopping_guide_poster'               => UpdateShoppingGuidePoster::class,
        'auto:notice_marketing_plan_task'            => AutoNoticeMarketingPlanTask::class,
        'auto:update_marketing_plan_task_resource'   => AutoUpdateMarketingPlanTaskResource::class,
        'auto:update_marketing_plan_poster'          => AutoUpdateMarketingPlanPoster::class,
        'auto:create_marketing_plan_preview_poster'  => AutoCreateMarketingPlanPreviewPoster::class,
        'UnbindStaffShare'                           => UnbindStaffShare::class,
        'init_tomorrow_customer_statistic'           => InitTomorrowCustomerStatistic::class,
        'update_today_cumulative_customer_statistic' => UpdateTodayCumulativeCustomerStatistic::class,
        'sync_daily_behavior_statistic_data'         => SyncDailyBehaviorStatisticData::class,
        'daily-sync-user-behavior-data'              => DailySyncUserBehaviorData::class,
        //        'migrate-history-binding-customer-statistic' => MigrateHistoryBindingCustomerStatistic::class,
        'migrate-history-member-customer-statistic'  => MigrateHistoryMemberCustomerStatistic::class,
        'migrate-history-customer-statistic'         => MigrateHistoryCustomerStatistic::class,
        'initialize-member-behavior-data'            => InitializeMemberBehaviorData::class,
        'sync_datalake_info'                         => SyncDataLakeInfo::class,
        'update-member-data-statistic'               => UpdateMemberDataStatistic::class,
        'migrate_history_external'                   => MigrateHistoryExternal::class,
        'auto_external_missing'                      => AutoExternalMissing::class,
        'auto:update-store-worker-info'              => UpdateStoreWorkerInfo::class,
        'auto:sync-member-user-relation'             => SyncMemberUserRelation::class,
        'test'                                       => TestCommand::class,
        'batch_add_member_to_guide'                  => BatchAddMemberToGuide::class,
        'temp:generate-worker-poster'                => TempGenerateWorkerPoster::class,
        'forbid-leave-admin'                         => ForbidLeaveAdmin::class,
        'auto_add_guide'                             => AutoAddGuide::class,
        'auto:get_store_by_sftp'                     => AutoGetStoreBySftp::class,
        'auto:update_uc_by_store'                    => AutoUpdateUcByStore::class,
        'sync-store-department-by-coach'             => SyncStoreDepartmentByCoach::class,
        'update-store-manager-by-position'           => UpdateStoreManagerByPosition::class,
        'save-user-behavior-data'                    => SaveUserBehaviorData::class,
        'add-qy-wechat-order'                        => AddQyWechatOrder::class,
        'test-crm-sdk'                               => TestCrmSdk::class,
        'auto:customer-transfer_result'              => CustomerTransferResult::class,
        'auto:external_pay_order_result'             => ExternalPayOrderResult::class,
        'auto:ciie:send_message'                     => AutoCIIESendMessage::class,
        'auto:ciie:expire'                           => AutoCIIEExpire::class,
        'auto:ciie:export'                           => AutoCIIEExport::class,
        'sync_uc_member'                             => SyncUCMember::class,
        UpdateQyWechatOrder::class,
        AutoDeleteCoachItBySftp::class,
        YilouDataTransfer::class,
        YilouDataSowTransfer::class,
        YilouDataOrderTransfer::class,
        StoreBillFilesByEmail::class
    ];

    /**
     * Define the application's command schedule.
     *
     * @param Schedule $schedule
     *
     * @return void
     */
    protected function schedule(Schedule $schedule)
    {
        // 每一分钟给客户分配导购
//        $schedule->command('auto:allocate_user_to_guide')
//            ->everyMinute()->withoutOverlapping();

        // 每一分钟关闭已经分配够足额的导购的分配权限
//        $schedule->command('auto:close_guide_allocation')
//            ->everyMinute()->withoutOverlapping();

        // 每天晚11:00调用cdp接口同步当日的客户访问量
//        $schedule->command('auto:pull_cdp_user_act_data')
//            ->dailyAt('23:00');


        // 定时任务更新欢迎语mediaId
//        $schedule->command('auto:refresh_welcome_speech_media_id')
//            ->everyFifteenMinutes()->withoutOverlapping();
//        $schedule->command('auto:refresh_welcome_speech_media_id')
//            ->dailyAt('03:00');

        // 每一分钟定时查询红包发送记录
//        $schedule->command('auto:query_cash_coupon_info')
//            ->everyMinute()->withoutOverlapping();

        // 每分钟执行一次从队列获取用户和分组关系
//        $schedule->command('customer_group_relation')
//            ->everyMinute();

        // 定时任务同步通讯录
        $schedule->command('auto:sync_contact')
            ->dailyAt('03:00');

        // 定时获取不同企业的同步状态
        $schedule->command('auto:get_sync_contact_status')
            ->everyMinute();

//        // 定时任务同步不同企业内购员工联系人二维码
//        $schedule->command('auto:sync_staff_data')
//            ->everyMinute()->withoutOverlapping();
//
        // 定时任务同步不同企业联系人二维码
//        $schedule->command('auto:sync_contact_code')
//            ->everyFifteenMinutes()->withoutOverlapping();

        // 定时任务更新素材图片mediaIdRedis
//        $schedule->command('auto:source_media_id')
//            ->dailyAt('03:00');

        // 每天早9点向导购发送营销计划任务
        $schedule->command('auto:notice_marketing_plan_task')
            ->dailyAt('09:00')->runInBackground()->withoutOverlapping();
        // 每分钟执行一次营销任务素材是否需要更新检测
        $schedule->command('auto:update_marketing_plan_task_resource')
            ->everyMinute()->runInBackground()->withoutOverlapping();
        //营销任务生成海报
        $schedule->command('auto:update_marketing_plan_poster')
            ->everyMinute()->runInBackground()->withoutOverlapping();
        //营销任务生成预览海报
        $schedule->command('auto:create_marketing_plan_preview_poster')
            ->everyMinute()->runInBackground()->withoutOverlapping();
        // 定时任务更新素材图片mediaIdRedis
//        $schedule->command('auto:source_media_id_new')
//            ->dailyAt('03:00');

        // 每分钟执行一次同步导购数据
//        $schedule->command('auto:sync_guide_data')
//            ->everyFiveMinutes();
        // 每分钟执行一次同步自动审批
//        $schedule->command('auto:withdrawal-auto-audit')
//            ->everyMinute();

        $schedule->command('auto:customer-relation-bind-succeed-new')
            ->everyMinute();
        //每分钟更新在职导购没有mem_userid的用户
        $schedule->command('auto:update-store-worker-info')
            ->everyMinute();


        // 每天02:00同步UC外部联系人缺失数据
        $schedule->command('auto_external_missing')
            ->dailyAt('02:00');

        //DL会员数据更新
        $schedule->command('update-member-data-statistic')
            ->dailyAt('03:00');

        // 每天04:00更新顾客行为数据
        $schedule->command('daily-sync-user-behavior-data')
            ->dailyAt('04:00');
        // 每天00:01更新顾客累计统计数据
        $schedule->command('update_today_cumulative_customer_statistic')
            ->dailyAt('00:01');
        // 每天04:00初始化明天顾客统计数据
        $schedule->command('init_tomorrow_customer_statistic')
            ->dailyAt('04:01');
        // 每半个小时执行CIIE吊牌预约活动模板消息发送
        $schedule->command('auto:ciie:send_message')
            ->everyThirtyMinutes();
        //  每天0点01分执行CIIE吊牌预约活动过期
        $schedule->command('auto:ciie:expire')
            ->dailyAt('00:01');

    }
}
