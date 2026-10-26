<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Validator;

class UpdateController extends BaseController
{

    /**
     * 更新信息
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'where',
            'data'
        ]);
        $this->validateParams($params);
        StoreWorkerService::update($params['where'],$params['data']);
        return success([]);
    }

    private function validateParams($params)
    {
        $rules    = [
            'where' => 'required|array',
            'data'  => 'required|array',
        ];
        $messages = [
            'where.required' => 'where参数必填',
            'where.integer'  => 'where必须为array类型',
            'data.required'  => 'data参数必填',
            'data.integer'   => 'data必须为array类型',
        ];
        Validator::make($params, $rules, $messages)->validate();
    }
}
