<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StoreService;

class SaveDepartmentInfoController extends BaseController
{

    /**
     * @return array
     * User: zhaoyang@vchangyi.com
     * Date: 2022/2/7 15:45
     */
    public function index()
    {
        $params = $this->request->only([
            'id',
            'store_name',
            'store_type',
            'store_code',
            'operate_name',
            'enterprise_id',
            'address',
            'email',
            'popup_store_status',
            'popup_store_keys',
            'sa_share_order_url',
            'sa_check_order_url',
            'sa_check_role_status',
            'sa_check_role',
            'city',
            'phone',
            'inventory_name',

        ]);
        if (empty($params['id'])) {
            if (empty($params['store_name']) || empty($params['store_type']) || empty($params['store_code'])) {
                throw new Error(600016, '门店名称、类型、编码均不能为空');
            }
            StoreService::addDepartment($params);
        } else {
            StoreService::saveDepartment($params);
        }

        return success([]);
    }
}
