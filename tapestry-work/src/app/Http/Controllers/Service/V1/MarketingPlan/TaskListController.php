<?php
/**
 * 获取任务列表
 */

namespace App\Http\Controllers\Service\V1\MarketingPlan;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskProgress;
use App\Models\StoreWorker;
use App\Services\MarketingPlanService;
use App\Services\StoreWorkerService;

class TaskListController extends BaseController
{
    //查询任务列表范围为前后两天
    private $dayLimit = 2;

    /**
     * 获取营销任务列表
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/12 16:48
     * @return array
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
            'mobile',
            'page',
            'size',
            'default_date',  //请求列表时定位的默认中间位置标识
            'task_type',
            'task_title',
            'is_finished',
            'is_index',
            'start_time',
            'end_time'
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
        $input['store_id'] = $workerInfo['store_id'];

        //处理分页参数
        if (empty($input['page']) || $input['page'] <= 0) {
            $input['page'] = 1;
        }
        if (empty($input['size']) || $input['size'] <= 0) {
            $input['size'] = 10;
        }

        $currentDatetime = date('Y-m-d');
        if ($input['current_date'] > 0) {
            $currentDatetime = date('Y-m-d', $input['current_date']);
        }

        //获取列表时记录被选中的日期
        if ($input['default_date'] > 0) {
            $input['default_date'] = date('Y-m-d', $input['default_date']);
        } else {
            $input['default_date'] = $currentDatetime;
        }
        //确定首页已经的时间范围
        $betweenNum = $this->dayLimit;
        if($input['is_index']){
            $input['start_time'] = strtotime("{$input['default_date']} -{$betweenNum} day");
            $input['end_time'] = strtotime("{$input['default_date']} +{$betweenNum} day");

            $paramTimestamp = strtotime($input['default_date']);
            if (!$this->processDateRange(
                $paramTimestamp,
                $input['start_time'],
                $input['end_time']
            )
            ) {
                error(1200001);
            }
        }else{
            $input['start_time'] = isset($input['start_time']) ? $input['start_time'] : strtotime(date('Y-m-d'));
            $input['end_time'] = isset($input['end_time']) ? $input['end_time'] : strtotime(date('Y-m-d 23:59:59'));
        }
        $input['current_date'] = strtotime($currentDatetime);

        $marketingPlanService = new MarketingPlanService();

//        $input['plan_ids'] = $marketingPlanService->getRangePlanList(
//            $input['enterprise_id'],
//            $workerInfo['id']
//        );

        $input['task_ids'] = array_filter($marketingPlanService->getRangeTaskList(
            $input['enterprise_id'],
            $workerInfo['id']
        ));

        // 完成的任务
        $completed = array_filter($marketingPlanService->getCompleteInfoByTaskIds($input['task_ids'], $workerInfo['id']));

        // 首页
        if($input['is_index']){
            $input['task_ids'] = array_diff($input['task_ids'], $completed);
        }else{
            $input['task_ids'] = array_merge(array_diff($input['task_ids'], $completed), $completed);
            if($input['is_finished'] == MarketingPlanTaskProgress::TASK_NOT_FINISHED){
                $input['task_ids'] = array_diff($input['task_ids'], $completed);
            }else if($input['is_finished'] == MarketingPlanTaskProgress::TASK_ALREADY_FINISHED){
                $input['task_ids'] = $completed;
            }
        }

        $input['task_ids'] = array_unique($input['task_ids']);

        [$result, $total] = $marketingPlanService->getMarketingPlanTaskList($input);

        foreach ($result as &$item) {
            foreach ($item['resource'] as &$res) {
                //检查海报
//                if (empty($res['suncode_url']) || empty($res['suncode_mediaid'])) {
//                    $returnData = $marketingPlanService->checkSuncodeBanner(
//                        $input['enterprise_id'],
//                        $item['id'],
//                        $workerInfo
//                    );
//                    if (!empty($returnData)) {
//                        $res['suncode_url'] = $returnData['background_url'];
//                        $res['suncode_mediaid'] = $returnData['suncode_mediaid'];
//                    }
//                }
                //检查资源media_id
                if (empty($res['mediaid'])) {
                    $returnData = $marketingPlanService->checkResourceMediaId(
                        $input['enterprise_id'],
                        $res['id'],
                        $workerInfo
                    );
                    if (!empty($returnData)) {
                        $res['mediaid'] = $returnData['mediaid'] ?? '';
                    }
                }
            }
        }

        // 获取任务数量
        if(!$input['is_index']){
            $return['page']        = $input['page'];
            $return['size']        = $input['size'];
            $return['total']       = $total;
            $return['list']        = array_values($result);
            $return['task_number'] = [
                'all_num'               => $marketingPlanService->getTaskNumber(array_merge($input, ['task_type' => 0])),
                'mass_distribution_num' => $marketingPlanService->getTaskNumber(array_merge($input, ['task_type' => MarketingPlanTask::TYPE_GROUP_TASK])),
                'train_num'             => $marketingPlanService->getTaskNumber(array_merge($input, ['task_type' => MarketingPlanTask::TYPE_TRAINING_TASK])),
                'invitation_num'        => $marketingPlanService->getTaskNumber(array_merge($input, ['task_type' => MarketingPlanTask::TYPE_INVITATION_TASK])),
            ];

            return success($return);
        }

        $result = $marketingPlanService->formatDayData(
            $result,
            $input['start_time'],
            $input['end_time'],
            $input['current_date'],
            true
        );

        foreach ($result as $key => &$val) {
            $val['page'] = $input['page'];
            $val['size'] = $input['size'];
            $val['is_current'] = false;
            $val['total'] = 0;

            if ($key != $currentDatetime) {
                $val['task_list'] = [];
                continue;
            }
            $val['is_current'] = true;
            $val['total'] = $total;
        }

        return success(array_values($result));
    }

    /**
     * 处理开始时间与结束时间，如果选中的是6个月前的1号，则显示1，2，3，4，5，如果是1个月后的最后一天，则显示27，28，29，30，31
     * User: Hui Yang
     * Email: 2361313833@qq.com
     * DateTime: 2021/1/26 16:38
     * @param $currentDatetime string 选中的时间戳
     * @param $startTime string 开始时间戳
     * @param $endTime string 结束时间戳
     */
    public function processDateRange($currentTimestamp, &$startTimestamp, &$endTimestamp)
    {
        $aDaySeconds = 86400;
        $betweenNumSeconds = $this->dayLimit * $aDaySeconds;
        //最小为6个月前的1号
        $monthAgo = MarketingPlanTask::QUERY_BETWEEN_MONTH_AGO;
        $minTimestamp = strtotime(date('Y-m-1') . " -{$monthAgo} month");
        $monthLater = MarketingPlanTask::QUERY_BETWEEN_MONTH_LATER + 1;
        $maxTimestamp = strtotime(date('Y-m-1') . " +{$monthLater} month -1 day");

        //验证时间范围有效性，如果超出范围，抛出错误
        if ($currentTimestamp < $minTimestamp || $currentTimestamp > $maxTimestamp) {
            return false;
        }

        //开始结束时间同加同减同一个数
        $calcMinSeconds = $currentTimestamp - $minTimestamp;
        if ($calcMinSeconds <= $betweenNumSeconds) {
            $startTimestamp = $currentTimestamp - $calcMinSeconds;
            $endTimestamp = $currentTimestamp + $betweenNumSeconds + ($betweenNumSeconds - $calcMinSeconds);
        }
        //最大为一个月后最后一天
        $calcMaxSeconds = $maxTimestamp - $currentTimestamp;
        if ($calcMaxSeconds <= $betweenNumSeconds) {
            $startTimestamp = $currentTimestamp - $betweenNumSeconds - ($betweenNumSeconds - $calcMaxSeconds);
            $endTimestamp = $currentTimestamp + $calcMaxSeconds;
        }
        return true;
    }
}
