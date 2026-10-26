<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class LeaveController extends BaseController
{
    /**
     * 员工离职（删除员工）
     * @param Request $request
     * @return JsonResponse|void
     */
    public function index(Request $request)
    {
        // 员工id
        $memUid = $request->get('mem_uid', '');
        if (!$memUid) {
            error(320151);
        }

        // 调用sdk删除员工
        $condition = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => ['memUid' => $memUid]
        ];
        $detail    = (new Member())->detail($condition);
        (new Member())->delete($condition);

        // 写入记录日志
        $memUsername = $detail['memUsername'] ? $detail['memUsername'] : '';
        $this->writeOperatorLog('离职员工：' . $memUsername);

        return success();
    }
}
