<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;

class SyncDataByDepartmentsController extends BaseController
{

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(): \Illuminate\Http\JsonResponse
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = Sdk::work()
            ->setApi('service/v1/store/sync-data-by-departments')
            ->post($params);
        return success($result);
    }
}
