<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class DetailByIdOrMemUidController extends BaseController
{

    /**
     * 获取导购信息
     *
     * @return array
     * @author zhangjunchang<zhangjunchang@vchangyi.com>
     * @date 2020/10/30 14:54
     */
    public function index()
    {
        $workerId = $this->request->input('worker_id');
        if (empty($workerId)) {
            throw new Error(2300003, '门店员工ID必传');
        }
        $workerInfo = StoreWorkerService::detailByIdOrMemUid($workerId);
        if (empty($workerInfo)) {
            throw new Error(2300006, '门店员工不存在');
        }
        return success($workerInfo);
    }
}
