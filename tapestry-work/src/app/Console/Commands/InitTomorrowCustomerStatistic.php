<?php


namespace App\Console\Commands;


use App\Common\DateUtils;
use App\Services\DataStatistic\CustomerStatisticTaskService;
use Illuminate\Console\Command;

class InitTomorrowCustomerStatistic extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'init_tomorrow_customer_statistic {--statTime=} {--force} {--enterpriseId=}';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '初始化明天所有企业顾客统计数据';


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
            // 明日
            $statTime = DateUtils::tomorrowStatTime();
        }
        $force = $this->option('force');
        $enterpriseId = intval($this->option('enterpriseId'));
        if ($enterpriseId > 0) {
            (new CustomerStatisticTaskService())->initEnterpriseCustomerStatistic($enterpriseId, $statTime, $force);
        } else {
            (new CustomerStatisticTaskService())->initAllEnterpriseCustomerStatistic($statTime, $force);
        }
    }
}