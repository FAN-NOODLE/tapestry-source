<?php
/**
 * Created By PhpStorm
 * @Date: 2021/5/20
 * @Time: 11:21
 * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
 * @File: SyncNoticeMarketingPlanTask.php
 */

namespace App\Jobs;


use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanTaskProgress;
use App\Models\MarketingPlanVisibleRange;
use App\Models\StoreWorker;
use App\Services\Message\MessageService;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use Illuminate\Support\Carbon;
use Illuminate\Support\Facades\Log;

class SyncNoticeMarketingPlanTask extends Job
{
    protected $data;

    public $queue = 'sync_notice_marketing_plan_task';

    public function __construct($data)
    {
        $this->data = $data;
    }

    public function handle()
    {
        Log::info('后台添加任务异步发送任务--队列开始', ['$data' => $this->data, 'time' => Carbon::now()]);

        // 根据传过来的数据获取当前任务
        $task = MarketingPlanTask::query()
            ->where('id', $this->data['id'])
            ->where('deleted_at', 0)
            ->first();

        if (!$task) {
            Log::info('任务不存在或已删除', ['time' => Carbon::now()]);
            return;
        }

        $task = $task->toArray();

        // 获取企业信息
        try {
            $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $task['enterprise_id']])['data'];
        } catch (\Exception $e) {
            Log::info('获取企业相关信息失败', [
                'enterprise_id' => $task['enterprise_id']
            ]);
            return ;
        }

        // 查询出当前任务可见范围
        $ranges = MarketingPlanVisibleRange::query()->where('task_id', $task['id'])->get()->toArray();

        $memUidArr = array_unique(array_values(array_column($ranges, 'mem_uid')));
        if ($memUidArr) {
            $workerList = StoreWorkerService::list([
                    'enterprise_id' => $task['enterprise_id'],
                    'mem_uids'      => $memUidArr,
                    'worker_status' => StoreWorker::IN_SERVICE
                ], false)['list'] ?? [];

            // 查询已经完成的导购
            $finishedMemUids = MarketingPlanTaskProgress::query()
                ->where('enterprise_id', $task['enterprise_id'])
                ->where('task_id', $task['id'])
                ->pluck('mem_uid')
                ->toArray();

            foreach ($workerList as $value) {
                if ($finishedMemUids && in_array($value['mem_uid'], $finishedMemUids)) {
                    continue;
                }

                $this->sendTaskStartMessage($task['enterprise_id'], [$value['mem_uid']], [
                    'body'     => $task['task_title'],
                    'store_id' => $value['store_id'],
                    'task_id'  => $this->data['id'],
                    'corp_id'  => $enterpriseInfo['corp_id']
                ]);
            }
        }

        Log::info('后台添加任务异步发送任务--队列结束', ['time' => Carbon::now()]);
    }


    /**
     * @param $enterpriseId
     * @param $memUidArr
     * @param array $params
     * @return bool
     * @Date: 2021/5/20
     * @Time: 15:13
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function sendTaskStartMessage($enterpriseId, $memUidArr, array $params)
    {
        $content =
            [
                [
                    'key'   => '今日营销任务',
                    'value' => $params['body'],
                ],
                [
                    'key'   => '说明',
                    'value' => '点击查看和获取你今天的营销任务素材哦~',
                ],
            ];
        $title   = '导购营销任务提醒';

        $page = 'pages/task-detail/task-detail?id=' . $params['task_id'] . '&corp_id=' . $params['corp_id'];
        try {
            $result = (new MessageService())->sendMiniMessage($enterpriseId, $content, $page, $title, $memUidArr);

            Log::info('导购营销任务提醒发送记录', [
                'enterprise_id' => $enterpriseId,
                'store_id'      => $params['store_id'],
                'mem_uid'       => $memUidArr
            ]);
        } catch (\Exception $e) {
            return false;
        }
        return $result;
    }
}