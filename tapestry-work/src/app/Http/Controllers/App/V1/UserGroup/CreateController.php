<?php

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;
use Illuminate\Http\Request;

class CreateController extends BaseController
{
    /**
     * 新建用户分组
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
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

        $params = [
            'name'          => $request->input('name'),
            'enterprise_id' => $this->enterpriseId,
            'mem_uid'       => $this->memUid,
        ];

        $group = UserGroup::where($params)->first();

        if ($group) {
            throw new Error(100000, '分组已存在');
        }

        $ids = $request->ids;
        array_walk($ids, function (&$item) {
            $item['enterprise_id'] = $this->enterpriseId;
        });

        $group = UserGroup::create($params);

        $customers = $group->members()->createMany($ids);

        return success([
            'data' => [
                'group'     => $group,
                'customers' => $customers
            ]
        ]);
    }
}
