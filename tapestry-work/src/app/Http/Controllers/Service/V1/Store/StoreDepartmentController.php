<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class StoreDepartmentController extends BaseController
{

    /**
     * 导购部门信息
     * @return array
     * @throws \Exception
     */
    public function detail(): array
    {
        $params = $this->request->only([
            'enterprise_id',
            'mobile',
        ]);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $result = StoreService::departmentDetail($params);
        return success($result);
    }
}
