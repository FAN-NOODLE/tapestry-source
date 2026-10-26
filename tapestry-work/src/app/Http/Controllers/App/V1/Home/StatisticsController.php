<?php

namespace App\Http\Controllers\App\V1\Home;

use App\Http\Controllers\App\V1\BaseController;
use App\Services\CustomerStatisticsService;
use App\Services\OrderStatisticsService;

class StatisticsController extends BaseController
{
    /**
     * @var CustomerStatisticsService
     */
    protected $customerStatisticsService;
    /**
     * @var OrderStatisticsService
     */
    protected $orderStatisticsService;

    /**
     * StatisticsController constructor.
     *
     * @param CustomerStatisticsService $customerStatisticsService
     * @param OrderStatisticsService    $orderStatisticsService
     */
    public function __construct(
        CustomerStatisticsService $customerStatisticsService,
        OrderStatisticsService $orderStatisticsService
    )
    {
        parent::__construct();
        $this->customerStatisticsService = $customerStatisticsService;
        $this->orderStatisticsService    = $orderStatisticsService;
    }

    /**
     * 首页统计数据
     *
     * @return array
     */
    public function index()
    {
        // 今日时间戳
        $dayStartTime = strtotime(date('Y-m-d', time()) . ' 00:00:00');
        $dayEndTime   = strtotime(date('Y-m-d', time()) . ' 23:59:59');

        // 客户统计数据
        $userStatistics = $this->customerStatisticsService->statistics($this->enterpriseId, $this->memUid, $dayStartTime, $dayEndTime);
        // 订单统计数据
        $orderStatistics = $this->orderStatisticsService->statistics($this->enterpriseId, $this->memUid, $dayStartTime, $dayEndTime);

        $statistics = array_merge($userStatistics, $orderStatistics);

        return success(['data' => $statistics]);
    }
}
