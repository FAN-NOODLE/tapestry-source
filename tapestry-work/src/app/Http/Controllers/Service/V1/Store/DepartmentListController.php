<?php

namespace App\Http\Controllers\Service\V1\Store;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Scrm\Members;
use App\Services\StoreService;

class DepartmentListController extends BaseController
{

    /**
     * 门店列表
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = $this->request->only([
            'enterprise_id',
            'page',
            'size',
            'code_or_name',
            'store_type',
            'close_status',
            'is_page'
        ]);
        if (empty($params['enterprise_id'])) {
            throw new Error(600000, '企业id不能为空');
        }
        $fields = ['id', 'uc_department_id', 'wechat_department_id', 'third_store_code', 'store_name', 'store_type',
            'operate_name', 'sync_status', 'sync_at', 'address', 'email', 'popup_store_status', 'popup_store_keys','inventory_name',
            'sa_share_order_url', 'sa_check_order_url', 'sa_check_role_status', 'sa_check_role', 'phone', 'city', 'closed_at', 'created_at', 'updated_at'];
        $result = StoreService::departmentList($params, $params['is_page'] ?? true, $fields);

        $userIdArr = [];
        $useridStrArr = array_column($result['list'], 'sa_check_role');
        foreach ($useridStrArr as $val) {
            if (!empty($val)) {
                $userIdArr = array_merge($userIdArr, explode(',', $val));
            }
        }
        $userIdArr = array_unique($userIdArr);
        $memberList = Members::query()->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('userid', $userIdArr)
            ->get(['userid', 'name', 'gender'])
            ->keyBy('userid')
            ->toArray();
        foreach ($result['list'] as &$val) {
            $val['sa_check_role_arr'] = [];
            if (!empty($val['sa_check_role'])) {
                foreach (explode(',', $val['sa_check_role']) as $userid) {
                    $val['sa_check_role_arr'][] = $memberList[$userid] ?? [];
                }
            }
        }
        return success($result);
    }
}
