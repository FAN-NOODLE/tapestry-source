<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;

class GetWorkerByIdsController extends BaseController
{

    /**
     * 获取导购列表
     * @return array
     */
    public function index()
    {
        $ids = $this->request->input('ids');
        if (!empty($ids) && array_filter($ids)) {
            $result = StoreWorkerService::getWorkerListByIds($ids);
            return success($result);
        }

        return success();
    }
}
