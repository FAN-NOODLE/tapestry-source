<?php


namespace App\Services\DataStatistic;


use App\Common\DateUtils;
use App\Models\DailyWorkerCustomerStatistic;
use App\Services\BaseService;
use Illuminate\Support\Facades\Log;

class DailyWorkerCustomerStatisticService extends BaseService
{

    /**
     * 获取导购指定日期统计数据，不存在则创建
     * @param int $enterpriseId
     * @param int $workerId
     * @param string $memUid
     * @param int $statTime
     * @return \Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function getOrCreateWorkerStatistic(int $enterpriseId, int $workerId, string $memUid, int $statTime)
    {
        $statistic = DailyWorkerCustomerStatistic::getWorkerStatisticByStatTime($enterpriseId, $workerId, $statTime);
        if (empty($statistic)) {
            $statistic = $this->initWorkerStatistic($enterpriseId, $workerId, $memUid, $statTime);
        }
        return $statistic;
    }

    public function initWorkerStatistic($enterpriseId, $workerId, $memUid, $statTime)
    {
        Log::info("initWorkerStatistic params", func_get_args());
        // 生成统计字段数据
        $data = DailyWorkerCustomerStatistic::generateWorkerStatisticInitData($enterpriseId, $workerId, $memUid, $statTime, time());

        // 获取导购之前的统计数据，并将累计字段赋值到今日数据上
        $prevDayStatTime = $statTime - DateUtils::DAY_SECONDS;
        $latestStatistic = DailyWorkerCustomerStatistic::getWorkerLatestStatistic($enterpriseId, $workerId, $prevDayStatTime);
        if (!empty($latestStatistic)) {
            Log::info("initWorkerStatistic latestStatistic", $latestStatistic->toArray());
            $data = array_merge($data, [
                'cumulative_contact_card_way_customer_num' => $latestStatistic['cumulative_contact_card_way_customer_num'],
                'cumulative_other_way_customer_num' => $latestStatistic['cumulative_other_way_customer_num'],
                'cumulative_binding_customer_num' => $latestStatistic['cumulative_binding_customer_num']
            ]);
        }

        return DailyWorkerCustomerStatistic::query()->create($data);
    }

    /**
     * 递增/递减更新指定导购指定日期统计字段数据
     * @param $enterpriseId
     * @param $workerId
     * @param $memUid
     * @param $statTime
     * @param array $setData
     * @param bool $createIfNotExists
     * @return bool
     */
    public function incrWorkerStatisticField($enterpriseId, $workerId, $memUid, $statTime, array $setData,
                                             $createIfNotExists = false)
    {
        Log::info("incrWorkerStatisticField params", func_get_args());
        if ($createIfNotExists) {
            $statistic = $this->getOrCreateWorkerStatistic($enterpriseId, $workerId, $memUid, $statTime);
        } else {
            $statistic = DailyWorkerCustomerStatistic::getWorkerStatisticByStatTime($enterpriseId, $workerId, $statTime);
        }

        if (!empty($statistic)) {
            return DailyWorkerCustomerStatistic::incrStatFieldsByIds($statistic['id'], $setData) == 1;
        } else {
            Log::info("incrWorkerStatisticField statistic is not exists", func_get_args());
            return false;
        }
    }

    /**
     * 递增/递减更新指定导购指定日期范围统计字段数据
     * @param $enterpriseId
     * @param $workerId
     * @param $startStatTime
     * @param $endStatTime
     * @param array $setData
     */
    public function incrWorkerStatisticFieldByStatTimeRange($enterpriseId, $workerId, $startStatTime, $endStatTime,
                                                            array $setData)
    {
        Log::info("incrWorkerStatisticFieldByStatTimeRange params", func_get_args());
        $statisticList = DailyWorkerCustomerStatistic::listWorkerStatisticByStatTimeRange($enterpriseId, $workerId,
            $startStatTime, $endStatTime, 'id');

        if (!empty($statisticList)) {
            // 更新添加好友日期之后所有日期的累计统计数据
            // 根据字段循环更新，防止有些数据异常更新报错
            $idList = array_column($statisticList, 'id');
            foreach ($setData as $field => $value) {
                $setData = [$field => $value];
                $additionalWhere = $value > 0 ? [] : [[$field, '>=', abs($value)]];
                $affectRows = DailyWorkerCustomerStatistic::incrStatFieldsByIds($idList, $setData, $additionalWhere);
                Log::info("incrWorkerStatisticFieldByStatTimeRange result", [
                    'field' => $field,
                    'affectRows' => $affectRows
                ]);
            }
        }
    }

    /**
     * 更新导购指定日期统计字段数据
     * @param $enterpriseId
     * @param $memUid
     * @param $workerId
     * @param $statTime
     * @param array $setData
     * @param bool $createIfNotExists
     * @return bool
     */
    public function updateWorkerStatisticField($enterpriseId, $memUid, $workerId, $statTime, array $setData,
                                               $createIfNotExists = false)
    {
        Log::info("updateWorkerStatisticField params", func_get_args());
        if ($createIfNotExists) {
            $statistic = $this->getOrCreateWorkerStatistic($enterpriseId, $workerId, $memUid, $statTime);
        } else {
            $statistic = DailyWorkerCustomerStatistic::getWorkerStatisticByStatTime($enterpriseId, $workerId, $statTime);
        }
        if (!empty($statistic)) {
            if (!isset($setData['updated_at'])) {
                $setData['updated_at'] = time();
            }
            return DailyWorkerCustomerStatistic::query()->where('id', $statistic['id'])->update($setData) == 1;
        } else {
            Log::info("updateWorkerStatisticField statistic is not exists");
            return false;
        }
    }

    /**
     * 获取指定企业指定日期统计数据列表
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
        return DailyWorkerCustomerStatistic::listStatisticByStatTime($enterpriseId, $statTime, $workerIdList, $columns);
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
    public function listGroupAggStatistics(int $enterpriseId,
                                           int $statTimeStart,
                                           int $statTimeEnd,
                                           array $workerIdList = [],
                                           string $groupBy = DailyWorkerCustomerStatistic::GROUP_BY_STAT_TIME,
                                           array $statFields = [])
    {
        return DailyWorkerCustomerStatistic::listAggStatistic($enterpriseId, $statTimeStart, $statTimeEnd,
            $workerIdList, $groupBy, $statFields);
    }


    /**
     * 获取指定企业指定日期范围聚合统计数据列表
     * @param int $enterpriseId
     * @param int $statTimeStart
     * @param int $statTimeEnd
     * @param array $workerIdList
     * @param array $statFields
     * @return array
     */
    public function getAggStatistics(int $enterpriseId,
                                     int $statTimeStart,
                                     int $statTimeEnd,
                                     array $workerIdList = [],
                                     array $statFields = [])
    {
        $listAggStatistic = DailyWorkerCustomerStatistic::listAggStatistic($enterpriseId, $statTimeStart, $statTimeEnd,
            $workerIdList, '', $statFields);
        return !empty($listAggStatistic) ? $listAggStatistic[0] : [];
    }


}