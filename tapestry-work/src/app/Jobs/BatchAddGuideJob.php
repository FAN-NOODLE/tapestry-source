<?php

namespace App\Jobs;


use App\Services\StoreWorker\AddGuideInMember\AddGuideService;

class BatchAddGuideJob extends Job
{
    public $queue = 'batch-add-guide';

    private $enterpriseId;

    private $qyDomain;

    public function __construct($enterpriseId,$qyDomain)
    {
       $this->enterpriseId = $enterpriseId;
       $this->qyDomain = $qyDomain;
    }

    public function handle()
    {
        AddGuideService::addGuideNotInStoreWorker($this->enterpriseId,$this->qyDomain);
    }
}