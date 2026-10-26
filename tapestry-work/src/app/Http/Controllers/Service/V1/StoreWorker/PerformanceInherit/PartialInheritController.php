<?php


namespace App\Http\Controllers\Service\V1\StoreWorker\PerformanceInherit;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorker\PerformanceInheritService;

class PartialInheritController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'new_worker_id' => 'required|integer',
            'old_worker_id' => 'required|integer',
            'user_ids' => 'required|array',
        ];

        $this->validate($this->request, $rules);
        $params = $this->request->only([
            'enterprise_id',
            'new_worker_id',
            'old_worker_id',
            'user_ids'
        ]);

        switchAppHost(['enterprise_id' => $params['enterprise_id']]);

        (new PerformanceInheritService())->partialInherit($params);
        return success([]);
    }
}