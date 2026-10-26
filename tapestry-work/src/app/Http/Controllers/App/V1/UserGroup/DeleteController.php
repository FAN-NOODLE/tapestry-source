<?php

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;

class DeleteController extends BaseController
{
    /**
     * 删除用户分组
     *
     * @param         $groupId
     *
     * @return array
     */
    public function index($groupId)
    {
        $group = UserGroup::find($groupId);

        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        $deleteGroup = $group->forceDelete();

        $deleteMember = $group->members()->forceDelete();

        return success([
            'data' => [
                'deleteGroup'  => $deleteGroup,
                'deleteMember' => $deleteMember
            ]
        ]);
    }
}
