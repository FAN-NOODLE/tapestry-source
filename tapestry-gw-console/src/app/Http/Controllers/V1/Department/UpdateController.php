<?php

namespace App\Http\Controllers\V1\Department;

use App\Common\Constant;
use CrsSdk\Work\Department;
use Illuminate\Http\Request;
use App\Http\Controllers\V1\BaseController;

class UpdateController extends BaseController
{
    /**
     * 编辑组织
     *
     * @param Request $request
     * @return array|mixed
     */
    public function index(Request $request)
    {
        $input = $request->input();

        // 组织id
        if (!isset($input['dp_id'])) {
            notice(320101);
        }
        // 组织名称
        if (!isset($input['dp_name'])) {
            notice(320102);
        }
        // 上级组织id
        if (!isset($input['dp_parent_id'])) {
            notice(320103);
        }
        // 是否设置为门店 默认非门店
        $isStore = $input['is_store'] ?? Constant::DEPARTMENT_IS_STORE_FALSE;

        // 组织类型id
        $dptId      = '';
        $Department = new Department();

        // 设置为门店
        if ($isStore) {
            // 调用sdk获取门店dptId
            $dptIdData = $Department->storeDptId([
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier
            ]);

            // 组织门店类型id
            $dptId = $dptIdData['data']['dptId'] ?? '';
        }

        // 调用sdk
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => [
                'dpId'       => $input['dp_id'],
                'dpName'     => $input['dp_name'],
                'dpParentid' => $input['dp_parent_id'],
                'dptId'      => $dptId
            ]
        ];
        $detail    = $Department->edit($condition);

        $data = [];

        // 数据格式化
        if ($detail) {

            // 调用详情SDK返回最新数据
            $detail = $Department->detail([
                'qy_domain'      => $this->qyDomain,
                'app_identifier' => $this->appIdentifier,
                'dp_id'          => $input['dp_id']
            ]);

            $data = [
                'dp_id'               => $detail['dpId'],
                'dp_name'             => $detail['dpName'],
                'dp_parent_id'        => $detail['dpParentid'],
                'is_child_department' => $detail['isChildDepartment'],
                'is_store'            => $isStore
            ];
        }
        // 写入记录日志
        $this->writeOperatorLog('编辑组织:' . $input['dp_name']);

        return success($data);
    }
}
