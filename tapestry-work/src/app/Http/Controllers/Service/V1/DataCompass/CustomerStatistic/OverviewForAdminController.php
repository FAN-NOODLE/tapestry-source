<?php


namespace App\Http\Controllers\Service\V1\DataCompass\CustomerStatistic;


use App\Common\DateUtils;
use App\Services\DataStatistic\CustomerDataCompassService;

class OverviewForAdminController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'qy_domain' => 'required|string',
            'start_time' => 'required|integer',
            'end_time' => 'required|integer',
            'worker_ids' => 'array',
        ];
        $this->validate($this->request, $rules);
        $params['enterprise_id'] = $this->request->input('enterprise_id');
        $params['qy_domain'] = $this->request->input('qy_domain');
        $params['start_time'] = $this->request->input('start_time');
        $params['end_time'] = $this->request->input('end_time');
        $params['worker_ids'] = $this->request->input('worker_ids', []);

        $params['start_stat_time'] = DateUtils::timeToStatTime($params['start_time']);
        $params['end_stat_time'] = DateUtils::timeToStatTime($params['end_time']);

        $workerList = [];
        if (!empty($params['worker_ids'])) {
            $workerList = $this->listServicingWorker($params['enterprise_id'], $params['worker_ids']);
            if (empty($workerList)) {
                return success([
                    'overview' => CustomerDataCompassService::DEFAULT_STATISTIC,
                    'detail_list' => []
                ]);
            }
        }

        $overview = $this->customerDataCompassService->overview(
            $params['enterprise_id'],
            $params['start_stat_time'],
            $params['end_stat_time'],
            $workerList
        );

        $memUidList = array_values(array_unique(array_filter(array_column($workerList, 'mem_uid'))));
        // 类型为1-截止当前选择日期之前的数据
        $memberLevelStatistics = $this->customerDataCompassService->getMemberLevelStatistics(
            $params['enterprise_id'],
            $params['start_time'],
            $params['end_time'],
            $memUidList
        );

        // 类型为2-按照时间范围筛选）
        $memberLevelStatisticsByTimeRange = $this->customerDataCompassService->getMemberLevelStatistics(
            $params['enterprise_id'],
            $params['start_time'],
            $params['end_time'],
            $memUidList,
            2
        );

        $listByDate = $this->customerDataCompassService->listByDate(
            $params['enterprise_id'],
            $params['start_stat_time'],
            $params['end_stat_time'],
            $workerList,
            '',
            0,
            0,
            0,
            false,
            false
        );
        return success([
            'overview' => $overview,
            'detail_list' => $listByDate['list'],
            'member_level_statistics' => $memberLevelStatistics['memberGrade'],
            'member_total_count' => $memberLevelStatistics['memberTotal'],
            'member_level_statistics_by_time_range' => $memberLevelStatisticsByTimeRange['memberGrade'],
            'member_total_count_by_time_range' => $memberLevelStatisticsByTimeRange['memberTotal']
        ]);
    }
}