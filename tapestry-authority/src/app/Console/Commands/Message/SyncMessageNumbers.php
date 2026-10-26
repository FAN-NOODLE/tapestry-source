<?php


namespace App\Console\Commands\Message;


use App\Models\EnterpriseMessageScene;
use App\Models\SubscribeMessageRecord;
use App\Services\ActivityMessageService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class SyncMessageNumbers extends Command
{
    protected $signature = 'sync_message_number';

    protected $description = '同步订阅消息订阅人数和发送人数';

    public function handle()
    {
        $result = (new ActivityMessageService())->SyncMessageNumbers();
        $this->info(json_encode($result));
    }
}