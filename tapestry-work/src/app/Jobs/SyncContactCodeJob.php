<?php

namespace App\Jobs;

use App\Services\StaffManagerService;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class SyncContactCodeJob extends Job
{

    protected $input;

    /**
     * Create a new job instance.
     *
     * @return void
     */
    public function __construct($input)
    {
        $this->input = $input;
//        $this->connection = 'redis';
//        $this->queue = 'sync-contact';
    }

    /**
     * Execute the job.
     *
     * @return void
     */
    public function handle()
    {
        Log::info('后台内购设置自建联系人二维码--队列开始', ['$input' => $this->input, 'time' => Carbon::now()]);

        //为内购员工自建联系人二维码
        $result = (new StaffManagerService())->syncContactCodeByMem($this->input);

        Log::info(
            '后台内购设置自建联系人二维码--队列结束', ['$input'    => $this->input, 'time' => Carbon::now(),
                                     '$ucResult' => $result ?? '失败']
        );
    }
}
