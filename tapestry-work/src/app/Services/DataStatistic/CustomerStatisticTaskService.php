<?php


namespace App\Services\DataStatistic;


use App\Common\DateUtils;
use App\Exceptions\Error;
use App\Library\Uc\Member;
use App\Models\DailyMemberCustomerStatistic;
use App\Models\DailyWorkerCustomerStatistic;
use App\Models\StoreWorker;
use App\Services\BaseService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Facades\Log;

class CustomerStatisticTaskService extends BaseService
{

    /**
     * 初始化所有企业指定日期顾客统计数据
     * @param $statTime
     * @param bool $force
     */
    public function initAllEnterpriseCustomerStatistic($statTime, bool $force)
    {
        Log::info("initAllEnterpriseCustomerStatistic start", ['statTime' => $statTime]);
        $enterpriseIdList = StoreWorkerService::getAllEnterpriseId();
        if (empty($enterpriseIdList)) {
            Log::info("initAllEnterpriseCustomerStatistic enterpriseIdList is empty");
            return;
        }

        foreach ($enterpriseIdList as $enterpriseId) {
            $this->initEnterpriseCustomerStatistic($enterpriseId, $statTime, $force);
        }
        Log::info("initAllEnterpriseCustomerStatistic end");
    }

    /**
     * 初始化指定企业指定日期顾客统计数据
     * @param int $enterpriseId
     * @param int $statTime
     * @param bool $force
     */
    public function initEnterpriseCustomerStatistic(int $enterpriseId, int $statTime, bool $force)
    {
        Log::info("initEnterpriseCustomerStatistic params", func_get_args());
        $flagKey = redisKey('ENTERPRISE_CUSTOMER_STATISTIC_INIT_FLAG', [
            'enterprise_id' => $enterpriseId,
            'stat_time' => $statTime
        ]);
        // 判断今日是否已经执行过
        if (!$force && getOrSetData($flagKey)) {
            Log::info("initEnterpriseCustomerStatistic initialized");
            return;
        }

        try {
            // 获取在职的导购列表
            $workerList = StoreWorker::query()
                ->select(['id', 'mem_uid'])
                ->where('enterprise_id', $enterpriseId)
                ->get()
                ->toArray();
            if (empty($workerList)) {
                Log::info("initEnterpriseCustomerStatistic workerList is empty");
                return;
            }
            $memUidList = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));

            // 查询已经创建数据的导购，不再进行重复创建
            $workerStatisticList = DailyWorkerCustomerStatistic::listStatisticByStatTime($enterpriseId, $statTime, [],
                ['worker_id']);
            $workerStatisticMap = array_flip(array_column($workerStatisticList, 'worker_id'));

            $memberStatisticList = DailyMemberCustomerStatistic::listStatisticByStatTime($enterpriseId, $statTime, [],
                ['mem_uid']);
            $memberStatisticMap = array_flip(array_column($memberStatisticList, 'mem_uid'));

            $createdAt = time();
            $workerStatisticData = [];
            foreach ($workerList as $worker) {
                if (!isset($workerStatisticMap[$worker['id']])) {
                    $workerStatisticData[] = DailyWorkerCustomerStatistic::generateWorkerStatisticInitData($enterpriseId,
                        $worker['id'], $worker['mem_uid'], $statTime, $createdAt);
                }
            }
            $memberStatisticData = [];
            foreach ($memUidList as $memUid) {
                if (!isset($memberStatisticMap[$memUid])) {
                    $memberStatisticData[] = DailyMemberCustomerStatistic::generateMemberStatisticInitData($enterpriseId,
                        $memUid, $statTime, $createdAt);
                }
            }

            if (!empty($workerStatisticData)) {
                DailyWorkerCustomerStatistic::query()->insert($workerStatisticData);
            }
            if (!empty($memberStatisticData)) {
                DailyMemberCustomerStatistic::query()->insert($memberStatisticData);
            }

