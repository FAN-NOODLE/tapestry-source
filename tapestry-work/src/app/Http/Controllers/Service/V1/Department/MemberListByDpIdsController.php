<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:27
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Scrm\Departments;
use Illuminate\Http\Request;

class MemberListByDpIdsController extends BaseController
{


    /**
     * 根据部门id获取员工列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request): array
    {
        $enterpriseId = $request->post('enterprise_id', '');
        $dpIds = $request->post('dp_ids', '');
        if (!is_array($dpIds)) {
            $dpIds = [$dpIds];
        }
        $dpList = Departments::query()->with("membersRel.members")
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('third_id', $dpIds)
            ->get(['id'])
            ->toArray();
        $memberList = [];
        foreach ($dpList as $val) {
            foreach ($val['members_rel'] as $rel) {
                $memberList[] = $rel['members'];
            }
        }
        return success($memberList);
    }
}
