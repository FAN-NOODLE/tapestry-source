<?php


namespace App\Console\Commands;


use App\Services\Message\CustomerRelationBindMessageService;
use Illuminate\Console\Command;

class AutoCustomerRelationBindSucceedMessage extends Command
{

    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:customer-relation-bind-succeed';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时检测顾客关系绑定成功的客户给导购消息';

    /**
     * 检测数据时间间隔
     * @var int
     */
    protected $interval = 300;

    public function handle()
    {
        (new CustomerRelationBindMessageService())->autoExec($this->interval);
    }
}