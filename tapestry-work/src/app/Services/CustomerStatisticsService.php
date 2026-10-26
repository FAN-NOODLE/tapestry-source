<?php

namespace App\Services;

use App\Models\MemberUserRelation;
use App\Models\MemberUserRelationVisit;
use App\Library\Cdp\CdpService;
use Carbon\Carbon;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class CustomerStatisticsService
{
    /**
     * @var CdpService
     */
    protected $cdpService;

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        $this->cdpService = new CdpService();
    }

    /**
     * 客户统计数据
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
        // 添加客户总数(人)
        $totalUsers = MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->count();

        // 今日添加客户总数(人)
        $todayUsers = MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('relation_at', '>=', $dayStartTime)
            ->where('relation_at', '<=', $dayEndTime)
            ->count();

        // 昨日添加客户总数(人)
        $yesterdayUsers = MemberUserRelation::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('relation_at', '>=', $dayStartTime - 24 * 60 * 60)
            ->where('relation_at', '<=', $dayEndTime - 24 * 60 * 60)
            ->count();

        // 总共添加客户集合
        $totalUsersSet = MemberUserRelation::query()
            ->select('user.id as uid', 'member_user_relation.relation_at as tm')
            ->join('user', function ($join) {
                $join->on([
                    ['user.union_id', '=', 'member_user_relation.union_id'],
                ]);
            })
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('user.deleted_at', 0)
            ->where('member_user_relation.deleted_at', 0)
            ->get()
            ->toArray();

        // 今日访问人数
        $todayVisit = 0;

        if ($totalUsersSet) {
            Log::debug('今日添加客户集合: ' . json_encode($totalUsersSet));
            try {
                $todayVisitResult = $this->cdpService->userAct($totalUsersSet);

                Log::debug('调用cdp接口查询今日添加客户访问量: ' . json_encode($todayVisitResult));

                $todayVisit = ($todayVisitResult['count'] ?? 0) == -1 ? 0 : ($todayVisitResult['count'] ?? 0);
            } catch (\Exception $exception) {
                Log::error('调用cdp接口查询今日添加客户访问量: ' . json_encode($exception->getMessage()));
                $todayVisit = 0;
            }
        }

        // 昨日访问人数
        $yesterdayVisitData = MemberUserRelationVisit::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('result_time', '>=', $dayStartTime - 24 * 60 * 60)
            ->where('result_time', '<=', $dayEndTime - 24 * 60 * 60)
            ->first();

        return [
            'total_users'     => intval($totalUsers),
            'today_users'     => intval($todayUsers),
            'yesterday_users' => intval($yesterdayUsers),
            'today_visit'     => intval($todayVisit),
            'yesterday_visit' => intval($yesterdayVisitData->count ?? 0),
        ];
    }

    /**
     * 客户统计数据趋势图
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
        // 初始化添加客户统计数(人)
        $mapAddUserCounts = [];
        // 初始化客户浏览统计数(人)
        $mapVisitUserCounts = [];

        // 历史数据查询redis
        $key = 'customer_statistics_map_cache:' . $memUid . ':' . $endTime . ':' . $period;

        $statisticsMap = Redis::get($key) ?? '';

        $todayStartTime = Carbon::now()->startOfDay()->timestamp;
        $todayEndTime   = Carbon::now()->endOfDay()->timestamp;
        $today          = date('m-d', $todayStartTime);

        // 添加客户今日统计数(人)
        $todayAddUserCountsData = MemberUserRelation::query()
            ->select(DB::raw("count(union_id) map_add_user_count"))
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('relation_at', '>=', $todayStartTime)
            ->where('relation_at', '<=', $todayEndTime)
            ->first();

        // 总共添加客户集合
        $totalUsersSet = MemberUserRelation::query()
            ->select('user.id as uid', 'member_user_relation.relation_at as tm')
            ->join('user', function ($join) {
                $join->on([
                    ['user.union_id', '=', 'member_user_relation.union_id'],
                ]);
            })
            ->where('member_user_relation.enterprise_id', $enterpriseId)
            ->where('member_user_relation.mem_uid', $memUid)
            ->where('member_user_relation.delete_type', MemberUserRelation::NORMAL_STATUS)
            ->where('member_user_relation.is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
            ->where('user.deleted_at', 0)
            ->where('member_user_relation.deleted_at', 0)
            ->get()
            ->toArray();

        $todayVisit = 0;

        if ($totalUsersSet) {
            Log::debug('今日添加客户集合: ' . json_encode($totalUsersSet));
            try {
                $todayVisitResult = $this->cdpService->userAct($totalUsersSet);

                Log::debug('调用cdp接口查询今日添加客户访问量: ' . json_encode($todayVisitResult));

                $todayVisit = ($todayVisitResult['count'] ?? 0) == -1 ? 0 : ($todayVisitResult['count'] ?? 0);
            } catch (\Exception $exception) {
                Log::error('调用cdp接口查询今日添加客户访问量: ' . json_encode($exception->getMessage()));
                $todayVisit = 0;
            }
        }

        if ($statisticsMap) {
            $statisticsMap = json_decode($statisticsMap, true);

            $mapAddUserCounts   = $statisticsMap['map_add_user_counts'];
            $mapVisitUserCounts = $statisticsMap['map_visit_user_counts'];

            $mapAddUserCounts[]   = ['day' => $today, 'value' => $todayAddUserCountsData->map_add_user_count ?? 0];
            $mapVisitUserCounts[] = ['day' => $today, 'value' => $todayVisit];
        } else {
            // 添加客户按时间统计数(人)
            $mapAddUserCountsData = MemberUserRelation::query()
                ->select(DB::raw("count(union_id) map_add_user_count, FROM_UNIXTIME(relation_at, '%m-%d') relation_time"))
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $memUid)
                ->where('delete_type', MemberUserRelation::NORMAL_STATUS)
                ->where('is_binding_customer', MemberUserRelation::BINDING_CUSTOMER_YES)
                ->where('relation_at', '>=', $startTime)
                ->where('relation_at', '<=', $endTime)
                ->groupBy('relation_time')
                ->get()
                ->toArray();

            $mapAddUserCountsData = array_column($mapAddUserCountsData, null, 'relation_time');

            Log::debug('添加客户按时间统计数: ' . json_encode($mapAddUserCountsData));

            // 客户浏览按时间统计数(人)
            $mapVisitUserCountsData = MemberUserRelationVisit::query()
                ->select(DB::raw("count as map_visit_user_count, FROM_UNIXTIME(result_time, '%m-%d') result_time"))
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $memUid)
                ->where('result_time', '>=', $startTime)
                ->where('result_time', '<=', $endTime)
                ->get()
                ->toArray();

            $mapVisitUserCountsData = array_column($mapVisitUserCountsData, null, 'result_time');

            Log::debug('客户浏览按时间统计数: ' . json_encode($mapVisitUserCountsData));

            foreach ($dayArr as $day) {
                $mapAddUserCounts[]   = ['day' => $day, 'value' => ($mapAddUserCountsData[$day] ?? [])['map_add_user_count'] ?? 0];
                $mapVisitUserCounts[] = ['day' => $day, 'value' => ($mapVisitUserCountsData[$day] ?? [])['map_visit_user_count'] ?? 0];
            }

            Redis::set($key, json_encode([
                'map_add_user_counts'   => $mapAddUserCounts,
                'map_visit_user_counts' => $mapVisitUserCounts,
            ]));

            Redis::expire($key, 25 * 60 * 60);

            $mapAddUserCounts[]   = ['day' => $today, 'value' => $todayAddUserCountsData->map_add_user_count ?? 0];
            $mapVisitUserCounts[] = ['day' => $today, 'value' => $todayVisit];
        }

        $data = [
            'map_add_user_counts'   => $mapAddUserCounts,
            'map_visit_user_counts' => $mapVisitUserCounts,
        ];

        return $data;
    }
}
