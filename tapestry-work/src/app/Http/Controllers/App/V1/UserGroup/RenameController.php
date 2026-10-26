<?php
/**
 * 修改分组昵称
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/24
 * Time: 22:06
 */

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use Illuminate\Http\Request;

class RenameController extends BaseController
{
    public function index(Request $request)
    {
        $rules = [
            'group_id' => 'required',
            'name'  => 'required|max:6'
        ];
        $message = [
            'group_id.required' => '分组id不能为空',
            'name.required' => '分组名称不能为空',
            'name.max'  => '分组名称不能超过:max个字符',
        ];

        $this->validate($request, $rules, $message);
        $groupId = $request->input('group_id');
        $name = $request->input('name');

        $group = UserGroup::find($groupId);

        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        $updateGroup = $group->update(['name' => $name]);
        return success([
            'data' => [
                'updateGroup'   => $updateGroup
            ]
        ]);
    }
}
