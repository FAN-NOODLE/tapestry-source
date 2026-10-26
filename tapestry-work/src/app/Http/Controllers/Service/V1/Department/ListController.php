<?php

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use App\Models\Scrm\Departments;
use App\Models\StoreDepartment;
use Exception;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 部门列表
     *
     * @param Request $request
     *
     * @return array
     * @throws Exception
     */
    public function index(Request $request): array
    {
        $qyDomain = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params = $request->post('params', []);
        $ignorePermissions = $request->post('ignorePermissions', false);

        $result = (new Department())->list($qyDomain, $appIdentifier, $params, $ignorePermissions);

        return success($result);
    }

    /**
     * 部门列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function lists(Request $request): array
    {
        $enterpriseId = $request->post('enterprise_id');
        $params = $request->post('params');
        if(!empty($params['dp_name'])) {
            $thirdIds = StoreDepartment::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('store_name', 'like', '%' . $params['dp_name'] . '%')
                ->pluck('wechat_department_id')
                ->toArray();
            if (!empty($thirdIds)) {
                $params['store_name_third_id'] = $thirdIds;
            }
        }
        $departments = $this->getDepartmentList($params, $enterpriseId);

        $departments = $this->formatData($departments, $enterpriseId);

        return success($departments);
    }

    public function getDepartmentList($params, $enterpriseId): array
    {
        return Departments::query()
            ->with(['child' => function ($query) {
                $query->where('deleted_at', 0);
            }])
            ->select('id as department_id', 'third_id', 'third_id as dp_id', 'third_id as dp_third_id', 'third_parent_id as dp_third_parent_id', 'third_parent_id as dp_parent_id', 'name as dp_name')
            ->when(isset($params['dp_name']) && isset($params['store_name_third_id']), function ($query) use ($params) {
                return $query->where(function ($q) use ($params) {
                    $q->where('name', 'like', '%' . $params['dp_name'] . '%')
                        ->orWhereIn('third_id', $params['store_name_third_id']);
                });
            })
            ->when(isset($params['dp_name']) && !isset($params['store_name_third_id']), function ($query) use ($params) {
                return $query->where('name', 'like', '%' . $params['dp_name'] . '%');
            })
            ->when(isset($params['dp_id']), function ($query) use ($params) {
                return $query->where('third_parent_id', $params['dp_id']);
            })
            ->when(isset($params['third_ids']), function ($query) use ($params) {
                return $query->whereIn('third_id', $params['third_ids']);
            })
            ->when(isset($params['third_id']), function ($query) use ($params) {
                return $query->where('third_id', $params['third_id']);
            })
            ->when(isset($params['third_parent_id']), function ($query) use ($params) {
                return $query->where('third_parent_id', $params['third_parent_id']);
            })
            ->where('enterprise_id', $enterpriseId)
            ->where("deleted_at",0)
            ->orderBy('id', 'ASC')
            ->get()
            ->toArray();
    }

    public function formatData($departments, $enterpriseId): array
    {
        $storeDepartments = StoreDepartment::query()
            ->where('enterprise_id', $enterpriseId)
            ->get(['wechat_department_id', 'third_store_code', 'store_name', 'store_type'])
            ->keyBy('wechat_department_id')
            ->toArray();

        foreach ($departments as $k => $v) {
            $departments[$k]['store_type'] = isset($storeDepartments[$v['dp_third_id']]) ? $storeDepartments[$v['dp_third_id']]['store_type'] : 0;
            $departments[$k]['store_name'] = isset($storeDepartments[$v['dp_third_id']]) ? $storeDepartments[$v['dp_third_id']]['store_name'] : $v['dp_name'];
            $departments[$k]['is_child_department'] = count($v['child']) > 0 ? 1 : 0;
            $departments[$k]['dp_store_name'] = isset($storeDepartments[$v['dp_third_id']]) ? $storeDepartments[$v['dp_third_id']]['third_store_code'] . "-" . $storeDepartments[$v['dp_third_id']]["store_name"] : $v['dp_name'];
            unset($departments[$k]['child']);
        }

        return $departments;

    }
}
