<?php

namespace App\Http\Controllers\App\V1\UserGroup;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\UserGroup;

class ListController extends BaseController
{
    /**
     * 用户分组列表
     * @return array
     */
    public function index()
    {
        $data = UserGroup::select('id', 'mem_uid', 'name')
            ->where('mem_uid', $this->memUid)->withCount([
                'members' => function ($query) {
                    $query->whereHas('members', function ($query) {
                        $query->where('mem_uid', '=', $this->memUid);
                    });
                }
            ])->get();

        return success(['data' => $data]);
    }
}
