<?php
/**
 * Created by PhpStorm.
 * User: Hui Yang
 * Date: 2020/01/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Models\MarketingPlanTask;
use App\Models\MarketingPlanVisibleRange;
use Illuminate\Support\Carbon;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;
use App\Services\Message\MessageService;
use App\Services\StoreWorkerService;

class AutoNoticeMarketingPlanTask extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:notice_marketing_plan_task';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '定时为符合营销计划任务条件的人员发送任务提醒';

    /**
     * 今天的时候戳
     * @var string
     */
    private $todayTimestamp = 0;

    public function __construct()
    {
        parent::__construct();

        $this->todayTimestamp = strtotime(date('Y-m-d'));
    }

    public function handle()
    {
        Log::info($this->description . '--开始', ['time' => Carbon::now()]);

        $tasks = MarketingPlanTask::query()
            ->where('marketing_plan_task.start_time', '<=', $this->todayTimestamp)
            ->where('marketing_plan_task.end_time', '>=', $this->todayTimestamp)
            ->select([
                'marketing_plan_task.*'
            ])
            ->orderBy('id', 'DESC')
            ->get()
            ->toArray();
        if (empty($tasks)) {
            return true;
        }
        $tasks = array_column($tasks, null, 'id');
        $planIdArr = array_column($tasks, 'plan_id');

        //将计划任务依据计划分组，并加入计划名称
        $planGroup = [];
        foreach ($tasks as $task) {
            $planGroup[$task['plan_id']][] = $task;
        }

        //查询范围内的人员
        $ranges = MarketingPlanVisibleRange::query()->whereIn('plan_id', $planIdArr)->get()->toArray();

        //存储计划下面的人员mem_uid
        //格式 [
        //      enterprise_id - plan_id => [ mem_uid1,mem_uid2 ]
        //     ]
        $planHasMemUidArr = [];
        foreach ($ranges as $ran) {
            if (empty($ran['mem_uid'])) {
                continue;
            }
            $key = $ran['enterprise_id'] . '-' . $ran['plan_id'];
            if (isset($planHasMemUidArr[$key])) {
                if (in_array($ran['mem_uid'], $planHasMemUidArr[$key])) {
                    continue;
                }
            }
            $planHasMemUidArr[$key][] = $ran['mem_uid'];
        }


        //整理用户的信息
        $allMemUidArr = [];
        foreach ($planHasMemUidArr as $key => $item) {
            try {
                [$enterpriseId, $planId] = explode('-', $key);
                $workerList = StoreWorkerService::list([
                        'enterprise_id' => $enterpriseId,
                        'mem_uids' => $item
                    ], false)['list'] ?? [];
                if (empty($workerList)) {
                    continue;
                }
            } catch (\Exception $e) {
                continue;
            }
            foreach ($workerList as $worker) {
                $memKey = $enterpriseId . '-' . $worker['mem_uid'];
                // 整合用户信息以及任务列表
                // 如果有门店，获取最近登录过的门店展示，但是要将任务列表加入到task_list中
                if (isset($allMemUidArr[$memKey])
                    && $allMemUidArr[$memKey]['last_login'] >= $worker['recently_logged_time']
                ) {
                    $allMemUidArr[$memKey]['task_list'] = array_merge(
                        $allMemUidArr[$memKey]['task_list'] ?? [],
                        $planGroup[$planId] ?? []
                    );
                    continue;
                }
                $allMemUidArr[$memKey] = [
                    'enterprise_id' => $enterpriseId,
                    'mem_uid' => $worker['mem_uid'],
                    'last_login' => $worker['recently_logged_time'],
                    'task_list' => array_merge(
                        isset($allMemUidArr[$memKey]) ? $allMemUidArr[$memKey]['task_list'] : [],
                        $planGroup[$planId] ?? []
                    ),
                    'store_id' => $worker['store_id'],
                ];
            }

        }
        //查询用户可以使用的计划
        foreach ($allMemUidArr as $k => $mem) {
            if (empty($mem['task_list'])) {
                continue;
            }
            $taskListIdArr = array_column($mem['task_list'], 'id');
            $mem['task_list'] = array_combine($taskListIdArr, $mem['task_list']);
            krsort($mem['task_list']);

            $firstItem = array_first($mem['task_list']);
            $warningText = $firstItem['task_title'];
            $taskCount = count(array_unique($taskListIdArr));
            if ($taskCount > 1) {
                $warningText .= "等" . $taskCount . '个';
            }
            $memUidArr = [
                $mem['mem_uid']
            ];
            $this->sendTaskStartMessage($mem['enterprise_id'], $memUidArr, [
                'body' => $warningText,
                'store_id' => $mem['store_id']
            ]);
        }

        Log::info($this->description . '--结束', ['time' => Carbon::now()]);
    }

    /**
     * 发送分销订单结算消息提醒
     * @param $enterpriseId
     * @param $sendOrderMsg array 结算订单数据
     * @param $settlementStatus integer 结算状态
     * @return bool
     * @throws \Exception
     * @author Calm
     * @time 2020-10-19 15:04:40
     */
    public function sendTaskStartMessage($enterpriseId, $memUidArr, array $params)
    {
        $content =
            [
                [
                    'key' => '今日营销任务',
                    'value' => $params['body'],
                ],
                [
                    'key' => '说明',
                    'value' => '点击查看和获取你今天的营销任务素材哦~',
                ],
            ];
        $title = '导购营销任务提醒';
        $page = 'pages/index/index?store_id=' . $params['store_id'];
        try {
            $result = (new MessageService())->sendMiniMessage($enterpriseId, $content, $page, $title, $memUidArr);

            Log::info('导购营销任务提醒发送记录', [
                'enterprise_id' => $enterpriseId,
                'store_id' => $params['store_id'],
                'mem_uid' => $memUidArr
            ]);
        } catch (\Exception $e) {
            return false;
        }
        return $result;
    }
}
