<?php


namespace App\Console\Commands;

use App\Models\DailyMemberCustomerStatistic;
use App\Models\StoreWorker;
use App\Services\DataStatistic\CustomerStatisticTaskService;
use App\Services\DataStatistic\DailyMemberCustomerStatisticService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;


class DailySyncUserBehaviorData extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'daily-sync-user-behavior-data';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '前一天导购行为数据同步';


    public function handle()
    {
        Log::info("DailySyncUserBehaviorData.autoExec start", ['time' => time()]);
        //初始化昨日新增导购数据
        $initializeTime = strtotime('yesterday');
        StoreWorker::query()
            ->select(['enterprise_id', 'mem_uid'])
            ->where(['worker_status' => StoreWorker::IN_SERVICE])
            ->where('created_at', '>', $initializeTime)
            ->groupBy(['mem_uid'])
            ->orderBy('enterprise_id', 'ASC')
            ->chunkById(700, function ($workerInfo) use ($initializeTime) {
                foreach ($workerInfo as $v) {
                    (new DailyMemberCustomerStatisticService)->getOrCreateMemberStatistic($v->enterprise_id, $v->mem_uid, $initializeTime);
                }
            });
        $syncEnterprise = DailyMemberCustomerStatistic::query()
            ->select('enterprise_id')
            ->where('stat_time', $initializeTime)
            ->groupBy(['enterprise_id'])
            ->get();
        foreach ($syncEnterprise as $v) {
            (new CustomerStatisticTaskService)->syncUserBehaviorByEnterpriseId($v->enterprise_id, [$initializeTime]);
        }
        Log::info("DailySyncUserBehaviorData.autoExec end", ['time' => time()]);
        return true;
    }
}