<?php


namespace App\Http\Controllers\Service\V1\StoreWorker;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class MemUidListController extends BaseController
{

    public function index()
    {
        $rules   = [
            'enterprise_id' => 'required',
            'ids'           => 'required|array'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'ids.required'           => '导购ID不能为空',
            'ids.array'              => '导购ID数据类型错误'
        ];
        $this->validate($this->request, $rules, $message);
        $params = $this->request->only(['ids', 'enterprise_id']);
        $result = StoreWorkerService::getMemUid($params['enterprise_id'], $params['ids']);
        return success($result);
    }
}
