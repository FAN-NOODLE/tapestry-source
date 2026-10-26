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
use App\Models\MemberUserRelation;
use App\Models\Order;
use App\Models\User;

class CountController extends BaseController
{
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 首页 导购数据统计
     *
     * @return array
     */
    public function index()
    {
        $additionTotal = $this->additionTotal();
        $orderTotal = $this->orderTotal();

        return success(['data' => [
            'addition_total' => $additionTotal,
            'order_total' => $orderTotal
        ]]);
    }

    /**
     * 待添加客户数量
     *
     * @return int
     */
    protected function additionTotal()
    {
        // 已添加的客户
        $list = MemberUserRelation::query()
            ->where('mem_uid', $this->memUid)
            ->where('delete_type', '<>', MemberUserRelation::DEL_EXTERNAL_CONTACT)
            ->pluck('union_id')
            ->toArray();

        // 未添加的客户
        $total = User::query()
            ->join('guide_user_relation', function ($join) {
                $join->on([
                    ['user.id', '=', 'guide_user_relation.user_id']
                ]);
            })
            ->where('guide_user_relation.mem_uid', $this->memUid)
            ->where('guide_user_relation.user_id', '>', 0)
            ->whereNotIn('guide_user_relation.union_id', $list)
            ->where('user.mobile', '<>', '')
            ->where('user.deleted_at', 0)
            ->where('guide_user_relation.deleted_at', 0)
            ->count();

        return $total;
    }

    /**
     * 业绩订单未读数
     *
     * @return int
     */
    protected function orderTotal()
    {
        // 已添加的客户未读订单数
        $total = MemberUserRelation::query()
            ->join('order', function ($join) {
                $join->on([
                    ['order.user_id', '=', 'member_user_relation.user_id']
                ]);
            })
            ->where('member_user_relation.delete_type', '<>', MemberUserRelation::DEL_EXTERNAL_CONTACT)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('member_user_relation.mem_uid', $this->memUid)
            ->where('order.read_status', Order::READ_STATUS_NO)
            ->whereIn('order.status', Order::STATUS_PAY_ARRAY)
            ->whereRaw('order.order_at >= member_user_relation.relation_at')
            ->where('member_user_relation.deleted_at', 0)
            ->where('order.deleted_at', 0)
            ->count();

        return $total;
    }
}
