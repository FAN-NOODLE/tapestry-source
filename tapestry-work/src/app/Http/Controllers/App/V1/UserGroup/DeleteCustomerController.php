<?php
/**
 * 分组删除客户
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/22
 * Time: 18:25
 */

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use App\Models\UserGroupRelation;
use Illuminate\Http\Request;

class DeleteCustomerController extends BaseController
{
    public function index(Request $request)
    {
        $rules = [
            'group_id' => 'required',
            'user_ids' => 'required|array'
        ];
        $message = [
            'group_id.required' => '分组Id不能为空',
            'user_ids.required' => '客户Id集合不能为空',
            'user_ids.array' => '客户Id集合必须是数组',
        ];

        $this->validate($request, $rules, $message);

        $groupId = $request->input('group_id');
        $userIds = $request->input('user_ids');

        $group = UserGroup::query()->with(['members' => function ($query) use ($userIds) {
            $query->whereIn('ec_userid', $userIds);
        }])->find($groupId);

        if (!$group->toArray()) {
            throw new Error(100000, '该分组不存在');
        }
        $groupUser = UserGroupRelation::find($userIds);
        //检测分组是否存在客户
        if (!$groupUser) {
            throw new Error(100000, '该分组不存在此成员');
        }

        //不一定是全部删除 ，所以循环删除
        foreach ($group->members as $member) {
            $member->delete();
        }

        return success([
            'data' => []
        ]);
    }
}
