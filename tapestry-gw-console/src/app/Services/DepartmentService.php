<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/3/20 10:24
 */

namespace App\Services;

use App\Common\Constant;

class DepartmentService extends BaseService
{

    /**
     * 部门列表数据格式化
     *
     * @param array  $list 待格式化部门列表
     * @param string $dptId 门店dptId
     *
     * @return array
     */
    public function formatData($list, $dptId, $departmentLists = [])
    {
        $data = [];
        if (empty($list)) {
            return $data;
        }

        foreach ($list as $department) {
            // 是否设置为门店：否
            $isStore = Constant::DEPARTMENT_IS_STORE_FALSE;
            // 部门门店类型id
            $departmentDptId = $department['dptId'] ?? '';

            if ($dptId && $departmentDptId == $dptId) {
                // 是否设置为门店：是
                $isStore = Constant::DEPARTMENT_IS_STORE_TRUE;
            }
            $data[] = [
                'dp_id'               => $department['dpId'],
                'dp_name'             => $department['dpName'],
                'dp_parent_id'        => $department['dpParentid'],
                'is_child_department' => $department['isChildDepartment'],
                'is_store'            => $isStore,
                'dp_user_count'       => $department['departmentMemberCount'],
                'dp_display_order'    => $department['dpDisplayorder'],
                'dp_third_id'         => $department['dpThirdid'],
                'dp_third_parent_id'  => isset($departmentLists[$department['dpThirdid']]['third_parent_id']) ? $departmentLists[$department['dpThirdid']]['third_parent_id'] : '',
                'department_id'       => isset($departmentLists[$department['dpThirdid']]['id']) ? $departmentLists[$department['dpThirdid']]['id'] : '',
                'store_type'          => isset($departmentLists[$department['dpThirdid']]['store_type']) ? $departmentLists[$department['dpThirdid']]['store_type'] : 0,
            ];
        }

        return $data;
    }

}
