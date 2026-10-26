<?php


namespace App\Services\DataStatistic;


use App\Common\DateUtils;
use App\Models\DailyMemberCustomerStatistic;
use App\Services\BaseService;
use Illuminate\Support\Facades\Log;

class DailyMemberCustomerStatisticService extends BaseService
{

    /**
     * 获取员工指定日期统计数据，不存在则新建
     * @param int $enterpriseId
     * @param string $memUid
     * @param int $statTime
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function getOrCreateMemberStatistic(int $enterpriseId, string $memUid, int $statTime)
    {
        $statistic = DailyMemberCustomerStatistic::getMemberStatisticByStatTime($enterpriseId, $memUid, $statTime);
        if (empty($statistic)) {
            $statistic = $this->initMemberStatistic($enterpriseId, $memUid, $statTime);
        }
        return $statistic;
    }

    /**
     * 初始化员工指定日期统计数据
     * @param $enterpriseId
     * @param $memUid
     * @param $statTime
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model
     */
    public function initMemberStatistic($enterpriseId, $memUid, $statTime)
    {
        Log::info("incrMemberStatisticField params", func_get_args());

        $data = DailyMemberCustomerStatistic::generateMemberStatisticInitData($enterpriseId, $memUid, $statTime, time());

        // 获取员工之前的统计数据，并将累计字段赋值到今日数据上
        $prevDayStatTime = $statTime - DateUtils::DAY_SECONDS;
        $latestStatistic = DailyMemberCustomerStatistic::getMemberLatestStatistic($enterpriseId, $memUid, $prevDayStatTime);
        if (!empty($latestStatistic)) {
            $data = array_merge($data, [
                'cumulative_customer_num'   => $latestStatistic['cumulative_customer_num'],
                'cumulative_member_num'     => $latestStatistic['cumulative_member_num'],
                'cumulative_non_member_num' => $latestStatistic['cumulative_non_member_num'],
            ]);
        }

        return DailyMemberCustomerStatistic::query()->create($data);
    }

    /**
     * 递增/递减更新指定员工指定日期统计字段数据
     * @param $enterpriseId
     * @param $memUid
     * @param $statTime
     * @param array $setData
     * @param bool $createIfNotExists
     * @return bool
     */
    public function incrMemberStatisticField($enterpriseId, $memUid, $statTime, array $setData,
                                             $createIfNotExists = false)
    {
        Log::info("incrMemberStatisticField params", func_get_args());
        if ($createIfNotExists) {
            $statistic = $this->getOrCreateMemberStatistic($enterpriseId, $memUid, $statTime);
        } else {
            $statistic = DailyMemberCustomerStatistic::getMemberStatisticByStatTime($enterpriseId, $memUid, $statTime);
        }

        if (!empty($statistic)) {
            return DailyMemberCustomerStatistic::incrStatFieldsByIds($statistic['id'], $setData) == 1;
        } else {
            Log::info("member statistic is not exists", func_get_args());
            return false;
        }
    }

    /**
     * 递增/递减更新指定员工指定日期范围统计字段数据
     * @param $enterpriseId
     * @param $memUid
     * @param $startStatTime
     * @param $endStatTime
     * @param array $setData
     */
    public function incrMemberStatisticFieldByStatTimeRange($enterpriseId, $memUid, $startStatTime, $endStatTime,
                                                            array $setData)
    {
        Log::info("incrMemberStatisticFieldByStatTimeRange params", func_get_args());
        $statisticList = DailyMemberCustomerStatistic::listMemberStatisticByStatTimeRange($enterpriseId, $memUid,
            $startStatTime, $endStatTime, 'id');

        if (!empty($statisticList)) {
            // 更新添加好友日期之后所有日期的累计统计数据
            // 根据字段循环更新，防止有些数据异常更新报错
            $idList = array_column($statisticList, 'id');
            foreach ($setData as $field => $value) {
                $setData         = [$field => $value];
                $additionalWhere = $value > 0 ? [] : [[$field, '>=', abs($value)]];
                $affectRows      = DailyMemberCustomerStatistic::incrStatFieldsByIds($idList, $setData, $additionalWhere);
                Log::info("incrMemberStatisticFieldByStatTimeRange result", [
                    'field'      => $field,
                    'affectRows' => $affectRows
                ]);
            }
        }
    }

    /**
     * 获取指定日期统计数据列表
     * @param int $enterpriseId
     * @param int $statTime
     * @param array $memUidList
     * @param array $columns
     * @return array
     */
    public function listStatisticByStatTime(int $enterpriseId,
                                            int $statTime,
                                            array $memUidList = [],
                                            array $columns = [])
    {
        $listStatistic = DailyMemberCustomerStatistic::listStatisticByStatTime($enterpriseId, $statTime, $memUidList, $columns);

        if (empty($columns) || in_array('reply_percentage', $columns)) {
            foreach ($listStatistic as $key => &$value) {
                $value['reply_percentage'] = self::convertReplyPercentage($value['reply_percentage']);
            }
        }
        return $listStatistic;
    }

