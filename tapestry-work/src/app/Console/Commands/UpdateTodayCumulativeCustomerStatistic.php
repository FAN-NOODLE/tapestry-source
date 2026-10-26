<?php


namespace App\Console\Commands;


use App\Common\DateUtils;
use App\Services\DataStatistic\CustomerStatisticTaskService;
use Illuminate\Console\Command;

class UpdateTodayCumulativeCustomerStatistic extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'update_today_cumulative_customer_statistic {--statTime=} {--force}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '跟新所有企业今日顾客累计统计数据';


    /**
     * CustomerStatisticService constructor.
     */
    public function __construct()
    {
        parent::__construct();
    }


    public function handle()
    {
        $statTime = intval($this->option('statTime'));
        if (empty($statTime)) {
            // 今日
            $statTime = DateUtils::todayStatTime();
        }
        $force = $this->option('force');
        (new CustomerStatisticTaskService())->updateAllEnterpriseCumulativeCustomerStatistic($statTime, $force);
    }
}