<?php

namespace App\Http\Controllers\Service\V1\StoreWorker;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreWorkerService;
use Illuminate\Support\Facades\Validator;

class ListController extends BaseController
{

    /**
     * 获取导购信息
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'page',
            'size',
            'worker_no',
            'name',
            'mobile',
            'worker_status',
            'start_date',
            'end_date',
            'enterprise_id',
            'enterprise_ids',
            'store_id',
            'is_boss',
            'type',
            'is_page',
            'mem_uids',
            'user_ids',
            'store_ids',
            'ids',
            'contact_qr_code_id',
            'contact_qr_code_ids'
        ]);
        if ((!isset($params['enterprise_id']) || empty($params['enterprise_id'])) && (!isset($params['enterprise_ids']) || empty($params['enterprise_ids']))) {
            throw new Error(600000, '企业id不能为空');
        }
        $result = StoreWorkerService::list($params, $params['is_page'] ?? true);
        return success($result);
    }
}
