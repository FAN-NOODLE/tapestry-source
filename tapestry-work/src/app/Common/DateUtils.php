<?php


namespace App\Common;


class DateUtils
{

    const DAY_SECONDS = 86400;

    private function __construct()
    {
    }

    /**
     * 时间戳转换为0点时间戳
     * @param int $time
     * @return int
     */
    public static function timeToStatTime(int $time): int
    {
        return strtotime(date("Y-m-d", $time));
    }

    public static function todayStatTime()
    {
        return self::timeToStatTime(time());
    }

    public static function tomorrowStatTime()
    {
        return self::todayStatTime() + self::DAY_SECONDS;
    }

    public static function yesterdayStatTime()
    {
        return self::todayStatTime() - self::DAY_SECONDS;
    }
}