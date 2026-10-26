<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\User;

class AddedListController extends BaseController
{
    /**
     * 已添加客户列表
     *
     *
     * @return array
     */
    public function index()
    {
        $size = request()->get('size', 10);
        $name = request()->get('name', '');

        $memberUserRelationList = MemberUserRelation::query()
            ->where('mem_uid', $this->memUid)
            ->where('delete_type', '<>', MemberUserRelation::DEL_EXTERNAL_CONTACT)
            ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->get()
            ->toArray();

        $userIds                = array_column($memberUserRelationList, 'user_id');
        $memberUserRelationList = array_column($memberUserRelationList, null, 'user_id');

        $list = User::query()
            ->where('enterprise_id', $this->enterpriseId)
            ->when(!empty($name), function ($query) use ($name) {
                $query->where(function ($query) use ($name) {
                    $query->where('nick_name', 'like', '%' . $name . '%')
                        ->orWhere('real_name', 'like', '%' . $name . '%');
                });
            })
            ->whereIn('id', $userIds)
            ->orderBy('id', 'desc')
            ->paginate($size, [
                'id',
                'nick_name',
                'avatar',
                'type',
            ])
            ->toArray();

        foreach ($list['data'] as &$value) {
            $value['delete_type'] = $memberUserRelationList[$value['id']]['delete_type'];
        }

        return success([
            'data' => $list['data'],
            'meta' => [
                'total' => $list['total'],
                'size'  => (int)$list['per_page'],
                'page'  => $list['current_page']
            ]
        ]);
    }
}
