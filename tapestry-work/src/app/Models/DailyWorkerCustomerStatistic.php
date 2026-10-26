<?php


namespace App\Models;


use App\Common\DateUtils;
use Illuminate\Database\Eloquent\Builder;

class DailyWorkerCustomerStatistic extends CustomerStatisticsBaseModel
{
    const GROUP_BY_STAT_TIME = 'stat_time';

    const GROUP_BY_WORKER = 'worker_id';

    protected $table = 'daily_worker_customer_statistic';

    protected $guarded = [

    ];

    /**
     * 统计字段列表
     * @var array
     */
    protected static $statFields = [
        'cumulative_contact_card_way_customer_num',
        'cumulative_other_way_customer_num',
        'cumulative_binding_customer_num',
        'today_contact_card_way_customer_num',
        'today_other_way_customer_num',
        'today_binding_customer_num',
    ];

    /**
     * 生成导购统计数据初始化数据
     * @param $enterpriseId
     * @param $workerId
     * @param $memUid
     * @param $statTime
     * @param int $createdAt
     * @return array
     */
    public static function generateWorkerStatisticInitData($enterpriseId, $workerId, $memUid, $statTime, $createdAt = 0)
    {
        $createdAt = $createdAt == 0 ? time() : $createdAt;
        $data = [
            'enterprise_id' => $enterpriseId,
            'stat_time' => $statTime,
            'worker_id' => $workerId,
            'mem_uid' => $memUid,
            'created_at' => $createdAt,
            'updated_at' => $createdAt
        ];
        foreach (self::$statFields as $field) {
            $data[$field] = 0;
        }
        return $data;
    }

    /**
     * 获取导购指定日期的统计数据
     * @param int $enterpriseId
     * @param int $workerId
     * @param int $statTime
     * @return Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function getWorkerStatisticByStatTime(int $enterpriseId, int $workerId, int $statTime)
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('stat_time', $statTime)
            ->where('worker_id', $workerId)
            ->first();
    }

    /**
     * 获取指定导购指定日期最近一条统计数据
     * @param int $enterpriseId
     * @param int $workerId
     * @param int $statTime
     * @return Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function getWorkerLatestStatistic(int $enterpriseId, int $workerId, int $statTime)
    {
        $statistic = self::getWorkerStatisticByStatTime($enterpriseId, $workerId, $statTime);
        if (empty($statistic)) {
            $statistic = self::query()->where('enterprise_id', $enterpriseId)
                ->where('stat_time', '<', $statTime)
                ->where('worker_id', $workerId)
                ->orderByDesc('stat_time')
                ->first();
        }
        return $statistic;
    }

    /**
     * 查询指定导购指定日期范围统计数据列表
     * @param int $enterpriseId
     * @param int $workerId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param string|array $columns
     * @return array
     */
    public static function listWorkerStatisticByStatTimeRange(int $enterpriseId,
                                                              int $workerId,
                                                              int $statTimeStart,
                                                              int $statTimeEnd,
                                                              $columns = '*')
    {
        $queryBuilder = self::query()->select($columns)->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        self::setWorkerIdListWhere($queryBuilder, [$workerId]);
        return $queryBuilder->get()->toArray();
    }

    /**
     * 查询指定日期统计数据列表
     * @param int $enterpriseId
     * @param int $statTime
     * @param array $workerIdList
     * @param array $columns
     * @return array
     */
    public static function listStatisticByStatTime(int $enterpriseId,
                                                   int $statTime,
                                                   array $workerIdList = [],
                                                   array $columns = [])
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('stat_time', $statTime)
            ->when(!empty($workerIdList), function (Builder $query) use ($workerIdList) {
                $query->whereIn('worker_id', $workerIdList);
            })
            ->when(!empty($columns), function (Builder $query) use ($columns) {
                $query->select($columns);
            })
            ->get()
            ->toArray();
    }

    /**
     * 获取指定企业指定日期范围分组聚合统计数据列表（按日期分组、按导购分组）
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param array $workerIdList
     * @param string $groupBy
     * @param array $statFields
     * @return array
     */
    public static function listAggStatistic(int $enterpriseId,
                                            int $statTimeStart,
                                            int $statTimeEnd,
                                            array $workerIdList = [],
                                            string $groupBy = '',
                                            array $statFields = [])
    {
        $columns = self::generateStatFieldsAggColumns($statFields);
        if ($groupBy != '') {
            array_unshift($columns, $groupBy);
        }
        $queryBuilder = self::query()
            ->select($columns)
            ->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        self::setWorkerIdListWhere($queryBuilder, $workerIdList);
        if ($groupBy != '') {
            $queryBuilder->groupBy($groupBy);
        }
        $statisticList = $queryBuilder->get()->toArray();
        self::formatStatColumn($statisticList);
        return $statisticList;
    }

    public static function setStatTimeRangeWhere(Builder $queryBuilder, $statTimeStart, $statTimeEnd)
    {
        if ($statTimeStart == $statTimeEnd) {
            $queryBuilder->where('stat_time', $statTimeStart);
        } else {
            $queryBuilder->whereBetween('stat_time', [$statTimeStart, $statTimeEnd]);
        }
    }

    public static function setWorkerIdListWhere(Builder $queryBuilder, array $workerIdList)
    {
        if (!empty($workerIdList)) {
            if (count($workerIdList) == 1) {
                $queryBuilder->where('worker_id', $workerIdList[0]);
            } else {
                $queryBuilder->whereIn('worker_id', $workerIdList);
            }
        }
    }

}