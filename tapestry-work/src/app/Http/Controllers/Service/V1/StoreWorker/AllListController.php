<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Validator;

class AllListController extends BaseController
{

    /**
     * 获取导购列表
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'worker_no',
            'name',
            'mobile',
            'worker_status',
            'enterprise_id',
            'store_id',
            'is_boss',
            'type',
            'mem_uids',
            'store_ids',
            'ids',
            'columns'
        ]);
        //$this->validateParams($params);
        $result = StoreWorkerService::allList($params);
        return success($result);
    }

    private function validateParams($params)
    {
        $rules    = [
            'enterprise_id' => 'required|integer',
        ];
        $messages = [
            'enterprise_id.required' => 'enterprise_id参数必填',
            'enterprise_id.integer'  => 'enterprise_id必须为int类型',
        ];
        // 参数验证
        Validator::make($params, $rules, $messages)->validate();
    }
}
