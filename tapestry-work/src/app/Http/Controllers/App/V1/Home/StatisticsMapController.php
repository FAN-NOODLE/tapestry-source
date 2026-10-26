<?php

namespace App\Http\Controllers\App\V1\Home;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\CustomerStatisticsService;
use App\Services\OrderStatisticsService;
use Illuminate\Http\Request;
use Illuminate\Support\Carbon;

class StatisticsMapController extends BaseController
{
    const PERIOD_SHORT  = 1; // 近7天
    const PERIOD_MEDIUM = 2; // 近30天
    const PERIOD_LONG   = 3; // 近90天

    /**
     * @var CustomerStatisticsService
     */
    protected $customerStatisticsService;
    /**
     * @var OrderStatisticsService
     */
    protected $orderStatisticsService;

    /**
     * StatisticsMapController constructor.
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
     * 统计数据趋势图
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $period = $request->get('period');

        if (empty($period)) {
            throw new Error(111001, '参数错误');
        }

        switch ($period) {
            case self::PERIOD_SHORT:
                $startTime = Carbon::now()->subDay(6)->startOfDay()->timestamp;
                break;
            case self::PERIOD_MEDIUM:
                $startTime = Carbon::now()->subDay(29)->startOfDay()->timestamp;
                break;
            case self::PERIOD_LONG:
                $startTime = Carbon::now()->subDay(89)->startOfDay()->timestamp;
                break;
            default:
                $startTime = Carbon::now()->subDay(1)->startOfDay()->timestamp;
                break;
        }

        $endTime = Carbon::now()->subDay(1)->endOfDay()->timestamp;

        $dayArr = [];
        for ($i = $startTime; $i < $endTime; $i += 86400) {
            $dayArr[] = date('m-d', $i);
        }

        $statisticsMap = [];

        if ($this->memUid) {
            // 客户统计数据趋势图
            $userStatisticsMap = $this->customerStatisticsService->statisticsMap($this->enterpriseId, $this->memUid, $startTime, $endTime, $dayArr, $period);
            // 订单统计数据趋势图
            $orderStatisticsMap = $this->orderStatisticsService->statisticsMap($this->enterpriseId, $this->memUid, $startTime, $endTime, $dayArr, $period);

            $statisticsMap = array_merge($userStatisticsMap, $orderStatisticsMap);
        }

        return success(['data' => $statisticsMap]);
    }
}