    /**
     * 获取指定日期范围分组聚合统计数据列表（按统计日期分组、按员工分组）
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param array $memUidList
     * @param string $groupBy
     * @param array $statFields
     * @return array
     */
    public function listGroupAggStatistics(int $enterpriseId,
                                           int $statTimeStart,
                                           int $statTimeEnd,
                                           array $memUidList = [],
                                           string $groupBy = DailyMemberCustomerStatistic::GROUP_BY_STAT_TIME,
                                           array $statFields = [])
    {
        $listAggStatistic = DailyMemberCustomerStatistic::listAggStatistic($enterpriseId, $statTimeStart, $statTimeEnd,
            $memUidList, $groupBy, $statFields);

        foreach ($listAggStatistic as $key => &$value) {
            $value['cumulative_customer_join_group_rate'] = round($value['cumulative_customer_join_group_rate'] ?? 0 / 100, 2) . '%';
            $value['cumulative_vip_external_user_rate']   = round($value['cumulative_vip_external_user_rate'] ?? 0 / 100, 2) . '%';
        }

        if (empty($statFields) || in_array('reply_percentage', $statFields)) {
            // 计算咨询回复率平均值
            if ($groupBy == DailyMemberCustomerStatistic::GROUP_BY_STAT_TIME) {
                if (!empty($memUidList)) {
                    $memUidCountMap = array_fill_keys(range($statTimeStart, $statTimeEnd, DateUtils::DAY_SECONDS), count($memUidList));
                } else {
                    // 获取指定日期范围，每天的导购数量
                    $memUidCountMap = DailyMemberCustomerStatistic::getStatTimeMemUidCountMap($enterpriseId, $statTimeStart, $statTimeEnd);
                }

                // 计算每天的咨询回复率平均值，除以导购数量
                foreach ($listAggStatistic as $key => &$value) {
                    $memUidCount               = $memUidCountMap[$value['stat_time']] ?? 0;
                    $value['reply_percentage'] = self::convertReplyPercentage($value['reply_percentage'], null, $memUidCount);
                }
            } elseif ($groupBy == DailyMemberCustomerStatistic::GROUP_BY_MEMBER) {
                // 计算每个导购指定日期范围的咨询回复率，除以间隔天数
                $intervalDays = ((int)(($statTimeEnd - $statTimeStart) / DateUtils::DAY_SECONDS)) + 1;
                foreach ($listAggStatistic as $key => &$value) {
                    $value['reply_percentage']                    = self::convertReplyPercentage($value['reply_percentage'], $intervalDays);
                    $value['cumulative_customer_join_group_rate'] = 200;
                    $value['cumulative_customer_join_group_rate'] = round($value['cumulative_customer_join_group_rate'] ?? 0 / 10000, 4);
                    $value['cumulative_vip_external_user_rate']   = round($value['cumulative_vip_external_user_rate'] ?? 0 / 10000, 4);
                }
            }
        }
        return $listAggStatistic;
    }

    /**
     * 获取指定日期范围聚合统计数据列表
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param array $memUidList
     * @param array $statFields
     * @return array
     */
    public function getAggStatistics(int $enterpriseId,
                                     int $statTimeStart,
                                     int $statTimeEnd,
                                     array $memUidList = [],
                                     array $statFields = [])
    {
        $listAggStatistic = DailyMemberCustomerStatistic::listAggStatistic($enterpriseId, $statTimeStart, $statTimeEnd,
            $memUidList, '', $statFields);
        if (!empty($listAggStatistic)) {
            $statistic = $listAggStatistic[0];

            // 计算咨询回复率平均值，除以间隔天数，除以导购数量
            if (empty($statFields) || in_array('reply_percentage', $statFields)) {
                $intervalDays = ((int)(($statTimeEnd - $statTimeStart) / DateUtils::DAY_SECONDS)) + 1;
                if (!empty($memUidList)) {
                    $memUidCount = count($memUidList);
                } else {
                    $memUidCount = DailyMemberCustomerStatistic::countMemUid($enterpriseId, $statTimeStart, $statTimeEnd);
                }
                $statistic['reply_percentage'] = self::convertReplyPercentage($statistic['reply_percentage'], $intervalDays, $memUidCount);
            }
            return $statistic;
        }
        return [];
    }

    /**
     * 计算咨询回复率平均值，保留2位小数，四舍五入
     * @param int $replyPercentage
     * @param int|null $intervalDays
     * @param int|null $memUidCount
     * @return float
     */
    public static function convertReplyPercentage(int $replyPercentage, ?int $intervalDays = null, ?int $memUidCount = null)
    {
        $avgReplyPercentage = 0;
        if ($replyPercentage > 0) {
            if ($intervalDays != null && $memUidCount != null) {
                $avgReplyPercentage = $replyPercentage / $intervalDays / $memUidCount;
            } elseif ($intervalDays != null) {
                $avgReplyPercentage = $replyPercentage / $intervalDays;
            } elseif ($memUidCount != null) {
                $avgReplyPercentage = $replyPercentage / $memUidCount;
            } else {
                $avgReplyPercentage = $replyPercentage;
            }
        }
        return round($avgReplyPercentage / 10000, 4);
    }
}