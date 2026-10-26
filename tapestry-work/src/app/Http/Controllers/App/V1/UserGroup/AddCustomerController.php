<?php
/**
 * 分组添加客户
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/22
 * Time: 18:23
 */

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use Illuminate\Http\Request;

class AddCustomerController extends BaseController
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

        $group = UserGroup::query()->with('members')->find($groupId);
        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        //检测是否存在已添加客户, 有则去重
        $groupUserIds = array_column($group->toArray()['members'], 'ec_userid');
        $userIds = array_diff($userIds, array_intersect($userIds, $groupUserIds));
        if (!$userIds) {
            throw new Error(100000, '该分组成员已存在');

        }

        //组装插入数据格式，开始插入
        $data = [];
        foreach ($userIds as $userId) {
            $data[] = [
                'ec_userid' => $userId,
                'enterprise_id' => $this->enterpriseId
            ];
        }
       $group->members()->createMany($data);

        return success([
            'data' => []
        ]);
    }
}
