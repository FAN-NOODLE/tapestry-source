<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class EmployeeListController extends BaseController
{

    /**
     * 获取导购信息
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'enterprise_id',
            'page',
            'size',
            'mobile_or_name',
            'third_store_code',
            'status',
            'is_page',
            'no_department'
        ]);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $fields = ['id', 'name', 'alias_name', 'employee_code',
                   'job_location', 'job_position', 'gender',
                   'status', 'third_store_code', 'mobile', 'created_at', 'sync_status', 'sync_at'];
        $result = StoreService::employeeList($params, $params['is_page'] ?? true, $fields);
        return success($result);
    }
}
