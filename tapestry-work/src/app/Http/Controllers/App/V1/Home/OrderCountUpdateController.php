<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/11/29
 * Time: 17:03
 */

namespace App\Http\Controllers\App\V1\Home;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\GuideUserRelation;
use App\Models\Order;

class OrderCountUpdateController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 业绩订单标为已读
     *
     * @return array
     */
    public function index()
    {

        // 已添加的客户id
        $userIds = GuideUserRelation::query()
            ->join('member_user_relation', function ($join) {
                $join->on([
                    ['guide_user_relation.mem_uid', '=', 'member_user_relation.mem_uid'],
                    ['guide_user_relation.union_id', '=', 'member_user_relation.union_id']
                ]);
            })
            ->where('guide_user_relation.mem_uid', $this->memUid)
            ->pluck('guide_user_relation.user_id')
            ->toArray();

        Order::query()->whereIn('user_id', $userIds)
            ->where('read_status', Order::READ_STATUS_NO)
            ->whereIn('status', Order::STATUS_PAY_ARRAY)
            ->update(['read_status' => Order::READ_STATUS_YES]);

        return success();
    }
}
