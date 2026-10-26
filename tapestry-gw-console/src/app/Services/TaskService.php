<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/12/3
 * Time: 17:30
 */

namespace App\Services;

use App\Common\Task;

class TaskService extends BaseService
{
    /**
     * 获取结束时间
     *
     * @param string $startTime 开始时间时间格式
     * @param int $taskCycle 任务周期：1:一天，2:一周，3:一个月，4:三个月,5:六个月，6:一年
     *
     * @return false|string
     */
    public function getEndTime($startTime, $taskCycle)
    {
        $data = '';

        switch ($taskCycle) {
            case Task::ONE_DAY:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+1day"));
                break;
            case Task::ONE_WEEKS:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+1week"));
                break;
            case Task::ONE_MONTH:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+1month"));
                break;
            case Task::THREE_MONTHS:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+3month"));
                break;
            case Task::SIX_MONTHS:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+6month"));
                break;
            case Task::ONE_YEAR:
                $data = date('Y-m-d H:i:s', strtotime("$startTime+1year"));
                break;
        }

        return $data;
    }
}
