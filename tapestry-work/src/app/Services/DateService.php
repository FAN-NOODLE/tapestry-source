<?php

namespace App\Services;

class DateService
{

    /**
     * 获取本周的时间起止
     *
     * @param int $timeType 1=时间戳（秒），2=时间格式（年-月-日 时:分:秒）
     * @return array
     */
    public function getWeekTime($timeType = 1)
    {
        // 当前日期
        $sdefaultDate = date("Y-m-d");
        // $first =1 表示每周星期一为开始日期 0表示每周日为开始日期
        $first = 1;
        // 获取当前周的第几天 周日是 0 周一到周六是 1 - 6
        $w = date('w', strtotime($sdefaultDate));
        // 获取本周开始日期，如果$w是0，则表示周日，减去 6 天
        $week_start = date('Y-m-d 00:00:00', strtotime("$sdefaultDate -" . ($w ? $w - $first : 6) . ' days'));
        // 本周结束日期
        $week_end_sorttime = strtotime("$week_start +6 days") > time() ? time() : strtotime("$week_start +6 days");
        $week_end = date('Y-m-d 23:59:59', $week_end_sorttime);

        return [
            'start_time' => $timeType == 1 ? strtotime($week_start) : $week_start,
            'end_time'   => $timeType == 1 ? strtotime($week_end) : $week_end,
        ];
    }

    /**
     * 获取本月时间起止
     *
     * @param int $timeType 1=时间戳（秒），2=时间格式（年-月-日 时:分:秒）
     * @return array
     */
    public function getMonthTime($timeType = 1)
    {
        // 本月1号0点
        $startDate = date("Y-m-1 00:00:00");
        // 当前日期
        $endDate = date("Y-m-d 23:59:59");

        return [
            'start_time' => $timeType == 1 ? strtotime($startDate) : $startDate,
            'end_time'   => $timeType == 1 ? strtotime($endDate) : $endDate,
        ];
    }

    /**
     * 获取当天时间起止
     *
     * @param int $timeType 1=时间戳（秒），2=时间格式（年-月-日 时:分:秒）
     * @return array
     */
    public function getNowDay($timeType = 1)
    {

        // 当天0点
        $startDate = date("Y-m-d 00:00:00");
        // 当天日期
        $endDate = date("Y-m-d 23:59:59");

        return [
            'start_time' => $timeType == 1 ? strtotime($startDate) : $startDate,
            'end_time'   => $timeType == 1 ? strtotime($endDate) : $endDate,
        ];
    }

    /**
     * 获取时间起止
     *
     * @param int $dateType  时间类型 1:当天，2：本周 3：本月
     * @param int $returnTimeType 返回时间格式 1：时间戳，2：时间
     * @return array
     */
    public function getDateTime($dateType = 1, $returnTimeType = 1)
    {

        if ($dateType == 1) {
            // 当天
            return $this->getNowDay($returnTimeType);
        } elseif ($dateType == 2) {
            // 本周
            return $this->getWeekTime($returnTimeType);
        } else {
            // 当月
            return $this->getMonthTime($returnTimeType);
        }
    }


    /**
     * 数组分页
     *
     * @param int $size 每页数量
     * @param int $page 页码
     * @param array $list 列表数组
     * @return array
     */
    public function pageArray($size=15, $page=1, $list=[])
    {

        // 判断当前页面是否为空 如果为空就表示为第一页面
        $page = (empty($page)) ? '1' : $page;

        $totals = count($list);
        // 计算总页数
        $countpage = ceil($totals / $size);

        if ($page > $countpage) {
            $page = $countpage;
        }

        // 计算每次分页的开始位置
        $start = ($page - 1) * $size;

        $pagedata = array_slice($list, $start, $size);

        return $pagedata;
    }


}
