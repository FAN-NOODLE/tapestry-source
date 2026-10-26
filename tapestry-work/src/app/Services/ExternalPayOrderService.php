<?php

namespace App\Services;

use App\Models\Scrm\MemberDepartments;
use App\Models\Scrm\Members;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use App\Models\ScrmEnterprise\ExternalPayRefundOrder;
use App\Services\Crm\CrmService;
use Exception;

class ExternalPayOrderService
{
    /**
     * @param $params
     *
     * @return mixed
     * 查导购信息
     */
    public function getMember($params)
    {
        return Members::query()->where([
            'enterprise_id' => $params['enterprise_id'],
            'mobile'        => $params['mobile']
        ])->select('id', 'name')->first();
    }

    /**
     * @param $params
     *
     * @return void
     * 订单列表
     */
    public function orderList($params): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;

        $query = ExternalPayOrder::query()
            ->setConn($params['enterprise_id'])
            ->where('member_id', $params['member_id'])
            ->orderBy('pay_time', 'desc');

        if (isset($params['order_state']) && $params['order_state']) {
            if ($params['order_state'] == 4) {
                $query->whereRaw("total_refund_fee=total_fee")->where("total_refund_fee", ">", 0);
            } else {
                $query->where('order_state', $params['order_state'])->whereRaw("total_fee>total_refund_fee");
            }

        }
        if (isset($params['sales_order_no']) && $params['sales_order_no']) {
            $query->where('sales_order_no', 'like', '%' . $params['sales_order_no'] . '%');
        }

        $query->select('id', 'pay_time', 'external_id', 'external_userid', 'total_fee', 'member_id', 'order_state', 'total_refund_fee', 'sales_order_no',
            'logistics_no', 'logistics_name', 'trade_state', 'delivery_type'
        );

        $data = $query->with('user')->paginate($size, ['*'], 'page', $page);

        foreach ($data as $row) {
            $row->avatar          = $row->user->avatar ?? '';
            $row->name            = $row->user->name ?? '';
            $row->nick            = $row->user->nick ?? '';
            $row->old_order_state = $row->order_state;//订单的原始状态
            $refund_state         = 0;//0为无退款,1为部分退款,2为全额退款
            if ($row->total_refund_fee > 0) {
                $refund_state     = $row->total_refund_fee == $row->total_fee ? 2 : 1;
                $row->order_state = $row->total_refund_fee == $row->total_fee ? 4 : $row->order_state;
            }
            $row->refund_state    = $refund_state;
            $row->order_state_str = ExternalPayOrder::ORDER_STATE[$row->order_state] ?? '';
            unset($row->user);
        }

