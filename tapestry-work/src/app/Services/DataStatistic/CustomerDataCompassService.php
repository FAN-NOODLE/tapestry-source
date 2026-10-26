<?php


namespace App\Services\DataStatistic;


use App\Common\Constant;
use App\Common\DateUtils;
use App\Library\Uc\Department;
use App\Models\DailyMemberCustomerStatistic;
use App\Models\DailyWorkerCustomerStatistic;
use App\Models\MemberUserRelation;
use App\Models\Scrm\WorkExternalUserRelation;
use App\Models\ScrmEnterprise\WxUserCorpExternal;
use App\Services\BaseService;
use App\Services\GuideManagerService;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\DB;
use App\Library\DataLake\UserRequestService;

/**
 * 数据罗盘-顾客统计数据
 * Class CompassCustomerService
 * @package App\Services\DataStatistic
 */
class CustomerDataCompassService extends BaseService
{

    const DEFAULT_WORKER_STATISTIC = [
        'cumulative_contact_card_way_customer_num' => 0,
        'cumulative_other_way_customer_num'        => 0,
        'cumulative_binding_customer_num'          => 0,
        'today_contact_card_way_customer_num'      => 0,
        'today_other_way_customer_num'             => 0,
        'today_binding_customer_num'               => 0
    ];

    const DEFAULT_MEMBER_STATISTIC = [
        'cumulative_customer_num'             => 0,
        'today_customer_num'                  => 0,
        'cumulative_member_num'               => 0,
        'today_member_num'                    => 0,
        'cumulative_non_member_num'           => 0,
        'today_non_member_num'                => 0,
        'cumulative_customer_join_group_rate' => 0,
        'cumulative_vip_external_user_rate'   => 0,
        'chat_cnt'                            => 0,
        'reply_percentage'                    => 0.00,
    ];

    const DEFAULT_STATISTIC = self::DEFAULT_WORKER_STATISTIC + self::DEFAULT_MEMBER_STATISTIC;


    const CUMULATIVE_COLUMNS = [
        'cumulative_contact_card_way_customer_num',
        'cumulative_other_way_customer_num',
        'cumulative_binding_customer_num',
        'cumulative_customer_num'
    ];

    private $dailyWorkerCustomerStatisticService;

    private $dailyMemberCustomerStatisticService;

    /**
     * CustomerStatisticService constructor.
     */
    public function __construct()
    {
        $this->dailyWorkerCustomerStatisticService = new DailyWorkerCustomerStatisticService();
        $this->dailyMemberCustomerStatisticService = new DailyMemberCustomerStatisticService();
    }

