<?php


namespace App\Http\Controllers\Service\V1\StaffShare;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorker\PerformanceInheritService;

class BindingCustomerDetailController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
        ];

        $this->validate($this->request, $rules);

        $params = $this->request->only([
            'enterprise_id',
            'worker_ids',
            'page',
            'size',
            'is_page'
        ]);

        $result = (new PerformanceInheritService())->StaffList($params);
        return success($result);
    }

}