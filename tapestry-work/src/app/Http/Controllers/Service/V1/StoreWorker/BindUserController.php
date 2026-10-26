<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class BindUserController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'mem_uid'       => 'required',
            'user_id'       => 'required|integer',
        ];
        $message = [
            'enterprise_id.required' => 'enterprise_id参数必填',
            'enterprise_id.integer'  => 'enterprise_id必须为int类型',
            'mem_uid.required'       => 'mem_uid参数必填',
            'user_id.required'       => 'user_id参数必填',
            'user_id.integer'        => 'user_id必须为int类型',
        ];
        $input = $this->validate($this->request, $rules, $message);
        StoreWorkerService::workerBindUser($input);
        return success();
    }
}
