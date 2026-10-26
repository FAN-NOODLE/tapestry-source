<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class MoveEmployeeDepartmentController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only(['employee_id_arr', 'department_id', 'enterprise_id']);

        if (empty($params['employee_id_arr']) || empty($params['department_id']) || empty($params['enterprise_id'])) {
            throw new Error(600009, '员工id、门店id均不能为空');
        } else {
            StoreService::updateEmployeeDepartment($params);
        }

        return success([]);
    }
}
