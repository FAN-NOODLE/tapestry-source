<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\User;

class AdditionTotalController extends BaseController
{
    /**
     * 待添加客户数量
     *
     *
     * @return array
     */
    public function index()
    {
        // 已添加的客户
        $list = MemberUserRelation::query()
            ->where('mem_uid', $this->memUid)
            ->pluck('union_id')
            ->toArray();

        // 未添加的客户
        $total = User::query()
            ->join('guide_user_relation', function ($join) {
                $join->on([
                    ['user.id', '=', 'guide_user_relation.user_id']
                ]);
            })
            ->where('guide_user_relation.mem_uid', $this->memUid)
            ->whereNotIn('user.union_id', $list)
            ->count();

        return success([
            'data' => ['total' => $total]
        ]);
    }
}
