<?php


namespace App\Console\Commands;


use App\Services\StoreWorker\AddGuideInMember\AddGuideService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Console\Command;

class AutoAddGuide extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto_add_guide';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '自动添加员工为导购';

    public function handle()
    {
        $data = (new Enterprise())->listsOfficial([]);
        foreach ($data['data'] as $enterprise){
            AddGuideService::addGuideNotInStoreWorker($enterprise['id'],$enterprise['domain']);
        }
        echo "同步完成";
    }
}