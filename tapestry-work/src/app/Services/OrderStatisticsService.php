<?php

namespace App\Services;

use App\Models\MemberUserRelation;
use App\Models\Order;
use Carbon\Carbon;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Redis;

class OrderStatisticsService
{
    /**
     * 订单统计数据
     *
     * @param $enterpriseId
     * @param $memUid
     * @param $dayStartTime
     * @param $dayEndTime
     *
     * @return array
     */
    public function statistics($enterpriseId, $memUid, $dayStartTime, $dayEndTime)
    {
        // 成交总额(元)，成交订单总数(笔)
        $totalOrders = MemberUserRelation::query()
            ->select(DB::raw('sum(order.final_amount) total_amount, count(order.order_no) total_count'))
            ->join('user', 'member_user_relation.union_id', 'user.union_id')
            ->leftJoin('order', 'user.id', 'order.user_id')
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->whereIn('order.status', Order::STATUS_TRADED)
            ->where('member_user_relation.deleted_at', 0)
            ->where('user.deleted_at', 0)
            ->where('order.deleted_at', 0)
            ->whereRaw('order.order_at >= member_user_relation.relation_at')
            ->first();

        // 今日成交总额(元)，今日成交订单总数(笔)
        $todayOrders = MemberUserRelation::query()
            ->select(DB::raw('sum(order.final_amount) today_amount, count(order.order_no) today_count'))
            ->join('user', 'member_user_relation.union_id', 'user.union_id')
            ->leftJoin('order', 'user.id', 'order.user_id')
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->whereIn('order.status', Order::STATUS_TRADED)
            ->where('member_user_relation.deleted_at', 0)
            ->where('user.deleted_at', 0)
            ->where('order.deleted_at', 0)
            ->whereRaw('order.order_at >= member_user_relation.relation_at')
            ->where('order.order_at', '>=', $dayStartTime)
            ->where('order.order_at', '<=', $dayEndTime)
            ->first();

        // 昨日成交总额(元)，昨日成交订单总数(笔)
        $yesterdayOrders = MemberUserRelation::query()
            ->select(DB::raw('sum(order.final_amount) yesterday_amount, count(order.order_no) yesterday_count'))
            ->join('user', 'member_user_relation.union_id', 'user.union_id')
            ->leftJoin('order', 'user.id', 'order.user_id')
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->whereIn('order.status', Order::STATUS_TRADED)
            ->where('member_user_relation.deleted_at', 0)
            ->where('user.deleted_at', 0)
            ->where('order.deleted_at', 0)
            ->whereRaw('order.order_at >= member_user_relation.relation_at')
            ->where('order.order_at', '>=', $dayStartTime - 24 * 60 * 60)
            ->where('order.order_at', '<=', $dayEndTime - 24 * 60 * 60)
            ->first();

        return [
            'total_amount'     => intval($totalOrders->total_amount ?? 0),
            'total_count'      => intval($totalOrders->total_count ?? 0),
            'today_amount'     => intval($todayOrders->today_amount ?? 0),
            'today_count'      => intval($todayOrders->today_count ?? 0),
            'yesterday_amount' => intval($yesterdayOrders->yesterday_amount ?? 0),
            'yesterday_count'  => intval($yesterdayOrders->yesterday_count ?? 0),
        ];
    }

    /**
     * 订单统计数据趋势图
     *
     * @param $enterpriseId
     * @param $memUid
     * @param $startTime
     * @param $endTime
     * @param $dayArr
     * @param $period
     *
     * @return array
     */
    public function statisticsMap($enterpriseId, $memUid, $startTime, $endTime, $dayArr, $period)
    {
        // 初始化成交总额(元)
        $mapOrderAmounts = [];
        // 初始化成交订单笔数(笔)
        $mapOrderCounts = [];

        // 历史数据查询redis
        $key = 'order_statistics_map_cache:' . $memUid . ':' . $endTime . ':' . $period;

        $statisticsMap = Redis::get($key) ?? '';

        $todayStartTime = Carbon::now()->startOfDay()->timestamp;
        $todayEndTime   = Carbon::now()->endOfDay()->timestamp;
        $today          = date('m-d', $todayStartTime);

        // 成交总额今日统计数(元)，成交订单今日统计数(笔)
        $todayOrderCountsData = MemberUserRelation::query()
            ->select(DB::raw("sum(order.final_amount) map_order_amount, count(order.order_no) map_order_count"))
            ->join('user', 'member_user_relation.union_id', 'user.union_id')
            ->leftJoin('order', 'user.id', 'order.user_id')
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->whereIn('order.status', Order::STATUS_TRADED)
            ->where('member_user_relation.deleted_at', 0)
            ->where('user.deleted_at', 0)
            ->where('order.deleted_at', 0)
            ->where('order.order_at', '>=', $todayStartTime)
            ->where('order.order_at', '<=', $todayEndTime)
            ->whereRaw('order.order_at >= member_user_relation.relation_at')
            ->first();

        if ($statisticsMap) {
            $statisticsMap = json_decode($statisticsMap, true);

            $mapOrderAmounts = $statisticsMap['map_order_amounts'];
            $mapOrderCounts  = $statisticsMap['map_order_counts'];

            $mapOrderAmounts[] = ['day' => $today, 'value' => $todayOrderCountsData->map_order_amount ?? 0];
            $mapOrderCounts[]  = ['day' => $today, 'value' => $todayOrderCountsData->map_order_count ?? 0];
        } else {
            // 成交总额按时间统计数(元)，成交订单按时间统计数(笔)
            $mapOrderCountsData = MemberUserRelation::query()
                ->select(DB::raw("sum(order.final_amount) map_order_amount, count(order.order_no) map_order_count, FROM_UNIXTIME(order.order_at, '%m-%d') order_time"))
                ->join('user', 'member_user_relation.union_id', 'user.union_id')
                ->leftJoin('order', 'user.id', 'order.user_id')
                ->where('member_user_relation.enterprise_id', $enterpriseId)
                ->where('member_user_relation.mem_uid', $memUid)
                ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
                ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
                ->whereIn('order.status', Order::STATUS_TRADED)
                ->where('member_user_relation.deleted_at', 0)
                ->where('user.deleted_at', 0)
                ->where('order.deleted_at', 0)
                ->where('order.order_at', '>=', $startTime)
                ->where('order.order_at', '<=', $endTime)
                ->whereRaw('order.order_at >= member_user_relation.relation_at')
                ->groupBy('order_time')
                ->get()
                ->toArray();

            $mapOrderCountsData = array_column($mapOrderCountsData, null, 'order_time');

            foreach ($dayArr as $day) {
                $mapOrderAmounts[] = ['day' => $day, 'value' => ($mapOrderCountsData[$day] ?? [])['map_order_amount'] ?? 0];
                $mapOrderCounts[]  = ['day' => $day, 'value' => ($mapOrderCountsData[$day] ?? [])['map_order_count'] ?? 0];
            }

            Redis::set($key, json_encode([
                'map_order_amounts' => $mapOrderAmounts,
                'map_order_counts'  => $mapOrderCounts,
            ]));

            Redis::expire($key, 25 * 60 * 60);

            $mapOrderAmounts[] = ['day' => $today, 'value' => $todayOrderCountsData->map_order_amount ?? 0];
            $mapOrderCounts[]  = ['day' => $today, 'value' => $todayOrderCountsData->map_order_count ?? 0];
        }

        $data = [
            'map_order_amounts' => $mapOrderAmounts,
            'map_order_counts'  => $mapOrderCounts,
        ];

        return $data;
    }
}
