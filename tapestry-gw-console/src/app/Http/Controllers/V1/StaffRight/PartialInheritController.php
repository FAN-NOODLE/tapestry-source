<?php


namespace App\Http\Controllers\V1\StaffRight;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;

class PartialInheritController extends BaseController
{

    public function index(Request $request)
    {
        $params = $request->only([
            'old_worker_id',
            'new_worker_id',
            'user_ids'
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new StoreWorker())->performancePartialInherit($params));
    }
}