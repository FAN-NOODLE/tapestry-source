<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class UpdateWorkerInfoController extends BaseController
{

    /**
     * 定时修改导购信息
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $rules = [
            'member_list' => 'required|array',
            'domain' => 'required',
        ];
        $message = [
            'domain.required' => 'domain参数必填',
            'member_list.required' => '员工参数必填',
        ];
        $input = $this->validate($this->request, $rules, $message);

        StoreWorkerService::updateWorkerInfo($input);
        return success();
    }
}
