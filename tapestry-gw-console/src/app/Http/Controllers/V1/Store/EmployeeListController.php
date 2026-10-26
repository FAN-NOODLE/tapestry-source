<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;

class EmployeeListController extends BaseController
{

    /**
     * 获取导购信息
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = Sdk::work()
            ->setApi('service/v1/store/employee-list')
            ->post($params);
        return success($result);
    }
}
