<?php
/**
 * 获取任务详情
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MarketingPlanTask;
use App\Models\StoreWorker;
use App\Services\MarketingPlanService;
use App\Services\StoreWorkerService;

class TaskDetailController extends BaseController
{
    /**
     * 获取任务详情
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 16:48
     * @return array
     */
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required',
            'task_id' => 'required|min:1'
        ];

        $message = [
            'enterprise_id.*' => '企业ID不能为空',
            'mobile.*' => '导购手机号不能为空',
            'task_id.required' => '任务ID不能为空',
            'task_id.min' => '任务不存在'
        ];
        $this->validate($this->request, $rules, $message);
        $input = $this->request->all([
            'enterprise_id',
            'mobile',
            "domain",
            "task_id",
            'store_id'
        ]);

        //查询导购信息
        $workerInfo = StoreWorkerService::detail([
            'enterprise_id' => $input['enterprise_id'],
            'mobile' => $input['mobile'],
            'store_id' => $input['store_id'],
            'worker_status' => StoreWorker::IN_SERVICE
        ]);
        $input['worker_id'] = $workerInfo['id'];
        $input['mem_uid'] = $workerInfo['mem_uid'];

        //最小为6个月前的1号
        $monthAgo = MarketingPlanTask::QUERY_BETWEEN_MONTH_AGO;
        $monthLater = MarketingPlanTask::QUERY_BETWEEN_MONTH_LATER + 1;
        $currentMonth = date('Y-m-1');

        $input['start_time'] = strtotime("{$currentMonth} -{$monthAgo} month");
        $input['end_time'] = strtotime("{$currentMonth} +{$monthLater} month -1 day");

        $marketingPlanService = new MarketingPlanService();

        //检查海报情况
        $marketingPlanService->checkSuncodeBanner(
            $input['enterprise_id'],
            $input['task_id'],
            $workerInfo
        );

//        $input['plan_ids'] = $marketingPlanService->getRangePlanList(
//            $input['enterprise_id'],
//            $workerInfo['id']
//        );

        $result = $marketingPlanService->getTaskDetail($input);

        return success($result ?? []);
    }
}
