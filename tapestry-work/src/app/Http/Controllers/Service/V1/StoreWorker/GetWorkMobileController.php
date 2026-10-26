<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Exceptions\Error;
use App\Models\StoreWorker;
use App\Http\Controllers\Service\V1\BaseController;

class GetWorkMobileController extends BaseController
{
    public function index()
    {
        $memUid       = $this->request->input('mem_uid');
        $enterpriseId = $this->request->input('enterprise_id');
        if (empty($memUid)) {
            throw new Error(407009, 'mem_uid不能为空');
        }
        if (empty($enterpriseId)) {
            throw new Error(110006, '企业ID不能为空');
        }
        $mobiles = StoreWorker::query()
            ->where('mem_uid', '=', $memUid)
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('worker_status', '=', StoreWorker::IN_SERVICE)
            ->pluck('mobile')
            ->toArray();
        return success($mobiles);
    }
}