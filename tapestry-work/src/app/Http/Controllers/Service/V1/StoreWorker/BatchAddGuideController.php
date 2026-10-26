<?php


namespace App\Http\Controllers\Service\V1\StoreWorker;


use App\Http\Controllers\Service\V1\BaseController;
use App\Jobs\BatchAddGuideJob;
use App\Library\Util\Logger;

class BatchAddGuideController extends BaseController
{
    public function index()
    {
        Logger::html("输入",$this->request->input());
        $enterpriseId = $this->request->input(['enterprise_id']);
        $qyDomain = $this->request->input(['qy_domain']);
        $job = new BatchAddGuideJob($enterpriseId,$qyDomain);
        $this->dispatch($job);
        return success([]);
    }
}