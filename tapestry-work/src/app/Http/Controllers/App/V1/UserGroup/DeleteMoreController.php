<?php
/**
 * 用户分组 批量删除
 * Created by PhpStorm.
 * User: zoulongbo
 * Date: 2019/1/18
 * Time: 10:07
 */

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use Illuminate\Http\Request;

class DeleteMoreController extends BaseController
{
    public function index(Request $request)
    {
        $groupIds = $request->input('groupIds');
        $group = UserGroup::find($groupIds);

        $deleteGroup = [];
        $deleteMember = [];

        if (!$group) {
            throw new Error(100000, '该分组不存在');
        }

        foreach ($group as $item) {
            $deleteGroup[] = $item->forceDelete();

            $deleteMember[] = $item->members()->forceDelete();
        }

        return success([
            'data' => [
                'deleteGroup' => array_unique($deleteGroup)
            ]
        ]);
    }
}
