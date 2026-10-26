<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\ScrmMemberService;
use App\Services\StoreWorkerService;

class DetailController extends BaseController
{

    /**
     * 获取导购信息
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();
        if (isset($params['use_scrm_db'])) {
            $result = ScrmMemberService::detail($params);
        } else {
            $result = StoreWorkerService::detail($params);

        }
        return success($result);
    }
}
