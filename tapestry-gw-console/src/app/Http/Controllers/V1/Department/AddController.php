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

class AddController extends BaseController
{
    /**
     * 新建部门
     *
     * @param Request $request
     * @return array|mixed
     */
    public function index(Request $request)
    {
        if (empty($this->qyDomain)) {
            return success(['data' => []]);
        }

        // 部门名称
        $dpName = trim($request->post('dp_name'));
        if (!$dpName) {
            error(320102);
        }

        // 上级部门id
        $dpParentId = $request->post('dp_parent_id', '');
        if (!$dpParentId) {
            error(320103);
        }

        // 是否设置为门店，默认非门店
        $isStore = $request->post('is_store', Constant::DEPARTMENT_IS_STORE_FALSE);

        $dptId = '';

        $Department = new Department();

        // 设置为门店
        if ($isStore) {
            // 调用sdk获取门店dptId
            $dptIdData = $Department->storeDptId([
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier
            ]);

            // 部门门店类型id
            $dptId = $dptIdData['data']['dptId'] ?? '';
        }

        // 调用SDK
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'dpName'     => $dpName,
                'dpParentid' => $dpParentId,
                'dptId'      => $dptId
            ]
        ];
        $detail    = $Department->add($condition);

        $data = [];
        // 数据格式化
        if ($detail) {
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
                'is_store'            => $isStore,
                'dp_display_order'    => $detail['dpDisplayorder']
            ];
        }

        // 写入记录日志
        $content = $isStore . ':' ? '门店' : '非门店';
        $this->writeOperatorLog('新增部门:' . $dpName . $content);

        return success($data);
    }
}
