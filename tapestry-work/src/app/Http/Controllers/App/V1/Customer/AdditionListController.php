<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\GuideUserRelation;
use App\Models\MemberUserRelation;
use App\Models\User;

class AdditionListController extends BaseController
{
    /**
     * 待添加客户列表
     * 
     *
     * @return array
     */
    public function index()
    {
        $size = request()->get('size', 10);
        $type = request()->get('type', 0);
        $isCopyMobile = request()->get('is_copy_mobile', GuideUserRelation::IS_COPY_MOBILE_NO);

        // 已添加的客户
        $list = MemberUserRelation::query()
            ->where('mem_uid', $this->memUid)
            ->where('delete_type', '<>', MemberUserRelation::DEL_EXTERNAL_CONTACT)
            ->pluck('union_id')
            ->toArray();

        // 未添加的客户
        $list = User::query()
            ->join('guide_user_relation', function ($join) {
                $join->on([
                    ['user.id', '=', 'guide_user_relation.user_id']
                ]);
            })
            ->where('guide_user_relation.mem_uid', $this->memUid)
            ->where('guide_user_relation.is_copy_mobile', $isCopyMobile)
            ->whereNotIn('user.union_id', $list)
            ->where('user.mobile', '<>', '')
            ->where('guide_user_relation.user_id', '>', 0)
            ->when($type > 0, function ($query) use ($type) {
                $query->where('user.type', $type);
            })
            ->where('user.deleted_at', 0)
            ->where('guide_user_relation.deleted_at', 0)
            ->paginate($size, [
                'user.id',
                'user.nick_name',
                'user.avatar',
                'user.type',
                'user.mobile',
                'guide_user_relation.is_copy_mobile'
            ])
            ->toArray();


        return success([
            'data' => $list['data'],
            'meta' => [
                'total' => $list['total'],
                'size' => (int)$list['per_page'],
                'page' => $list['current_page']
            ]
        ]);
    }
}
