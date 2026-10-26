<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;

class DepartmentListController extends BaseController
{

    /**
     * 门店列表
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = Sdk::work()
            ->setApi('service/v1/store/department-list')
            ->post($params);
        return success($result);
    }
}
