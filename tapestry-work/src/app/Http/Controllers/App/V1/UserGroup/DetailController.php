<?php

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\Customer\CommonController;
use App\Models\UserGroup;
use App\Models\UserGroupRelation;

class DetailController extends CommonController
{
    /**
     * 用户分组详情
     * @param         $groupId
     *
     * @return array
     */
    public function index($groupId)
    {
        $group = UserGroup::select(['id', 'name', 'mem_uid'])->find($groupId);
        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        $members = UserGroupRelation::where('group_id', $groupId)
            ->has('members')
            ->with([
                'members' => function ($query) {
                    $query->where('mem_uid', '=', $this->memUid)->select([
                        'mem_uid',
                        'ec_userid',
                        'ec_avatar',
                        'ec_name',
                        'union_id',
                        'user_remark_name',
                        'user_name_index'
                    ]);
                }
            ])
            ->get();

        $group->members_count = $members->count();

        $members = $members->pluck('members')->flatten();

        $group->members = $this->formatList($members);

        return success(['data' => $group]);
    }
}
