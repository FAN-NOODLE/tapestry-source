<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\Order;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;

class AchievementListController extends BaseController
{
    // 业绩统计用户列表
    public function index(Request $request)
    {
        $limit     = $request->input('size', 10);
        $orderType = $request->input('order_type', 0);

        // 测试使用
        //$this->enterpriseId = 2;
        //$this->memUid       = '2C8B36927F0000017AC9A7CF9045A1EF';

        //DB::enableQueryLog();
        //  userIds获取
        $joinSub = (new MemberUserRelation())
            ->setTable('mu')
            ->from('member_user_relation as mu')
            ->select([
                'mu.user_id as id',
                'mu.relation_at',
            ])
            ->where([
                'mu.enterprise_id'       => $this->enterpriseId,
                'mu.mem_uid'             => $this->memUid,
                'mu.is_binding_customer' => MemberUserRelation::BINDING_CUSTOMER_YES,
                'mu.deleted_at'          => 0,
            ])
            ->where('mu.delete_type', '<>', MemberUserRelation::DEL_EXTERNAL_CONTACT)
            ->where('mu.user_id', '>', 0)
            ->getQuery();

        $orders = (new Order())
            ->setTable('ord')
            ->from('order as ord')
            ->with([
                'goods' => function ($query) {
                    $query->select([
                        "order_no",
                        "spu_id",
                        "sku_id",
                        "spu_name",
                        "category_name",
                        "sku_picture",
                        "sku_attribute",
                        "quantity",
                        "sell_price",
                        "sub_price"
                    ]);
                }
            ])
            ->with([
                'user' => function ($query) {
                    $query->select([
                        'id',
                        'sex',
                        'nick_name',
                        'real_name',
                        'mobile',
                        'avatar',
                        'type',
                    ]);
                }
            ])
            ->joinSub($joinSub, 'son', 'son.id', '=', 'ord.user_id')
            ->orderBy('ord.order_at', 'desc')
            ->select([
                "ord.id",
                "ord.user_id",
                "ord.order_no",
                "ord.cash_fee as final_amount",
                "ord.order_at",
                "ord.cash_fee",
                "ord.pay_at",
                "ord.status",
            ])
            ->when($orderType, function ($query) use ($orderType) {
                ($orderType == order::STATUS_WAIT) && $query->whereIn('status', order::STATUS_WAIT_ARRAY);
                ($orderType == order::STATUS_PAY) && $query->whereIn('status', order::STATUS_PAY_ARRAY);
                ($orderType == order::STATUS_FINISH) && $query->whereIn('status', order::STATUS_FINISH_ARRAY);
            })
            ->whereRaw('ord.order_at >= son.relation_at')
            ->where([
                'ord.enterprise_id' => $this->enterpriseId,
            ])
            ->paginate($limit)
            ->toArray();

        //print_r(DB::getQueryLog());

        foreach ($orders['data'] as &$order) {
            // 状态别名返回
            $order['status_alias'] = order::getStatusAlias($order['status']);
            foreach ($order['goods'] as &$item) {
                $item['category_name'] = json_decode($item['category_name'] ?? '[]', true);
                $item['sku_attribute'] = json_decode($item['sku_attribute'] ?? '[]', true);
            }
        }

        return success([
            'data' => $orders['data'],
            'meta' => [
                'total' => $orders['total'],
                'size'  => (int)$orders['per_page'],
                'page'  => $orders['current_page']
            ]
        ]);
    }

}
