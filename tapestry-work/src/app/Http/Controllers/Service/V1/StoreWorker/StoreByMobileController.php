<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\StoreWorkerService;

class StoreByMobileController extends BaseController
{


    /**
     * 在职导购或手机号查询门店ID
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules   = [
            'enterprise_id' => 'required'
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
        ];
        $this->validate($this->request, $rules, $message);
        $conditions                  = $this->request->input();
        $conditions['worker_status'] = $conditions['worker_status'] ?? StoreWorker::IN_SERVICE;
        $result                      = StoreWorkerService::getStoreIdListByMobile($conditions);
        return success($result);
    }

}
