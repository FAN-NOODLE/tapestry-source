<?php
/**
 * 获取营销日历
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MarketingPlanTask;
use App\Models\StoreWorker;
use App\Services\MarketingPlanService;
use App\Services\StoreWorkerService;

class TaskCalendarController extends BaseController
{
    /**
     * 获取营销日历
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 17:07
     */
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required'
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
            'mobile.*' => '导购手机号不能为空'
        ];
        $this->validate($this->request, $rules, $message);
        $input = $this->request->all([
            'enterprise_id',
            'current_date',
            "domain",
            'store_id',
            'mobile'
        ]);

        //查询导购信息
        $workerInfo = StoreWorkerService::detail([
            'enterprise_id' => $input['enterprise_id'],
            'mobile' => $input['mobile'],
            'store_id' => $input['store_id'],
            'worker_status' => StoreWorker::IN_SERVICE
        ]);

        //6个月前
        $todayMonth = date('Y-m-1');
        $monthAgo = MarketingPlanTask::QUERY_BETWEEN_MONTH_AGO;
        $sixMonthAgo = strtotime("{$todayMonth} -{$monthAgo} month");
        //后一个 月
        $monthLater = MarketingPlanTask::QUERY_BETWEEN_MONTH_LATER + 1;
        $aMonthLater = strtotime("{$todayMonth} +{$monthLater} month -1 day");

        $input['start_time'] = $sixMonthAgo;
        $input['end_time'] = $aMonthLater;

        $marketingPlanService = new MarketingPlanService();
        $input['plan_ids'] = $marketingPlanService->getRangePlanList(
            $input['enterprise_id'],
            $workerInfo['id']
        );

        $result = $marketingPlanService->getMarketingPlanTaskCalendarList($input);
        $result = $marketingPlanService->formatDayData(
            $result,
            $input['start_time'],
            $input['end_time'],
            $input['current_date'],
            false
        );
        $result = $marketingPlanService->fullFormatCalendar($result);
        return success($result);
    }
}
