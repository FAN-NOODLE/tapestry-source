<?php
/**
 * 任务完成上报
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskProgress;
use App\Models\StoreWorker;
use App\Services\MarketingPlanService;
use App\Services\StoreWorkerService;
use Illuminate\Validation\Rule;

class TaskFinishController extends BaseController
{
    /**
     * 任务完成上报
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
            'task_id' => 'required|min:1',
            'finish_type' => [
                'required',
                Rule::in([
                    MarketingPlanTaskProgress::FINISHED_TYPE_COPY,
                    MarketingPlanTaskProgress::FINISHED_TYPE_FORWARD,
                    MarketingPlanTaskProgress::FINISHED_TYPE_SHARE,
                    MarketingPlanTaskProgress::FINISHED_TYPE_DOWNLOAD
                ])
            ]
        ];

        $message = [
            'enterprise_id.*' => '企业ID不合法',
            'mobile.*' => '导购手机号不能为空',
            'task_id.*' => '任务参数错误',
            'finish_type.*' => '完成类型参数错误'
        ];
        $this->validate($this->request, $rules, $message);
        $input = $this->request->all([
            'enterprise_id',
            'mobile',
            "domain",
            "task_id",
            'store_id',
            'finish_type'
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
        $monthLater = MarketingPlanTask::QUERY_BETWEEN_MONTH_LATER;
        $currentMonth = date('Y-m-1');

        $input['start_time'] = strtotime("{$currentMonth} -{$monthAgo} month");
        $input['end_time'] = strtotime("{$currentMonth} +{$monthLater} month");

        $marketingPlanService = new MarketingPlanService();
        $input['plan_ids'] = $marketingPlanService->getRangePlanList(
            $input['enterprise_id'],
            $workerInfo['id']
        );

        $todayTimestamp = strtotime(date('Y-m-d'));

        $result = $marketingPlanService->getTaskDetail($input);
        if (empty($result)) {
            return success();
        }
        //判断是否在时间范围内，如果不在时间范围内，则不写入数据库
//        if ($result['start_time'] > $todayTimestamp || $result['end_time'] < $todayTimestamp) {
//            return success();
//        }

        //写入记录
        $data = [
            'enterprise_id' => $input['enterprise_id'],
            'task_id' => $input['task_id'],
            'worker_id' => $workerInfo['id']
        ];
        $finishObj = MarketingPlanTaskProgress::query()->firstOrNew($data);
        $data['mem_uid'] = $workerInfo['mem_uid'];
        $data['finished_type'] = $input['finish_type'];
        $finishObj->fill($data);
        $finishObj->save();

        return success();
    }
}
