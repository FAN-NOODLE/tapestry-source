<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\Department;

use App\Common\Constant;
use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 部门详情
     *
     * @param Request $request
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 部门id
        $input = $request->only([
            'dp_id'
        ]);

        // 调用sdk
        $condition  = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'dp_id'          => $input['dp_id']
        ];
        $Department = new Department();
        $detail     = $Department->detail($condition);

        $data = [];
        // 数据格式化
        if ($detail) {
            // 调用sdk获取门店dptId
            $dptIdData = $Department->storeDptId([
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier
            ]);

            // 部门门店类型id
            $dptId = $dptIdData['data']['dptId'] ?? '';

            // 是否设置为门店：否
            $isStore = Constant::DEPARTMENT_IS_STORE_FALSE;
            // 部门门店类型id
            $detailDptId = $detail['dptId'] ?? '';
            if ($dptId && $detailDptId == $dptId) {
                // 是否设置为门店：是
                $isStore = Constant::DEPARTMENT_IS_STORE_TRUE;
            }

            $data = [
                'dp_id'               => $detail['dpId'],
                'dp_name'             => $detail['dpName'],
                'dp_parent_id'        => $detail['dpParentid'],
                'is_child_department' => $detail['isChildDepartment'],
                'dp_display_order'    => $detail['dpDisplayorder'],
                'is_store'            => $isStore
            ];
        }

        return success($data);
    }
}
