<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class EditController extends BaseController
{

    public function index()
    {
        $rules = [
            'id'               => 'required|integer',
            'enterprise_id'    => 'required',
            'mobile'           => 'required',
            'name'             => 'required',
            'store_id'         => 'required|integer'
        ];

        $message = [
            'id.required'           => 'id参数必填',
            'id.integer'            => 'id参数必须int类型',
            'enterprise_id.required'=> 'enterprise_id参数必填',
            'store_id.required'     => 'store_id参数必填',
            'store_id.integer'      => 'store_id参数必须int类型',
            'mobile.required'       => 'mobile参数必填',
            'name.required'         => 'name参数必填'
        ];

        $input = $this->validate($this->request, $rules, $message);
        StoreWorkerService::editStoreWorker($input);
        return success();
    }
}
