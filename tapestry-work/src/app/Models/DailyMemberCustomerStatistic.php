<?php


namespace App\Models;



use App\Common\DateUtils;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Support\Facades\DB;

class DailyMemberCustomerStatistic extends CustomerStatisticsBaseModel
{
    const GROUP_BY_STAT_TIME = 'stat_time';

    const GROUP_BY_MEMBER = 'mem_uid';

    protected $table = 'daily_member_customer_statistic';

    protected $guarded = [

    ];

    /**
     * 统计字段列表
     * @var array
     */
    public static $statFields = [
        'cumulative_customer_num',
        'today_customer_num',
        'new_apply_cnt',
        'new_contact_cnt',
        'chat_cnt',
        'message_cnt',
        'reply_percentage',
        'avg_reply_time',
        'negative_feedback_cnt',
        'cumulative_member_num',
        'today_member_num',
        'cumulative_non_member_num',
        'today_non_member_num',
        'cumulative_customer_join_group_rate',
        'cumulative_vip_external_user_rate',
    ];

    /**
     * 生成员工统计数据初始化数据
     * @param $enterpriseId
     * @param $memUid
     * @param $statTime
     * @param int $createdAt
     * @return array
     */
    public static function generateMemberStatisticInitData($enterpriseId, $memUid, $statTime, $createdAt = 0)
    {
        $createdAt = $createdAt == 0 ? time() : $createdAt;
        $data = [
            'enterprise_id' => $enterpriseId,
            'stat_time' => $statTime,
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
     * @param string $memUid
     * @param int $statTime
     * @return Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function getMemberStatisticByStatTime(int $enterpriseId, string $memUid, int $statTime)
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('stat_time', $statTime)
            ->where('mem_uid', $memUid)
            ->first();
    }

    /**
     * 获取指定导购指定日期之前最后的统计数据
     * @param int $enterpriseId
     * @param string $memUid
     * @param int $statTime
     * @return Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public static function getMemberLatestStatistic(int $enterpriseId, string $memUid, int $statTime)
    {
        $statistic = self::getMemberStatisticByStatTime($enterpriseId, $memUid, $statTime);
        if (empty($statistic)) {
            $statistic = self::query()->where('enterprise_id', $enterpriseId)
                ->where('stat_time', '<', $statTime)
                ->where('mem_uid', $memUid)
                ->orderByDesc('stat_time')
                ->first();
        }
        return $statistic;
    }

    /**
     * 查询指定导购指定日期范围统计数据列表
     * @param int $enterpriseId
     * @param string $memUid
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param string|array $columns
     * @return array
     */
    public static function listMemberStatisticByStatTimeRange(int $enterpriseId,
                                                              string $memUid,
                                                              int $statTimeStart,
                                                              int $statTimeEnd,
                                                              $columns = '*')
    {
        $queryBuilder = self::query()->select($columns)->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        self::setMemUidListWhere($queryBuilder, [$memUid]);
        return $queryBuilder->get()->toArray();
    }

    /**
     * 查询指定日期统计数据列表
     * @param int $enterpriseId
     * @param int $statTime
     * @param array $memUidList
     * @param array $columns
     * @return array
     */
    public static function listStatisticByStatTime(int $enterpriseId,
                                                   int $statTime,
                                                   array $memUidList = [],
                                                   array $columns = [])
    {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('stat_time', $statTime)
            ->when(!empty($memUidList), function (Builder $query) use ($memUidList) {
                $query->whereIn('mem_uid', $memUidList);
            })
            ->when(!empty($columns),function (Builder $query) use ($columns) {
                $query->select($columns);
            })
            ->get()
            ->toArray();
    }

    /**
     * 获取指定企业指定日期范围分组聚合统计数据列表（按日期分组、按员工分组）
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param array $memUidList
     * @param string $groupBy
     * @return array
     */
    public static function listAggStatistic(int $enterpriseId,
                                            int $statTimeStart,
                                            int $statTimeEnd,
                                            array $memUidList = [],
                                            string $groupBy = '',
                                            array $statFields = [])
    {
        $queryBuilder = self::query()->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        self::setMemUidListWhere($queryBuilder, $memUidList);

        $columns = self::generateStatFieldsAggColumns($statFields);
        if ($groupBy != '') {
            array_unshift($columns, $groupBy);
            $queryBuilder->groupBy($groupBy);
        }
        $queryBuilder->select($columns);
        $statisticList = $queryBuilder->get()->toArray();

        self::formatStatColumn($statisticList);
        return $statisticList;
    }

    /**
     * 查询指定日期范围导购数量
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @return int
     */
    public static function countMemUid(int $enterpriseId,
                                       int $statTimeStart,
                                       int $statTimeEnd)
    {
        $queryBuilder = self::query()->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        return $queryBuilder->count(DB::raw('DISTINCT(mem_uid)'));
    }

    /**
     * 查询指定日期范围每日的导购数量
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @return array
     */
    public static function getStatTimeMemUidCountMap(int $enterpriseId,
                                                     int $statTimeStart,
                                                     int $statTimeEnd)
    {
        $queryBuilder = self::query()->where('enterprise_id', $enterpriseId);
        self::setStatTimeRangeWhere($queryBuilder, $statTimeStart, $statTimeEnd);
        $queryBuilder->select([
            'stat_time',
            DB::raw('count(DISTINCT mem_uid) AS cnt')
        ]);
        $result = $queryBuilder->groupBy('stat_time')->get()->toArray();

        return !empty($result) ? array_column($result, 'cnt', 'stat_time') : [];
    }

    public static function setStatTimeRangeWhere(Builder $queryBuilder, $statTimeStart, $statTimeEnd)
    {
        if ($statTimeStart == $statTimeEnd) {
            $queryBuilder->where('stat_time', $statTimeStart);
        } else {
            $queryBuilder->whereBetween('stat_time', [$statTimeStart, $statTimeEnd]);
        }
    }

    public static function setMemUidListWhere(Builder $queryBuilder, array $memUidList)
    {
        if (!empty($memUidList)) {
            if (count($memUidList) == 1) {
                $queryBuilder->where('mem_uid', $memUidList[0]);
            } else {
                $queryBuilder->whereIn('mem_uid', $memUidList);
            }
        }
    }

}