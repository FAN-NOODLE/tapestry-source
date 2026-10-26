<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 14:11
 */

namespace App\Http\Controllers\V1\Department;

use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class DeleteController extends BaseController
{
    /**
     * 删除部门
     *
     * @param Request $request
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 部门id
        $dpId = $request->post('dp_id', '');
        if (!$dpId) {
            error(320101);
        }

        // 调用sdk
        $condition = [
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params' => ['dpId' => $dpId]
        ];
        // 获取部门名称
        $params = [
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'dp_id' => $dpId
        ];
        $Department = new Department();

        $detail = $Department->detail($params);

        // 存在子级不能删除
        if (isset($detail['isChildDepartment']) && $detail['isChildDepartment'] === 1) {
            notice(320105);
        } else if (isset($detail['departmentMemberCount']) && $detail['departmentMemberCount'] > 0) {
            // 存在员工不能删除
            notice(320106);
        }

        $Department->delete($condition);

        // 写入记录日志
        $dpName = $detail['dpName'] ? $detail['dpName'] : '';
        $this->writeOperatorLog('删除部门:'.$dpName);

        return success();
    }
}
