<?php


namespace App\Http\Controllers\Service\V1\StoreWorker\PerformanceInherit;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorker\PerformanceInheritService;
use Illuminate\Validation\Rule;

class CustomerListController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'worker_id' => 'required|integer',
            'search_type' => [
                'integer',
                Rule::in(PerformanceInheritService::SEARCH_TYPE_LIST)
            ],
            'search_value' => 'string',
            'search_time_type' => [
                'integer',
                Rule::in(PerformanceInheritService::SEARCH_TIME_TYPE_LIST)
            ],
            'start_time' => 'integer',
            'end_time' => 'integer',
            'page' => 'required|integer',
            'size' => 'required|integer',
        ];

        $this->validate($this->request, $rules);

        $params = $this->request->only([
            'enterprise_id',
            'mem_uid',
            'worker_id',
            'page',
            'size',
            'search_type',
            'search_value',
            'search_time_type',
            'start_time',
            'end_time'
        ]);

        switchAppHost(['enterprise_id' => $params['enterprise_id']]);
        $result = (new PerformanceInheritService())->customerList($params);
        $result['page'] = $params['page'];
        $result['size'] = $params['size'];
        return success($result);
    }
}