            getOrSetData($flagKey, 1, 86400);
            Log::info("initEnterpriseCustomerStatistic succeed");
        } catch (\Exception $e) {
            delRedis([$flagKey]);

            Log::info("initEnterpriseCustomerStatistic exception", [
                'code' => $e->getCode(),
                'statusCode' => $e instanceof Error ? $e->getStatusCode() : 0,
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
        }
    }

    /**
     * 初始化所有企业今日累计数据
     * @param int $statTime
     * @param bool $force
     */
    public function updateAllEnterpriseCumulativeCustomerStatistic(int $statTime, bool $force)
    {
        Log::info("updateAllEnterpriseCumulativeCustomerStatistic start", ['statTime' => $statTime]);
        $enterpriseIdList = StoreWorkerService::getAllEnterpriseId();
        if (empty($enterpriseIdList)) {
            Log::info("updateAllEnterpriseCumulativeCustomerStatistic enterpriseIdList is empty");
            return;
        }

        foreach ($enterpriseIdList as $enterpriseId) {
            $this->updateEnterpriseCumulativeCustomerStatistic($enterpriseId, $statTime, $force);
        }

        Log::info("updateAllEnterpriseCumulativeCustomerStatistic end");
    }

    /**
     * 初始化今日累计数据
     * @param int $enterpriseId
     * @param int $statTime
     * @param bool $force
     */
    public function updateEnterpriseCumulativeCustomerStatistic(int $enterpriseId, int $statTime, bool $force = false)
    {
        Log::info("updateEnterpriseCumulativeCustomerStatistic params", func_get_args());
        $flagKey = redisKey('ENTERPRISE_CUMULATIVE_CUSTOMER_STATISTIC_UPDATE_FLAG', [
            'enterprise_id' => $enterpriseId,
            'stat_time' => $statTime
        ]);

        if (!$force && getOrSetData($flagKey)) {
            Log::info("updateEnterpriseCumulativeCustomerStatistic initialized");
            return;
        }

        try {

            $workerList = StoreWorker::query()
                ->select(['id', 'mem_uid'])
                ->where('enterprise_id', $enterpriseId)
                ->get()
                ->toArray();
            if (empty($workerList)) {
                Log::info("updateEnterpriseCumulativeCustomerStatistic, workerList is empty");
                return;
            }
            $memUidList = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));
            $prevDayStatTime = $statTime - DateUtils::DAY_SECONDS;
            $updatedAt = time();
            $prevDayWorkerStatisticList = DailyWorkerCustomerStatistic::listStatisticByStatTime($enterpriseId,
                $prevDayStatTime,
                [],
                [
                    'worker_id',
                    'cumulative_contact_card_way_customer_num',
                    'cumulative_other_way_customer_num',
                    'cumulative_binding_customer_num'
                ]
            );
            $prevDayWorkerStatisticMap = array_column($prevDayWorkerStatisticList, null, 'worker_id');
            $thenWorkerStatisticList = DailyWorkerCustomerStatistic::listStatisticByStatTime($enterpriseId,
                $statTime,
                [],
                [
                    'id',
                    'worker_id',
                ]
            );
            $thenWorkerStatisticMap = array_column($thenWorkerStatisticList, 'id', 'worker_id');
            foreach ($workerList as $worker) {
                $workerId = $worker['id'];
                if (isset($prevDayWorkerStatisticMap[$workerId]) && isset($thenWorkerStatisticMap[$workerId])) {
                    $prevDayWorkerStatistic = $prevDayWorkerStatisticMap[$workerId];
                    $setData = [
                        'cumulative_contact_card_way_customer_num' => $prevDayWorkerStatistic['cumulative_contact_card_way_customer_num'],
                        'cumulative_other_way_customer_num' => $prevDayWorkerStatistic['cumulative_other_way_customer_num'],
                        'cumulative_binding_customer_num' => $prevDayWorkerStatistic['cumulative_binding_customer_num'],
                        'updated_at' => $updatedAt
                    ];
                    DailyWorkerCustomerStatistic::incrStatFieldsByIds($thenWorkerStatisticMap[$workerId], $setData);
                }
            }
            $prevDayMemberStatisticList = DailyMemberCustomerStatistic::listStatisticByStatTime($enterpriseId,
                $prevDayStatTime,
                [],
                [
                    'mem_uid',
                    'cumulative_customer_num',
                    'cumulative_member_num',
                    'cumulative_non_member_num',
                    'cumulative_customer_join_group_rate',
                    'cumulative_vip_external_user_rate',
                ]
            );
            $prevDayMemberStatisticMap = array_column($prevDayMemberStatisticList, null, 'mem_uid');
            $thenMemberStatisticList = DailyMemberCustomerStatistic::listStatisticByStatTime($enterpriseId,
                $statTime,
                [],
                [
                    'id',
                    'mem_uid'
                ]
            );
            $thenMemberStatisticMap = array_column($thenMemberStatisticList, 'id', 'mem_uid');
            foreach ($memUidList as $memUid) {
                if (isset($prevDayMemberStatisticMap[$memUid]) && isset($thenMemberStatisticMap[$memUid])) {
                    $prevDayMemberStatistic = $prevDayMemberStatisticMap[$memUid];
                    $setData = [
                        'cumulative_customer_num'   => $prevDayMemberStatistic['cumulative_customer_num'],
                        'cumulative_member_num'     => $prevDayMemberStatistic['cumulative_member_num'],
                        'cumulative_non_member_num' => $prevDayMemberStatistic['cumulative_non_member_num'],
                        'cumulative_customer_join_group_rate' => $prevDayMemberStatistic['cumulative_customer_join_group_rate'],
                        'cumulative_vip_external_user_rate' => $prevDayMemberStatistic['cumulative_vip_external_user_rate'],
                        'updated_at' => $updatedAt
                    ];
                    DailyMemberCustomerStatistic::incrStatFieldsByIds($thenMemberStatisticMap[$memUid], $setData);
                }
            }

            getOrSetData($flagKey, 1, 86400);
            Log::info("updateEnterpriseCumulativeCustomerStatistic succeed");

        } catch (\Exception $e) {

            delRedis([$flagKey]);

            Log::info("updateEnterpriseCumulativeCustomerStatistic exception", [
                'code' => $e->getCode(),
                'statusCode' => $e instanceof Error ? $e->getStatusCode() : 0,
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
        }
    }

    /**
     * 根据企业ID 同步用户行为信息
     *
     * @param int $enterpriseId
     * @param int $startTime
     * @return bool
     */
    public function syncUserBehaviorByEnterpriseId(int $enterpriseId, array $startTime)
    {
        try {
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];
            $domain = $enterprise['domain'];
        } catch (\Exception $e) {
            Log::error('企业信息查询失败', [
                'enterpriseId' => $enterpriseId,
                'code' => $e->getCode(),
                'message' => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            return false;
        }
        DailyMemberCustomerStatistic::query()
            ->where(['enterprise_id' => $enterpriseId])
            ->whereIn('stat_time', $startTime)
            ->chunkById(100, function ($statistic) use ($domain) {
                $this->dailyData($statistic, $domain);
            });
        return true;
    }

    public function dailyData($statistic, $domain)
    {
        $member = new Member();
        $time = time();
        foreach ($statistic as $v) {
            try {
                $behavior = $member->getBehaviorData($domain,
                    ['memUids' => [$v->mem_uid], 'startTime' => $v->stat_time, 'endTime' => strtotime(date('Y-m-d 23:59:59', $v->stat_time)), 'dpIds' => []]);
                $behavior = (is_array($behavior) && $behavior[0]) ? $behavior[0] : [];
                $behaviorData = [
                    'new_apply_cnt' => $behavior['newApplyCnt'] ?? 0,
                    'new_contact_cnt' => $behavior['newContactCnt'] ?? 0,
                    'chat_cnt' => $behavior['chatCnt'] ?? 0,
                    'message_cnt' => $behavior['messageCnt'] ?? 0,
                    'reply_percentage' => ($behavior['replyPercentage'] ?? 0) * 100,
                    'avg_reply_time' => $behavior['avgReplyTime'] ?? 0,
                    'negative_feedback_cnt' => $behavior['negativeFeedbackCnt'] ?? 0,
                    'updated_at' => $time,
                ];
                DailyMemberCustomerStatistic::query()->where('id', $v->id)->update($behaviorData);
            } catch (\Exception $e) {
                Log::error('更新当前导购信息失败', [
                    'params' => json_encode($v),
                    'code' => $e->getCode(),
                    'message' => $e->getMessage(),
                    'exception' => $e->getTraceAsString()
                ]);
                continue;
            }
        }
        return true;
    }
}