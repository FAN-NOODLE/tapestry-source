<?php

namespace App\Services;

use App\Models\Scrm\MemberDepartments;
use App\Models\ScrmEnterprise\ExternalPayOrder;
use App\Models\ScrmEnterprise\ExternalPayRefundOrder;
use App\Models\ScrmEnterprise\WxUserCorpExternal;


class ExternalPayUserOrderService
{
    /**
     * 查询用户信息
     *
     * @param $params
     *
     * @return mixed
     */
    public function getUser($params)
    {
        return WxUserCorpExternal::getExternalUserInfo($params['enterprise_id'], $params['union_id']);
    }

    /**
     * 订单列表
     *
     * @return void
     */
    public function orderList($params): array
    {
        $page = $params['page'] ?? 1;
        $size = $params['size'] ?? 10;

        $query = ExternalPayOrder::query()
            ->setConn($params['enterprise_id'])
            ->where('external_userid', $params['external_userid'])
            ->orderBy('pay_time', 'desc');

        if (isset($params['order_state']) && $params['order_state']) {
            $query->where('order_state', $params['order_state']);
        }

        $query->select('id', 'pay_time', 'member_id', 'order_state', 'total_fee', 'sales_order_no',
            'logistics_no', 'logistics_name', 'trade_state', 'tickets', 'contact_info', 'total_refund_fee', 'delivery_type'
        );

        $data = $query->with('member')->paginate($size, ['*'], 'page', $page);

        //导购id
        $member_id = collect($data->items())->pluck('member_id')->all();

        //查询店铺
        $departments = $this->getStore($params, $member_id);

        foreach ($data as $row) {
            $row->member_name     = $row->member->name ?? '';
            $row->order_state_str = ExternalPayOrder::ORDER_STATE[$row->order_state];
            $row->store_name      = $departments[$row->member_id]['departments']['name'] ?? '';
            unset($row->member);
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
            ->where('external_userid', $params['external_userid'])
            ->where('refund_status', 2)
            ->orderBy('refund_reqtime', 'desc');

        $query->select('id', 'refund_member_id', 'out_trade_no', 'refund_status', 'refund_fee', 'refund_time', 'refund_reqtime');

        $data = $query->with('member')->paginate($size, ['*'], 'page', $page);

        //导购id
        $member_id = collect($data->items())->pluck('refund_member_id')->all();
        //查询店铺
        $departments = $this->getStore($params, $member_id);

        foreach ($data as $row) {
            $row->refund_time       = $row->refund_reqtime ?? 0;
            $row->member_name       = $row->member->name ?? '';
            $row->refund_status_str = ExternalPayRefundOrder::REFUND_STATUS[$row->refund_status];
            $row->store_name        = $departments[$row->refund_member_id]['departments']['name'] ?? '';
            unset($row->member);
        }

        return [
            'list'  => $data->items(),
            'total' => $data->total(),
            'page'  => (int)$data->currentPage(),
            'size'  => (int)$data->perPage()
        ];
    }

    /**
     * @return void
     */
    public function getStore($params, $member_id): array
    {
        return MemberDepartments::query()
            ->where(['enterprise_id' => $params['enterprise_id']])
            ->whereIn('member_id', $member_id)
            ->with('departments')
            ->groupBy('member_id')
            ->get()
            ->keyBy('member_id')
            ->toArray();
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
            ->where(['id' => $params['id'], 'external_userid' => $params['external_userid']])
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
            $order->union_code = '';
            //导购
            $order->sa_name = $order->member->name ?? '--';
            //门店
            $store             = MemberDepartments::query()
                ->where(['enterprise_id' => $params['enterprise_id']])
                ->where('member_id', $order['member_id'])
                ->with('departments')
                ->groupBy('member_id')
                ->first();
            $order->store_name = $store->departments->name ?? '';

            unset($order->member);
            unset($order->user);
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
            ->where(['id' => $params['id'], 'external_userid' => $params['external_userid']])
            ->with('order.user')
            ->select('id', 'out_trade_no', 'out_refund_no', 'refund_reqtime', 'refund_status', 'refund_fee', 'refund_comment','refund_member_id')
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
            $refund->union_code        = '';
            $refund->sa_name           = $refund->member->name ?? '--';
            //门店
            $store = MemberDepartments::query()
                ->where(['enterprise_id' => $params['enterprise_id']])
                ->where('member_id', $refund['refund_member_id'])
                ->with('departments')
                ->groupBy('member_id')
                ->first();
            $refund->store_name = $store->departments->name ?? '';
        }

        unset($refund->order);
        unset($refund->member);

        return $refund ? $refund->toArray() : [];
    }
}