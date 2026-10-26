<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\Order;
use App\Models\User;

class AdditionListController extends BaseController
{
    /**
     * 待添加客户列表
     * 
     *
     * @return array
     */
    public function index()
    {
        $size = request()->input('size', 10);
        $type = request()->input('type', 0);
        $memUid = request()->input('mem_uid', '');
        $haveMobile = request()->input('have_mobile');
        $isCopyMobile = request()->input('is_copy_mobile');
        $enterpriseId = request()->input('enterprise_id');
        if (!$memUid) {
            throw new Error(19112615, '员工memUid不能为空!');
        }

        // 已添加的客户
        $list = MemberUserRelation::query()
            ->where('mem_uid', $memUid)
            ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
            ->pluck('union_id')
            ->toArray();

        // join子查询
        $subQuery = Order::query()
            ->selectRaw('user_id, max(order_at) as order_at')
            ->whereIn('status', Order::STATUS_PAY_ARRAY)
            ->where('enterprise_id', $enterpriseId)
            ->groupBy('user_id')
            ->getQuery();

        // 未添加的客户
        $list = User::query()
            ->selectRaw('user.id, user.nick_name, user.avatar, user.type, user.mobile, guide_user_relation.is_copy_mobile, IFNULL(sub.order_at, 0) as order_at')
            ->join('guide_user_relation', function ($join) {
                $join->on([
                    ['user.id', '=', 'guide_user_relation.user_id']
                ]);
            })
            ->leftJoinSub($subQuery, 'sub', 'sub.user_id', '=', 'user.id')
            ->where('guide_user_relation.mem_uid', $memUid)
            ->whereNotIn('user.union_id', $list)
            ->where('guide_user_relation.user_id', '>', 0)
            ->where('guide_user_relation.is_copy_mobile', $isCopyMobile)
            ->where('guide_user_relation.deleted_at', 0)
            ->where('user.deleted_at', 0)
            ->when(isset($haveMobile), function($query) use ($haveMobile) {
                if ($haveMobile == 1) {
                    // 有手机号
                    $query->where('user.mobile', '<>', '');
                } else {
                    // 无手机号
                    $query->where('user.mobile', '');
                }
            })
            ->when($type > 0, function ($query) use ($type) {
                $query->where('user.type', $type);
            })
            ->orderBy('sub.order_at', 'desc')
            ->orderBy('user.id', 'desc')
            ->paginate($size)
            ->toArray();


        return success([
            'data' => $list['data'],
            'meta' => [
                'total' => $list['total'],
                'size' => (int)$list['per_page'],
                'page' => $list['current_page']
            ]
        ]);
    }
}
