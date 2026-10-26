<?php

namespace App\Jobs;

use App\Console\Commands\MigrateHistoryExternal;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class SyncHistoryExternalJob extends Job
{

    /**
     * 队列名称
     *
     * @var string
     */
    public $queue = 'sync_history_external';

    protected $data;

    /**
     * Create a new job instance.
     *
     * @param $data
     */
    public function __construct($data)
    {
        $this->data  = $data;

    }

    /**
     * Execute the job.
     *
     * @return void
     * @throws \Exception
     */
    public function handle()
    {
        Log::info('添加导购自动同步历史外部联系人', ['$input' => $this->data, 'time' => Carbon::now()]);
        //添加导购自动同步历史外部联系人
        (new MigrateHistoryExternal())->handle($this->data['enterprise_id'], $this->data['worker_ids']);

    }
}