    /**
     * 概览数据
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $workerList
     * @return array
     */
    public function overview(int $enterpriseId, int $startStatTime, int $endStatTime, array $workerList): array
    {
        $workerIdList = array_column($workerList, 'id');
        $memUidList   = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));

        // 累计类字段，查询今日聚合数据
        $todayStatTime             = DateUtils::todayStatTime();
        $workerCumulativeStatistic = $this->dailyWorkerCustomerStatisticService->getAggStatistics(
            $enterpriseId,
            $todayStatTime,
            $todayStatTime,
            $workerIdList,
            [
                'cumulative_contact_card_way_customer_num',
                'cumulative_other_way_customer_num',
                'cumulative_binding_customer_num',
            ]
        );
        $memberCumulativeStatistic = $this->dailyMemberCustomerStatisticService->getAggStatistics(
            $enterpriseId,
            $todayStatTime,
            $todayStatTime,
            $memUidList,
            [
                'cumulative_customer_num',
                'cumulative_member_num',
                'cumulative_non_member_num',
                'cumulative_customer_join_group_rate',
                'cumulative_vip_external_user_rate',
            ]
        );

        // 非累计类字段，查询时间范围内聚合数据
        $workerTimeRangeStatistic = $this->dailyWorkerCustomerStatisticService->getAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $workerIdList,
            [
                'today_contact_card_way_customer_num',
                'today_other_way_customer_num',
                'today_binding_customer_num',
            ]
        );
        $memberTimeRangeStatistic = $this->dailyMemberCustomerStatisticService->getAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memUidList,
            [
                'today_customer_num',
                'chat_cnt',
                'reply_percentage',
                'today_member_num',
                'today_non_member_num',
            ]
        );

        return array_merge(self::DEFAULT_STATISTIC, $workerCumulativeStatistic, $memberCumulativeStatistic,
            $workerTimeRangeStatistic, $memberTimeRangeStatistic);
    }

    /**
     * 按日期查询顾客数据明细
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $workerList
     * @param string $sortColumn
     * @param int $sortType
     * @param int $page
     * @param int $size
     * @param bool $forPage
     * @param bool $defaultSortDesc stat_time默认排序规则
     * @return array
     */
    public function listByDate(int $enterpriseId,
                               int $startStatTime,
                               int $endStatTime,
                               array $workerList,
                               string $sortColumn,
                               int $sortType,
                               int $page,
                               int $size,
                               bool $forPage = true,
                               bool $defaultSortDesc = true): array
    {
        $workerIdList = array_column($workerList, 'id');
        $memUidList   = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));

        $intervalDays = ((int)(($endStatTime - $startStatTime) / DateUtils::DAY_SECONDS)) + 1;
        $offset       = ($page - 1) * $size;
        if ($offset >= $intervalDays) {
            return [
                'total' => $intervalDays,
                'list'  => []
            ];
        }

        // 导购维度统计数据
        $workerStatisticList = $this->dailyWorkerCustomerStatisticService->listGroupAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $workerIdList,
            DailyWorkerCustomerStatistic::GROUP_BY_STAT_TIME
        );
        $workerStatisticMap  = array_column($workerStatisticList, null, DailyWorkerCustomerStatistic::GROUP_BY_STAT_TIME);

        // 员工维度统计数据
        $memberStatisticList = $this->dailyMemberCustomerStatisticService->listGroupAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memUidList,
            DailyMemberCustomerStatistic::GROUP_BY_STAT_TIME,
            [
                'cumulative_customer_num',
                'today_customer_num',
                'chat_cnt',
                'reply_percentage',
                'cumulative_member_num',
                'today_member_num',
                'cumulative_non_member_num',
                'today_non_member_num',
                'cumulative_customer_join_group_rate',
                'cumulative_vip_external_user_rate',
            ]
        );
        $memberStatisticMap  = array_column($memberStatisticList, null, DailyMemberCustomerStatistic::GROUP_BY_STAT_TIME);

        // 拼装筛选日期中每日数据
        $aggStatisticList = [];
        for ($n = 0; $n < $intervalDays; $n++) {
            $statTime  = $startStatTime + ($n * DateUtils::DAY_SECONDS);
            $statistic = array_merge(
                ['stat_time' => $statTime],
                $workerStatisticMap[$statTime] ?? self::DEFAULT_WORKER_STATISTIC,
                $memberStatisticMap[$statTime] ?? self::DEFAULT_MEMBER_STATISTIC
            );
            array_push($aggStatisticList, $statistic);
        }

        // 排序，一级排序值相同时按日期升序排序
        if (count($aggStatisticList) > 1) {
            $aggStatisticCollection = collect($aggStatisticList)->sortBy('stat_time', SORT_REGULAR, $defaultSortDesc);
            if (!empty($sortColumn)) {
                $aggStatisticList = $aggStatisticCollection->sortBy($sortColumn, SORT_REGULAR, SORT_DESC == $sortType)
                    ->values()->toArray();
            } else {
                $aggStatisticList = $aggStatisticCollection->values()->toArray();
            }
        }

        return [
            'total' => $intervalDays,
            'list'  => $forPage ? array_slice($aggStatisticList, $offset, $size) : $aggStatisticList
        ];
    }

    /**
     * 按导购维度查询顾客数据
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $workerList
     * @param string $sortColumn
     * @param int $sortType
     * @param int $page
     * @param int $size
     * @param bool $forPage
     * @return array
     */
    public function listByWorker(int $enterpriseId,
                                 int $startStatTime,
                                 int $endStatTime,
                                 array $workerList,
                                 string $sortColumn,
                                 int $sortType,
                                 int $page,
                                 int $size,
                                 bool $forPage = true): array
    {
        $total  = count($workerList);
        $offset = ($page - 1) * $size;
        if ($offset >= $total) {
            return [
                'total' => $total,
                'list'  => []
            ];
        }

        $workerIdList = array_column($workerList, 'id');
        $memUidList   = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));

        $todayStatTime             = DateUtils::todayStatTime();
        $workerCumulativeStatistic = $this->dailyWorkerCustomerStatisticService->listStatisticByStatTime(
            $enterpriseId,
            $todayStatTime,
            $workerIdList,
            [
                'worker_id',
                'cumulative_contact_card_way_customer_num',
                'cumulative_other_way_customer_num',
                'cumulative_binding_customer_num',
            ]
        );
        $workerCumulativeStatistic = array_column($workerCumulativeStatistic, null, DailyWorkerCustomerStatistic::GROUP_BY_WORKER);
        // 导购维度统计数据
        $workerTimeRangeStatistic = $this->dailyWorkerCustomerStatisticService->listGroupAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $workerIdList,
            DailyWorkerCustomerStatistic::GROUP_BY_WORKER,
            [
                'today_contact_card_way_customer_num',
                'today_other_way_customer_num',
                'today_binding_customer_num',
            ]
        );
        $workerTimeRangeStatistic = array_column($workerTimeRangeStatistic, null, DailyWorkerCustomerStatistic::GROUP_BY_WORKER);

        $memberCumulativeStatistic = $this->dailyMemberCustomerStatisticService->listStatisticByStatTime(
            $enterpriseId,
            $todayStatTime,
            $memUidList,
            [
                'mem_uid',
                'cumulative_customer_num',
                'cumulative_member_num',
                'cumulative_non_member_num',
            ]
        );
        $memberCumulativeStatistic = array_column($memberCumulativeStatistic, null, DailyMemberCustomerStatistic::GROUP_BY_MEMBER);

        // 员工维度统计数据
        $memberTimeRangeStatistic = $this->dailyMemberCustomerStatisticService->listGroupAggStatistics(
            $enterpriseId,
            $startStatTime,
            $endStatTime,
            $memUidList,
            DailyMemberCustomerStatistic::GROUP_BY_MEMBER,
            [
                'today_customer_num',
                'today_member_num',
                'today_non_member_num',
                'chat_cnt',
                'reply_percentage',
            ]
        );
        $memberTimeRangeStatistic = array_column($memberTimeRangeStatistic, null, DailyMemberCustomerStatistic::GROUP_BY_MEMBER);

        // 拼装数据，没有数据给默认值0
        $aggStatisticList = [];
        foreach ($workerList as $worker) {
            array_push($aggStatisticList, array_merge(
                [
                    'worker_id' => $worker['id'],
                    'mem_uid'   => $worker['mem_uid']
                ],
                self::DEFAULT_STATISTIC,
                $workerCumulativeStatistic[$worker['id']] ?? [],
                $workerTimeRangeStatistic[$worker['id']] ?? [],
                $memberCumulativeStatistic[$worker['mem_uid']] ?? [],
                $memberTimeRangeStatistic[$worker['mem_uid']] ?? []
            ));
        }

        // 排序，一级排序值相同时按worker_id降序排序
        if (count($aggStatisticList) > 1) {
            $aggStatisticCollection = collect($aggStatisticList)->sortByDesc('worker_id');
            if (!empty($sortColumn)) {
                $aggStatisticList = $aggStatisticCollection->sortBy($sortColumn, SORT_REGULAR, SORT_DESC == $sortType)
                    ->values()->toArray();
            } else {
                $aggStatisticList = $aggStatisticCollection->values()->toArray();
            }
        }

        $result = $forPage ? array_slice($aggStatisticList, $offset, $size) : $aggStatisticList;

        return [
            'total' => $total,
            'list'  => $result
        ];

    }

    /**
     * 顾客数据排行
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param string $rankColumn
     * @param int $rankLength
     * @param array $workerList
     * @return array
     */
    public function rank(int $enterpriseId, int $startStatTime, int $endStatTime, string $rankColumn, int $rankLength, array $workerList)
    {
        $workerIdList = array_column($workerList, 'id');
        $memUidList   = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));

        $isCumulativeColumn = in_array($rankColumn, self::CUMULATIVE_COLUMNS);
        $todayStatTime      = DateUtils::todayStatTime();
        $statisticList      = [];

        if (array_key_exists($rankColumn, self::DEFAULT_WORKER_STATISTIC)) {
            if ($isCumulativeColumn) {
                $workerStatisticList = $this->dailyWorkerCustomerStatisticService->listStatisticByStatTime(
                    $enterpriseId,
                    $todayStatTime,
                    $workerIdList,
                    ['worker_id', $rankColumn]
                );
            } else {
                $workerStatisticList = $this->dailyWorkerCustomerStatisticService->listGroupAggStatistics(
                    $enterpriseId,
                    $startStatTime,
                    $endStatTime,
                    $workerIdList,
                    DailyWorkerCustomerStatistic::GROUP_BY_WORKER,
                    [$rankColumn]
                );
            }
            $workerStatisticMap = array_column($workerStatisticList, $rankColumn, 'worker_id');
            foreach ($workerIdList as $workerId) {
                $statisticList[] = [
                    'worker_id' => $workerId,
                    'num'       => $workerStatisticMap[$workerId] ?? 0
                ];
            }
        } else {
            $memberStatisticList = [];
            if (!empty($memUidList)) {
                if ($isCumulativeColumn) {
                    $memberStatisticList = $this->dailyMemberCustomerStatisticService->listStatisticByStatTime(
                        $enterpriseId,
                        $todayStatTime,
                        $memUidList,
                        ['mem_uid', $rankColumn]
                    );
                } else {
                    $memberStatisticList = $this->dailyMemberCustomerStatisticService->listGroupAggStatistics(
                        $enterpriseId,
                        $startStatTime,
                        $endStatTime,
                        $memUidList,
                        DailyMemberCustomerStatistic::GROUP_BY_MEMBER,
                        [$rankColumn]
                    );
                }
            }
            $memberStatisticMap = array_column($memberStatisticList, $rankColumn, 'mem_uid');
            foreach ($workerList as $worker) {
                $statisticList[] = [
                    'worker_id' => $worker['id'],
                    'num'       => $memberStatisticMap[$worker['mem_uid']] ?? 0
                ];
            }
        }

        // 倒序排序，并获取指定名次的数量
        if (count($statisticList) > 1) {
            $statisticList = collect($statisticList)->sortByDesc('worker_id')->sortByDesc('num')->values()->toArray();
            return array_slice($statisticList, 0, $rankLength);
        } else {
            return $statisticList;
        }
    }

    /**
     * 顾客数据排行
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param int $rankLength
     * @param array $workerList
     * @return array
     */
    public function rankForConsole(int $enterpriseId, int $startStatTime, int $endStatTime, int $rankLength, array $workerList)
    {
        $list               = MemberUserRelation::query()->where('enterprise_id', $enterpriseId)
            ->where('created_at', '>', $startStatTime)
            ->where('created_at', '<', $endStatTime)
            ->select(DB::raw('count(*) as user_count, mem_uid'))
            ->groupBy('mem_uid')
            ->get()->toArray();
        $memberStatisticMap = array_column($list, 'user_count', 'mem_uid');
        $statisticList      = [];
        foreach ($workerList as $worker) {
            $statisticList[] = [
                'worker_id' => $worker['id'],
                'num'       => $memberStatisticMap[$worker['mem_uid']] ?? 0
            ];
        }
        // 倒序排序，并获取指定名次的数量
        if (count($statisticList) > 1) {
            $statisticList = collect($statisticList)->sortByDesc('worker_id')->sortByDesc('num')->values()->toArray();

            if ($rankLength) {
                $data = array_slice($statisticList, 0, $rankLength);
            } else {
                $data = array_slice($statisticList, 0);
            }
            return $data;
        } else {
            return $statisticList;
        }
    }

    /**
     * 添加统计数据导购相关信息
     * @param array $statisticList
     * @param array $workerList
     * @param string $qyDomain
     * @throws \Exception
     */
    public function decorateWorkerInfo(array &$statisticList, array $workerList, string $qyDomain)
    {
        $workerMap = array_column($workerList, null, 'id');

        $departmentMap     = [];
        $storeMap          = [];
        $departmentPathMap = [];
        // 查询导购门店和所属部门信息
        if (!empty($statisticList)) {
            $workerIdList = array_column($statisticList, 'worker_id');
            $memUidList   = $storeIdList = [];
            foreach ($workerIdList as $workerId) {
                $worker = $workerMap[$workerId];
                if (!empty($worker)) {
                    $storeIdList[] = $worker['store_id'];
                    $memUidList[]  = $worker['mem_uid'];
                }
            }
            $memUidList  = array_values(array_unique(array_filter($memUidList)));
            $storeIdList = array_values(array_unique(array_filter($storeIdList)));

            if (!empty($memUidList)) {
                $departmentList     = GuideManagerService::getGuideDepartment($memUidList, $qyDomain);
                $departmentMap      = array_column($departmentList, null, 'memUid');
                $departmentIdArr    = array_values(array_unique(array_column($departmentList, 'dpId')));
                $departmentListInfo = (new Department())->list($qyDomain, '', ['dpIds' => $departmentIdArr, 'pageSize' => 100])['list'];
                $departmentPathMap  = array_column($departmentListInfo, 'departmentPath', 'dpId');
            }
            if (!empty($storeIdList)) {
                $storeList = (new Store())->listByIds(['store_ids' => $storeIdList]);
                $storeMap  = array_column($storeList, null, 'id');
            }
        }

        // 拼装数据
        foreach ($statisticList as &$statistic) {
            $worker                 = $workerMap[$statistic['worker_id']];
            $statistic['worker_no'] = $worker['worker_no'];
            $statistic['name']      = $worker['name'];
            $statistic['avatar']    = $worker['avatar'];
            $statistic['store_id']  = $worker['store_id'];
            $statistic['mem_uid']   = $worker['mem_uid'];


            if ($worker['type'] == Constant::STORE_WORKER_TYPE_STORE) {
                $store                   = $storeMap[$worker['store_id']] ?? [];
                $statistic['store_name'] = $store['short_name'] ?? $store['name'] ?? '';
            } else {
                $statistic['store_name'] = '';
            }

            $statistic['dp_name'] = isset($departmentMap[$worker['mem_uid']])
            && isset($departmentMap[$worker['mem_uid']]['dpName']) ? $departmentMap[$worker['mem_uid']]['dpName'] : '';

            $statistic['dp_name_path'] = isset($departmentMap[$worker['mem_uid']]) ? $departmentPathMap[$departmentMap[$worker['mem_uid']]['dpId']] : '';
        }
    }

    /**
     * 获取会员等级分布
     * @param int $enterpriseId
     * @param int $startStatTime
     * @param int $endStatTime
     * @param array $memUid
     * @param int $type 类型（1-截止当前选择日期之前的数据，2-按照时间范围筛选）
     * @return array
     * @Date: 2021/3/29
     * @Time: 19:46
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function getMemberLevelStatistics(int $enterpriseId, int $startStatTime, int $endStatTime, array $memUid, int $type = 1): array
    {
        $grade = (new UserRequestService($enterpriseId))->getTierList();

        if (!$enterpriseId) {
            $temp = [];
            foreach ($grade as $key => $value) {
                if ($value['category'] != 'default') {
                    unset($grade[$key]);
                }

                if (empty(json_decode($value['remark'], true))) {
                    unset($grade[$key]);
                }
                array_push($temp, [
                    'name'  => $value['tierName'],
                    'value' => 0
                ]);
            }

            return [
                'memberGrade' => $temp,
                'memberTotal' => 0
            ];
        }

        $contactList = MemberUserRelation::query()
            ->select(['grade', DB::raw('count(*) as count')])
            ->where('enterprise_id', $enterpriseId)
            ->when(!empty($memUid), function ($query) use ($memUid) {
                return $query->whereIn('mem_uid', $memUid);
            })
            ->where('grade', '!=', '')
            ->when(!empty($startStatTime) && !empty($endStatTime), function ($query) use ($startStatTime, $endStatTime, $type) {
                if ($type == 1) {
                    return $query->where('created_at', '<=', $endStatTime);
                } else {
                    return $query->whereBetween('created_at', [$startStatTime, $endStatTime]);
                }
            })
            ->where('delete_type', 0)
            ->groupBy('grade')
            ->get()
            ->toArray();

        $contactGradeMap = array_column($contactList, 'count', 'grade');

        if ($grade) {
            foreach ($grade as $key => &$value) {
                if ($value['category'] != 'default') {
                    unset($grade[$key]);
                }

                if (empty(json_decode($value['remark'], true))) {
                    unset($grade[$key]);
                }
            }
        }

        $grade = array_values($grade);

        $gradeNameArray = array_column($grade, 'tierName');


        $memberTotal = array_sum(array_values($contactGradeMap));
        $memberGrade = [];
        foreach ($gradeNameArray as $gradeName) {
            $value = 0;
            if (isset($contactGradeMap[$gradeName])) {
                $value = $contactGradeMap[$gradeName];
            }

            array_push($memberGrade, [
                'name'  => $gradeName,
                'value' => $value
            ]);
        }

        return [
            'memberGrade' => $memberGrade,
            'memberTotal' => $memberTotal ?? 0
        ];
    }


    public function getMemberLevelStatisticsNew($enterpriseId, $memId): array
    {
        $grade = (new UserRequestService($enterpriseId))->getTierList();

        if (!$enterpriseId) {
            $temp = [];
            foreach ($grade as $key => $value) {
                if ($value['category'] != 'default') {
                    unset($grade[$key]);
                }

                if (empty(json_decode($value['remark'], true))) {
                    unset($grade[$key]);
                }
                array_push($temp, [
                    'name'  => $value['tierName'],
                    'value' => 0
                ]);
            }

            return [
                'memberGrade' => $temp,
                'memberTotal' => 0
            ];
        }
        $userList    = WorkExternalUserRelation::query()
            ->select('work_external_user_relation.work_external_user_id')
            ->where('work_external_user_relation.enterprise_id', $enterpriseId)
            ->where('work_external_user_relation.member_id', $memId)
            ->where('work_external_user_relation.relation_status', 0)
            ->pluck('work_external_user_id')
            ->toArray();
        $contactList = WxUserCorpExternal::query()
            ->setConn($enterpriseId)
            ->select(['wx_user_corp_external.grade', DB::raw('count(*) as count')])
            ->whereIn('wx_user_corp_external.id', $userList)
            ->where('wx_user_corp_external.grade', '<>', '')
            ->get()
            ->toArray();

        $contactGradeMap = array_column($contactList, 'count', 'grade');
        if ($grade) {
            foreach ($grade as $key => &$value) {
                if ($value['category'] != 'default') {
                    unset($grade[$key]);
                }

                if (empty(json_decode($value['remark'], true))) {
                    unset($grade[$key]);
                }
            }
        }

        $grade = array_values($grade);

        $gradeNameArray = array_column($grade, 'tierName');


        $memberTotal = array_sum(array_values($contactGradeMap));
        $memberGrade = [];
        foreach ($gradeNameArray as $gradeName) {
            $value = 0;
            if (isset($contactGradeMap[$gradeName])) {
                $value = $contactGradeMap[$gradeName];
            }

            array_push($memberGrade, [
                'name'  => $gradeName,
                'value' => $value
            ]);
        }

        return [
            'memberGrade' => $memberGrade,
            'memberTotal' => $memberTotal ?? 0
        ];
    }


}
