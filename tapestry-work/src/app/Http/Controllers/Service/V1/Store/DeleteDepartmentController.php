<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class DeleteDepartmentController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only(['department_id', 'qy_domain', 'enterprise_id']);

        if (empty($params['department_id'])) {
            throw new Error(600009, '门店id不能为空');
        } else {
            StoreService::deleteDepartment($params);
        }

        return success([]);
    }
}
