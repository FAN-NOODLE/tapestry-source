<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StoreWorker;
use App\Services\StoreWorkerService;

class AddController extends BaseController
{

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input   = $this->request->input();
        $rules   = [
            'enterprise_id' => 'required',
            'store_id'      => 'required',
            'mem_uid'       => 'required',
            'operator_id'   => 'integer',
            'is_boss'       => 'required|in:' . StoreWorker::BOSS_YES . ',' . StoreWorker::BOSS_NO,
            'source'        => 'required|integer',
            'user_id'       => 'integer',
            'mobile'        => 'string'
        ];
        $message = [
            'enterprise_id.required' => 'enterprise_id参数必填',
            'is_boss.required'       => 'is_boss参数必填',
            'is_boss.in'             => 'is_boss参数错误',
            'operator_id.integer'    => 'operator_id参数必须int类型',
            'mem_uid.required'       => 'mem_uid参数必填',
            'source.required'        => '员工来源必填',
            'store_id.required'      => '门店ID必填',
            'user_id.integer'        => 'user_id参数必须int类型',
        ];
        $this->validate($this->request, $rules, $message);
        $result = StoreWorkerService::addStoreWorker($input, $input['mobile'] ?? '');
        return success($result);
    }
}
