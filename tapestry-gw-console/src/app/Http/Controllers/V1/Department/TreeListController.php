<?php

namespace App\Http\Controllers\V1\Department;

use App\Common\Constant;
use App\Micros\Work\Work;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\Work\Department;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use App\Services\DepartmentService;
use App\Http\Controllers\V1\BaseController;

class TreeListController extends BaseController
{

    const STORE_TYPE_KEY = 'store_type';
    const ROOT_DEPT_ID_KEY = 'root_department_id';

    /**
     * 部门列表树形结构
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws \Exception
     */
    public function index(Request $request): JsonResponse
    {
        $input = $request->input();
        if (empty($this->qyDomain)) {
            return success(['list' => []]);
        }


        //获取企业导购区域配置
        $enterpriseConfig = (new EnterpriseConfig())->info([
            'enterprise_id' => $this->enterpriseId,
            'key'           => [self::STORE_TYPE_KEY, self::ROOT_DEPT_ID_KEY,]
        ]);

        $enterpriseConfig    = array_column($enterpriseConfig['data'], NULL, 'key');
        $storeType           = $enterpriseConfig[self::STORE_TYPE_KEY]['value'] ?? 0;
        $rootDeptId          = $enterpriseConfig[self::ROOT_DEPT_ID_KEY]['value'] ?? 0;

        $departmentLists = Work::lists(['enterprise_id' => $this->enterpriseId, 'params' => $input]);

        if (isset($input['dp_name']) && !empty($input['dp_name'])) {
            return success([
                'list' => $departmentLists,
            ]);
        }

        $dataTree = $this->getTree($departmentLists, $storeType, $rootDeptId, $rootDeptId);
        return success([
            'list' => $dataTree,
        ]);
    }

    /**
     * 递归获取树形结构
     *
     * @param        $data
     * @param        $storeType
     * @param        $pId
     * @param int    $dpThirdId
     *
     * @return array
     */
    public function getTree(&$data, $storeType, $pId, int $dpThirdId = 0): array
    {
        $tree = [];
        foreach ($data as $k => &$v) {

            if ($v['dp_third_id'] == $dpThirdId) {
                $v['children_list'] = $this->getTree($data, $storeType, $v['dp_third_id']);
                $tree[]             = $v;
                unset($data[$k]);
            } else {

                if ($v['dp_third_parent_id'] == $pId) {   // 父级找到子级
                    if ($storeType && $v['store_type'] == 0) {
                        if (substr($v['dp_name'], 0, 2) === "OC") {
                            $v['store_type']     = 1;
                            $v['store_type_str'] = '零售';
                        } elseif (substr($v['dp_name'], 0, 3) === "OCF") {
                            $v['store_type']     = 2;
                            $v['store_type_str'] = '奥莱';
                        } else {
                            $v['store_type']     = 0;
                            $v['store_type_str'] = '其他';
                        }
                    }
                    $v['children_list'] = $this->getTree($data, $storeType, $v['dp_third_id']);
                    $tree[]             = $v;
                    unset($data[$k]);

                }
            }
        }
        return $tree;
    }


}
