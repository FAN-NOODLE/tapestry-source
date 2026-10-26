<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class CustomerBindingWorkerDetailController extends BaseController
{
    /**
     * 查询绑定关系
     *
     * @return array|\Illuminate\Http\JsonResponse|\Psr\Http\Message\ResponseInterface
     * @throws \Exception
     */
    public function index()
    {
        // 接收数据
        $rules = [
            'enterprise_id' => 'required|integer',
            'user_id'       => 'required|integer',
        ];
        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID不合法',
            'user_id.required'      => '用户ID不能为空',
            'user_id.integer'       => '用户ID不合法',
        ];
        $input = $this->validate($this->request, $rules, $message);
        $workerInfo = StoreWorkerService::getBindingWorker($input['enterprise_id'], $input['user_id']);

        return success($workerInfo);
    }
}
