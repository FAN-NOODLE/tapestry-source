<?php
namespace App\Console\Commands\Temp;

use App\Common\DateUtils;
use App\Models\DailyMemberCustomerStatistic;
use App\Models\DailyWorkerCustomerStatistic;
use App\Models\MemberUserRelation;
use App\Models\StoreWorker;
use App\Services\StoreWorkerService;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;
use App\Common\Constant;

class MigrateHistoryCustomerStatistic extends Command
{
    /**
     * @var string
     */
    protected $signature = 'migrate-history-customer-statistic {--enterpriseId=} {--memUid=} {--workerId=} {--endStatTime=}';

    /**
     * @var string
     */
    protected $description = '导购历史顾客统计数据迁移';


    public function handle()
    {
        $this->info("MigrateHistoryCustomerStatistic start");

        $enterpriseId = intval($this->option('enterpriseId'));
        $memUid = strval($this->option('memUid'));
        $workerId = intval($this->option('workerId'));
        $endStatTime = intval($this->option('endStatTime')) ?? DateUtils::todayStatTime();
        $endStatTime = $endStatTime > 0 ? $endStatTime : DateUtils::todayStatTime();

        $where = [];
        if ($enterpriseId > 0) {
            $where['enterprise_id'] = $enterpriseId;
        }
        if ($memUid != '') {
            $where['mem_uids'] = [$memUid];
        }
        if ($workerId > 0) {
            $where['ids'] = [$workerId];
        }
        $columns = ['id', 'enterprise_id', 'mem_uid', 'contact_qr_code_id', 'worker_status', 'created_at', 'updated_at'];
        $workerList = StoreWorkerService::list($where, false, $columns)['list'];
        if (empty($workerList)) {
            $this->info('workerList is empty');
            return;
        }

        $memUidWorkerMap = collect($workerList)->sortBy('id')->sortBy('enterprise_id')->groupBy('mem_uid')->toArray();
        $memUidList = array_keys($memUidWorkerMap);
        sort($memUidList);

        $memUidCount = count($memUidList);
        $failMemUidList = [];

        $redisKey = md5($this->signature) . '-' . md5(json_encode($this->options())) . "-" . DateUtils::todayStatTime();
        $this->info("redisKey: " . $redisKey);
        while (true) {
            $incrNum = Redis::incr($redisKey);
            if ($incrNum > $memUidCount) {
                $this->info('all member migrate finished');
                break;
            }

            $this->info(sprintf("migrate by mem_uid, total: %d, curr: %d", $memUidCount, $incrNum));
            $memUid = $memUidList[$incrNum - 1];
            $workerList = $memUidWorkerMap[$memUid];
            $enterpriseId = $workerList[0]['enterprise_id'];
            try {
                $this->migrateMemberCustomerStatistic($enterpriseId, $memUid, $workerList, $endStatTime);
            } catch (\Exception $e) {
                Log::info("migrateMemberCustomerStatistic exception", [
                    'memUid' => $memUid,
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
                $this->info("migrateMemberCustomerStatistic exception memUid:" . $memUid);
                array_push($failMemUidList, $memUid);
            }
        }

        $this->info("MigrateHistoryCustomerStatistic end, failMemUidList: " . json_encode($failMemUidList));
    }

    public function migrateMemberCustomerStatistic(int $enterpriseId, string $memUid, array $workerList, int $endStatTime)
    {
        $this->info(sprintf("migrateMemberCustomerStatistic memUid:%s start", $memUid));

        $now = time();

        $contactList = MemberUserRelation::query()
            ->select(['source_type', 'created_at'])
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('delete_type', 0)
            ->get()
            ->toArray();
        $contactDateMap = [];
        foreach ($contactList as $contact) {
            $statTime = DateUtils::timeToStatTime($contact['created_at']);
            if (!isset($contactDateMap[$statTime])) {
                $contactDateMap[$statTime] = [];
            }
            $contactDateMap[$statTime][] = $contact['source_type'];
        }

        list($contactMemberList, $contactMemberDateMap) = $this->getData($enterpriseId, $memUid, Constant::IS_MEMBER_STATISTICS);
        list($contactNonMemberList, $contactNonMemberDateMap) = $this->getData($enterpriseId, $memUid, Constant::NON_MEMBER_STATISTICS);

        $minTime = $contactDateMap ? min(array_keys($contactDateMap)) : 0;

        $createdAtList = array_column($workerList, 'created_at');
        sort($createdAtList);
        $memberStartStatTime = DateUtils::timeToStatTime($createdAtList[0] > $minTime && $minTime != 0 ? $minTime : $createdAtList[0]);

        // 员工维度每日的顾客统计数量
        $dateMemberCustomerNumDataList = [];
        $cumulativeCustomerNum = 0;
        $cumulativeMemberNum = 0;
        $cumulativeNonMemberNum = 0;
        $allTimeRange = [];
        for ($statTime = $memberStartStatTime; $statTime <= $endStatTime; $statTime += DateUtils::DAY_SECONDS) {
            array_push($allTimeRange, $statTime);

            $todayCustomerNum = isset($contactDateMap[$statTime]) ? count($contactDateMap[$statTime]) : 0;
            $cumulativeCustomerNum += $todayCustomerNum;

            $todayMemberNum = isset($contactMemberDateMap[$statTime]) ? count($contactMemberDateMap[$statTime]) : 0;
            $cumulativeMemberNum += $todayMemberNum;

            $todayNonMemberNum = isset($contactNonMemberDateMap[$statTime]) ? count($contactNonMemberDateMap[$statTime]) : 0;
            $cumulativeNonMemberNum += $todayNonMemberNum;

            array_push($dateMemberCustomerNumDataList, [
                'enterprise_id' => $enterpriseId,
                'mem_uid' => $memUid,
                'stat_time' => $statTime,
                'cumulative_customer_num' => $cumulativeCustomerNum,
                'today_customer_num' => $todayCustomerNum,
                'cumulative_member_num' => $cumulativeMemberNum,
                'today_member_num' => $todayMemberNum,
                'cumulative_non_member_num' => $cumulativeNonMemberNum,
                'today_non_member_num' => $todayNonMemberNum,
                'created_at' => $now,
                'updated_at' => $now
            ]);
        }
        // 导购维度每日的顾客统计数量
        $dateWorkerCustomerNumDataList = [];
        foreach ($workerList as $worker) {
            $cumulativeContactCardWayCustomerNum = 0;
            $cumulativeOtherWayCustomerNum = 0;
            //$workerStartStatTime = DateUtils::timeToStatTime($worker['created_at']);
            $workerStartStatTime = DateUtils::timeToStatTime($worker['created_at'] > $minTime && $minTime != 0 ? $minTime : $worker['created_at']);
            if (StoreWorker::IN_SERVICE == $worker['worker_status']) {
                $leaveNextStatTime = DateUtils::todayStatTime() + DateUtils::DAY_SECONDS;
            } else {
                $leaveNextStatTime = DateUtils::timeToStatTime($worker['updated_at']) + DateUtils::DAY_SECONDS;
            }
            for ($statTime = $workerStartStatTime; $statTime <= $endStatTime; $statTime += DateUtils::DAY_SECONDS) {
                $todayContactCardWayCustomerNum = 0;
                $todayOtherWayCustomerNum = 0;
                if ($statTime < $leaveNextStatTime && isset($contactDateMap[$statTime])) {
                    foreach ($contactDateMap[$statTime] as $sourceType) {
                        if ($sourceType == $worker['contact_qr_code_id']) {
                            $todayContactCardWayCustomerNum++;
                        } else {
                            $todayOtherWayCustomerNum++;
                        }
                    }
                }
                $cumulativeContactCardWayCustomerNum += $todayContactCardWayCustomerNum;
                $cumulativeOtherWayCustomerNum += $todayOtherWayCustomerNum;

                array_push($dateWorkerCustomerNumDataList, [
                    'enterprise_id' => $enterpriseId,
                    'worker_id' => $worker['id'],
                    'mem_uid' => $memUid,
                    'stat_time' => $statTime,
                    'cumulative_contact_card_way_customer_num' => $cumulativeContactCardWayCustomerNum,
                    'cumulative_other_way_customer_num' => $cumulativeOtherWayCustomerNum,
                    'today_contact_card_way_customer_num' => $todayContactCardWayCustomerNum,
                    'today_other_way_customer_num' => $todayOtherWayCustomerNum,
                    'created_at' => $now,
                    'updated_at' => $now
                ]);
            }
        }
        $limit = 500;
        $pageTotal = ceil(count($dateWorkerCustomerNumDataList) / $limit);
        for ($page = 1; $page <= $pageTotal; $page++) {
            $offset = ($page - 1) * $limit;
            DailyWorkerCustomerStatistic::query()->insert(array_slice($dateWorkerCustomerNumDataList, $offset, $limit));
        }

        $findStatisticsDataCurMemUid = DailyMemberCustomerStatistic::query()
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('stat_time', $allTimeRange)
            ->where('mem_uid', $memUid)
            ->pluck('stat_time')
            ->toArray();

        $updateData = $insertData = [];
        foreach ($dateMemberCustomerNumDataList as $key => $statTime) {
            if(in_array($statTime['stat_time'], $findStatisticsDataCurMemUid)){
                array_push($updateData, $dateMemberCustomerNumDataList[$key]);
            }else{
                array_push($insertData, $dateMemberCustomerNumDataList[$key]);
            }
        }

        // 更新当前导购的统计数据
        if($updateData){
            foreach ($updateData as $value){
                DailyMemberCustomerStatistic::query()
                    ->where([
                        'enterprise_id' => $value['enterprise_id'],
                        'stat_time' => $value['stat_time'],
                        'mem_uid' => $value['mem_uid']
                    ])
                    ->update($value);
            }
        }

        // 若没有导购某个日期的统计数据则添加
        $insertData && DailyMemberCustomerStatistic::query()->insert($insertData);

        unset($findStatisticsDataCurMemUid, $updateData, $insertData);
        unset($contactList, $contactDateMap, $dateMemberCustomerNumDataList, $dateWorkerCustomerNumDataList, $contactMemberList, $contactMemberDateMap, $contactNonMemberList, $contactNonMemberDateMap);

        $this->info(sprintf("migrateMemberCustomerStatistic memUid:%d succeed", $memUid));
    }

    /**
     * @param $enterpriseId
     * @param $memUid
     * @param int $isMember
     * @return array
     * @Date: 2021/3/29
     * @Time: 16:40
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function getData($enterpriseId, $memUid, $isMember = Constant::IS_MEMBER_STATISTICS)
    {
        $contactList = MemberUserRelation::query()
            ->select(['source_type', 'created_at'])
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->when(!empty($isMember), function ($query) use ($isMember) {
                if ($isMember == Constant::IS_MEMBER_STATISTICS) {
                    $query->where('union_code', '<>', '');
                } else if ($isMember == Constant::NON_MEMBER_STATISTICS) {
                    $query->where('union_code', '');
                }
            })
            ->where('delete_type', 0)
            ->get()
            ->toArray();
        $contactDateMap = [];
        foreach ($contactList as $contact) {
            $statTime = DateUtils::timeToStatTime($contact['created_at']);
            if (!isset($contactDateMap[$statTime])) {
                $contactDateMap[$statTime] = [];
            }
            $contactDateMap[$statTime][] = $contact['source_type'];
        }

        return [$contactList, $contactDateMap];
    }
}