        return [
            'list'  => $data->items(),
            'total' => $data->total(),
            'page'  => (int)$data->currentPage(),
            'size'  => (int)$data->perPage()
        ];
    }


    /**
     * 退款列表
     *
     * @param $params
     *
     * @return array
     */
    public function refundOrderList($params): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;

        $query = ExternalPayRefundOrder::query()
            ->setConn($params['enterprise_id'])
            ->where('refund_member_id', $params['member_id'])
            ->orderBy('refund_reqtime', 'desc');

        if (isset($params['sales_order_no']) && $params['sales_order_no']) {
            $query->whereHas('order', function ($q) use ($params) {
                $q->where('sales_order_no', 'like', '%' . $params['sales_order_no'] . '%');
            });
        }

        $query->select('id', 'out_trade_no', 'refund_member_id', 'external_userid', 'refund_comment', 'refund_status', 'refund_fee');

        $data = $query->with('order.user')->paginate($size, ['*'], 'page', $page);

        foreach ($data as $row) {
            $row->avatar            = $row->order->user->avatar ?? '';
            $row->name              = $row->order->user->name ?? '';
            $row->nick              = $row->order->user->nick ?? '';
            $row->pay_time          = $row->order->pay_time ?? '';
            $row->sales_order_no    = $row->order->sales_order_no ?? '';
            $row->refund_status_str = ExternalPayRefundOrder::REFUND_STATUS[$row->refund_status];
            unset($row->order);
        }

        return [
            'list'  => $data->items(),
            'total' => $data->total(),
            'page'  => (int)$data->currentPage(),
            'size'  => (int)$data->perPage()
        ];
    }

    /**
     * 发货
     *
     * @param $params
     *
     * @return mixed
     */
    public function deliverOrder($params): bool
    {
        $order = ExternalPayOrder::query()
            ->setConn($params['enterprise_id'])
            ->where(['id' => $params['id'], 'member_id' => $params['member_id']])
            ->first();

        if (!$order) {
            error("订单不存在");
        }

        if ($order->order_state == ExternalPayOrder::ORDER_STATE_SUCCESS) {
            error("订单已完成，不能再次发货");
        }

        if ($order->delivery_type == ExternalPayOrder::DELIVERY_TYPE_SELF) {
            //如果是自提订单就完成了
            $order->order_state = ExternalPayOrder::ORDER_STATE_SUCCESS;
        } else {
            $order->order_state = ExternalPayOrder::ORDER_STATE_DELIVER;
        }

        $order->logistics_no   = $params['logistics_no'];
        $order->logistics_name = $params['logistics_name'];
        $order->logistics_com  = $params['logistics_com'];
        $order->sales_order_no = $params['sales_order_no'];
        $order->delivery_time  = $params['delivery_time'];
        $order->tickets        = $params['tickets'];
        return $order->save();
    }

    /**
     * 订单详情
     *
     * @return void
     */
    public function orderDetail($params): array
    {
        $order = ExternalPayOrder::query()
            ->setConn($params['enterprise_id'])
            ->where(['id' => $params['id'], 'member_id' => $params['member_id']])
            ->with('user')
            ->first();

        if ($order) {
            $order->nick            = $order->user->nick ?? '';
            $order->unionid         = $order->user->unionid ?? '';
            $order->old_order_state = $order->order_state;//订单的原始状态
            if ($order->total_refund_fee == $order->total_fee) {
                $order->order_state = 4;
            }
            $order->order_state_str = ExternalPayOrder::ORDER_STATE[$order->order_state];
            unset($order->user);

            try {
                $union_code        = (new CrmService($params['enterprise_id']))->getMemberInfo([
                    'enterprise_id' => $params['enterprise_id'],
                    'union_id'      => $order->unionid
                ]);
                $order->union_code = $union_code->customer_code ?? '';
            } catch (Exception $exception) {
                $order->union_code = '';
            }

            $store             = MemberDepartments::query()
                ->where(['enterprise_id' => $params['enterprise_id']])
                ->where('member_id', $params['member_id'])
                ->with('departments')
                ->groupBy('member_id')
                ->first();
            $order->store_name = $store->departments->name ?? '';
        }

        return $order ? $order->toArray() : [];
    }

    /**
     * 退款详情
     *
     * @param $params
     *
     * @return array
     */
    public function refundOrderDetail($params)
    {
        $refund = ExternalPayRefundOrder::query()
            ->setConn($params['enterprise_id'])
            ->where(['id' => $params['id'], 'refund_member_id' => $params['member_id']])
            ->with('order.user')
            ->select('id', 'out_trade_no', 'out_refund_no', 'refund_reqtime', 'refund_status', 'refund_fee', 'refund_comment')
            ->first();

        if ($refund) {
            $refund->refund_status_str = ExternalPayRefundOrder::REFUND_STATUS[$refund->refund_status];
            $refund->avatar            = $refund->order->user->avatar ?? '';
            $refund->name              = $refund->order->user->name ?? '';
            $refund->nick              = $refund->order->user->nick ?? '';
            $refund->pay_time          = $refund->order->pay_time ?? '';
            $refund->total_fee         = $refund->order->total_fee ?? '';
            $refund->sales_order_no    = $refund->order->sales_order_no ?? '';
            $refund->unionid           = $refund->order->user->unionid ?? '';
            $refund->enterprise_id     = $params['enterprise_id'];

            try {
                $union_code         = (new CrmService($params['enterprise_id']))->getMemberInfo([
                    'enterprise_id' => $refund->enterprise_id,
                    'union_id'      => $refund->unionid
                ]);
                $refund->union_code = $union_code->customer_code ?? '';
            } catch (Exception $exception) {
                $refund->union_code = '';
            }
        }

        unset($refund->order);

        return $refund ? $refund->toArray() : [];
    }
}