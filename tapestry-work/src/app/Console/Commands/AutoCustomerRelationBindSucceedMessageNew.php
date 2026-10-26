<?php


namespace App\Console\Commands;


use App\Services\Message\CustomerRelationBindMessageNewService;
use Illuminate\Console\Command;

class AutoCustomerRelationBindSucceedMessageNew extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:customer-relation-bind-succeed-new';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时检测顾客关系绑定成功的客户给导购消息-新';

    /**
     * 检测数据时间间隔
     * @var int
     */
    protected $interval = 300;

    public function handle()
    {
        (new CustomerRelationBindMessageNewService())->autoExec($this->interval);
    }
}