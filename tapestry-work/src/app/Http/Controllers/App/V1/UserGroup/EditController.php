<?php

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use Illuminate\Http\Request;

class EditController extends BaseController
{
    /**
     * 编辑用户分组
     *
     * @param Request $request
     *
     * @param         $groupId
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request, $groupId)
    {
        $this->validate($request, [
            'name'            => 'required|max:6',
            'ids'             => 'required|array',
            'ids.*.ec_userid' => 'required'
        ], [
            'name.required'            => '分组名称不能为空',
            'name.max'                 => '分组名称不能超过:max个字符',
            'ids.required'             => 'uc外部联系人数组不能为空',
            'ids.*.ec_userid.required' => 'uc外部联系人不能为空'
        ]);

        $group = UserGroup::find($groupId);
        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        $updateGroup = null;
        if ($request->input('name')) {
            $updateGroup = $group->update(['name' => $request->input('name')]);
        }

        // 用户分组编辑组内客户是，采用硬删除，重新添加，防止生成过多冗余数据
        $updateMembers = null;
        if ($request->input('ids')) {
            $group->members()->forceDelete();

            $ids = $request->ids;
            array_walk($ids, function (&$item) {
                $item['enterprise_id'] = $this->enterpriseId;
            });
            $updateMembers = $group->members()->createMany($ids);
        }

        return success([
            'data' => [
                'updateGroup'   => $updateGroup,
                'updateMembers' => $updateMembers
            ]
        ]);
    }
}
