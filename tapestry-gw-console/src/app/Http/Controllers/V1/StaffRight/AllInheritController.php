<?php


namespace App\Http\Controllers\V1\StaffRight;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;

class AllInheritController extends BaseController
{

    public function index(Request $request)
    {
        $params = $request->only([
            'old_worker_id',
            'new_worker_id',
            'search_type',
            'search_value',
            'search_time_type',
            'start_time',
            'end_time'
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new StoreWorker())->performanceAllInherit($params));
    }
}