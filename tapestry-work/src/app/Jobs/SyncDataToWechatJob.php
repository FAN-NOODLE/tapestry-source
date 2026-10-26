<?php
/**
 * 同步业务架构数据至企微架构
 */

namespace App\Jobs;

use App\Console\Commands\AutoUpdateUcByStore;

class SyncDataToWechatJob extends Job
{
    public $queue = "sync_data_to_wechat_job";
    private $enterpriseId;

    public function __construct($enterpriseId)
    {
        $this->enterpriseId = $enterpriseId;
    }

    public function handle()
    {
        (new AutoUpdateUcByStore())->handle($this->enterpriseId);
    }
}