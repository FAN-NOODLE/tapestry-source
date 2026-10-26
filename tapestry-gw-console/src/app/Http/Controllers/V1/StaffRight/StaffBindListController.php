<?php


namespace App\Http\Controllers\V1\StaffRight;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\GuideManage;
use CrsSdk\Work\StoreWorker;
use Illuminate\Http\Request;
use CrsSdk\Shop\StaffRight;
use Illuminate\Validation\Rule;

class StaffBindListController extends BaseController
{

    public function index(Request $request)
    {
        $params = $request->only([
            'worker_id',
            'search_type',
            'search_value',
            'search_time_type',
            'start_time',
            'end_time',
            'page',
            'size',
        ]);
        $params['enterprise_id'] = $this->enterpriseId;
        return success((new StoreWorker())->getPerformanceInheritCustomerList($params));
    }


}