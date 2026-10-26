<?php


namespace App\Jobs;


use App\Library\Uc\Member;
use App\Models\StoreWorker;
use App\Services\StoreWorker\PerformanceInheritService;

class SyncStoreWorkerJob extends Job
{
    public $queue = "sync_store_worker_job";
    private $qyDomain = "";
    private $enterpriseId = "";

    public function __construct($qyDomain, $enterpriseId)
    {
        $this->qyDomain     = $qyDomain;
        $this->enterpriseId = $enterpriseId;
    }

    public function handle()
    {
        $appIdentifier = "";
        $page          = 1;
        $skip          = 0;
        do {
            $list = StoreWorker::query()
                ->where('enterprise_id', $this->enterpriseId)
                ->skip($skip)
                ->take(500)
                ->get(['mem_uid'])
                ->toArray();

            $memUid       = array_column($list, 'mem_uid');
            $memberResult = (new Member())->ignorePermissionsList($this->qyDomain, $appIdentifier, ['memUids' => $memUid]);
            foreach ($memberResult['list'] as $mem) {
                StoreWorker::query()
                    ->where('enterprise_id', $this->enterpriseId)
                    ->where('mem_uid', $mem['memUid'])
                    ->update([
                        'name' => $mem['memUsername'],
                        'avatar'=>$mem['memFace']
                    ]);
            }
            $page++;
            $skip = ($page - 1) * 500;
        } while ($list != []);
        return true;
    }
}