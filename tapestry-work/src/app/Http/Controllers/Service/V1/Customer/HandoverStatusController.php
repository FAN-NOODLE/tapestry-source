<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/6
 * Time: 16:18
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;

class HandoverStatusController extends BaseController
{
    /**
     * 批量获取离职员工客户交接状态
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 离职员工id集合
        $memUids = $request->post('mem_uids', []);
        if (!$memUids) {
            error(600007);
        }

        // 外部联系人标记列表
        $relationList = MemberUserRelation::whereIn('mem_uid', $memUids)->get()->toArray();

        // 以离职员工id分组的外部联系人标记列表
        $groupRelation = [];
        if ($relationList) {
            array_walk($relationList, function ($relation) use (&$groupRelation) {
                $groupRelation[$relation['mem_uid']][] = $relation;
            });
        }

        // 初始化返回数据
        $data = [];
        foreach ($memUids as $memUid) {
            // 离职员工客户是否分配完成状态：未完成
            $status = Constant::CUSTOMER_HANDOVER_COMPLETED_FALSE;

            // 离职员工外部联系人标记关系，有外部联系人标记关系->未分配完成
            $memRelation = $groupRelation[$memUid] ?? [];
            if (!$memRelation) {
                $status = Constant::CUSTOMER_HANDOVER_COMPLETED_TRUE;
            }

            $data[] = [
                'mem_uid'         => $memUid,
                'handover_status' => $status
            ];
        }

        return success(['data' => $data]);
    }
}
