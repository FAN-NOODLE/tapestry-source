<?php


namespace App\Console\Commands\Message;


use App\Services\ActivityMessageService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

class AutoSendActivityMessage extends Command
{
    /**
     * @var string
     */
    protected $signature = 'auto:send_activity_message';

    /**
     * @var string
     */
    protected $description = '发送活动开始提醒订阅消息';

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        parent::__construct();
    }

    public function handle()
    {
        (new ActivityMessageService)->SendMessage();
    }
}