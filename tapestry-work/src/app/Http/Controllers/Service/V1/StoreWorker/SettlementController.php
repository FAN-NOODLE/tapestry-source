<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class SettlementController extends BaseController
{

    public function getWorkersByName()
    {
        $enterpriseId = $this->request->input('enterprise_id');
        $workerName   = $this->request->input('work_name');
        $data         = StoreWorkerService::getWorkerIdByName($enterpriseId, $workerName);
        return success($data);
    }

    public function getBatchWorkersInfo()
    {
        $params = $this->request->input();
        $data         = StoreWorkerService::getWorkersInfo($params);
        return success($data);
    }

}