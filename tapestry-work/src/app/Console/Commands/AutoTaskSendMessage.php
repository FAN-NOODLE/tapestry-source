<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/28
 * Time: 17:55
 */

namespace App\Console\Commands;

use App\Models\Task;
use App\Services\TaskService;
use Illuminate\Console\Command;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;
use App\Models\TaskMemberRelation;
use App\Library\Uc\Common;
use Illuminate\Support\Facades\Log;

class AutoTaskSendMessage extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'auto:task:send_message';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = '任务开始或者结束时发送消息提醒';

    /**
     * TaskService
     */
    protected $taskService;

    /**
     * Create a new command instance.
     */
    public function __construct()
    {
        $this->taskService = new TaskService();
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return mixed
     *
     * @throws \Prettus\Repository\Exceptions\RepositoryException
     */
    public function handle()
    {
        // 查询当前未结束的任务
        $task = Task::where('end_time', '>=', time())->get()->toArray();

        // 为空
        if (empty($task)) {
            return;
        }

        foreach ($task as $val) {
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $val['enterprise_id']]);
            // 当前企业不存在
            if(empty($enterprise['data'])){
                Log::debug('正在进行中的任务_企业id: ' . json_encode($val['enterprise_id']));
                Log::debug('正在进行中的任务_获取授权信息: ' . json_encode($enterprise));
                continue;
            }
            $qyDomain = $enterprise['data']['domain'];
            $auth = (new Auth())->authDetail(['auth_type' => 2, 'enterprise_id' => $val['enterprise_id']]);
            // 授权信息为空
            if(empty($auth['data'])){
                Log::debug('正在进行中的任务_企业id: ' . json_encode($val['enterprise_id']));
                Log::debug('正在进行中的任务_获取授权信息: ' . json_encode($auth));
                continue;
            }
            $appId = $auth['data']['app_id'];
            // 类型为营销
            $val['task_index'] = ($val['task_type'] == Task::SALES_PERFORMANCE) ? $val['task_index'] / 100 : $val['task_index'];
            // 没有周期,结束时间和开始时间是固定的，发送消息只需要知道时间设置就可以了
            if ($val['is_loop'] == Task::TASK_NO_CYCLE_STATUS) {
                // 开始之前提醒
                if ($val['is_start_task'] == Task::IS_REMIND) {
                    $this->sendMessage($val['task_start_time'], $val['start_time'], Task::REMIND_TYPE_START, $val, $qyDomain, $appId);
                }
                // 结束之前提醒
                if ($val['is_end_task'] == Task::IS_REMIND) {
                    $this->sendMessage($val['task_end_time'], $val['end_time'], Task::REMIND_TYPE_END, $val, $qyDomain, $appId);
                }
            }

            // 有周期，任务循环，中间的时间点既是开始时间又是结束时间
            if ($val['is_loop'] == Task::TASK_LOOP_STATUS) {
                // 当前时间的任务的开始时间和结束时间
                $taskTime = $this->taskService->getTaskTimeSection($val['start_time'], $val['task_cycle'], $val['end_time']);
                $startTaskTime = $taskTime['start_time'];
                $endTaskTime = $taskTime['end_time'];
                // 任务开始时间，判断当前任务区间是不是开始区间
                if ($startTaskTime == $val['start_time']) {
                    // 开始之前提醒
                    if ($val['is_start_task'] == Task::IS_REMIND) {
                        $this->sendMessage($val['task_start_time'], $val['start_time'], Task::REMIND_TYPE_START, $val, $qyDomain, $appId);
                        $this->sendMessage($val['task_start_time'], $endTaskTime, Task::REMIND_TYPE_START, $val, $qyDomain, $appId);
                    }
                    // 结束之前提醒
                    if ($val['is_end_task'] == Task::IS_REMIND) {
                        $this->sendMessage($val['task_end_time'], $endTaskTime, Task::REMIND_TYPE_END, $val, $qyDomain, $appId);
                    }
                } else {
                    if ($endTaskTime == $val['end_time']) {
                        // 结束之前提醒
                        if ($val['is_end_task'] == Task::IS_REMIND) {
                            $this->sendMessage($val['task_end_time'], $val['end_time'], Task::REMIND_TYPE_END, $val, $qyDomain, $appId);
                        }
                    } else {
                        // 开始之前提醒
                        if ($val['is_start_task'] == Task::IS_REMIND) {
                            $this->sendMessage($val['task_start_time'], $endTaskTime, Task::REMIND_TYPE_START, $val, $qyDomain, $appId);
                        }
                        // 结束之前提醒
                        if ($val['is_end_task'] == Task::IS_REMIND) {
                            $this->sendMessage($val['task_end_time'], $endTaskTime, Task::REMIND_TYPE_END, $val, $qyDomain, $appId);
                        }
                    }
                }
            }
        }
        return 'ok';
    }

    /**
     * 任务消息发送
     *
     * @param int $timeInterval
     * @param int $sendTime
     */
    protected function sendMessage($timeInterval = 0, $sendTime = 0, $type, $task, $qyDomain, $appId)
    {
        // 当前时间格式
        $nowTime = date('Y-m-d H:i', time());
        // 相加时间之后的时间格式
        $addTime = date('Y-m-d H:i', strtotime("$nowTime+{$timeInterval}minute"));
        // 开始时间
        $time = date('Y-m-d H:i', $sendTime);
        if ($addTime == $time) {
            // 当前任务人员
            $taskMembers = TaskMemberRelation::where('task_id', $task['id'])->get()->toArray();
            foreach ($taskMembers as $taskMember) {
                $memUid = $taskMember['mem_uid'];
                $params = [
                    'toUser' => $memUid,
                    'msgType' => 'miniprogram',
                    'miniprogram' => [
                        'appid' => $appId,
                        'page' => 'pages/workbench/index/Index',
                        'title' => $type == 1 ? '任务开始提醒' : '任务到期提醒',
                        'description' => $type == 1 ? '你有一项任务即将开始' : '你有一项任务即将到期',
                        'emphasis_first_item' => false
                    ],
                    'minicontent' => $this->getContentItem($type, $task, $memUid)
                ];
                Log::debug('发送给uc的消息参数: ' . json_encode($params));
                $result = (new Common())->miniProgramSendMessage($qyDomain, $params);
                Log::debug('uc发送消息返回的参数: ' . json_encode($result));

            }
        }
    }

    /**
     * 发送消息内容
     *
     * @param int $type 消息类型：1开始消息，2结束消息
     * @param array $task 任务单个
     * @param string $memUid 员工id
     * @return array
     */
    protected function getContentItem($type, $task, $memUid)
    {
        $content_item = [];
        // 开始
        if ($type == Task::REMIND_TYPE_START && $task['task_type'] != Task::BIRTHDAY_WISHES) {
            $content_item =
                [
                    [
                        'key' => '任务名称',
                        'value' => $task['task_name'],
                    ],
                    [
                        'key' => '开始日期',
                        'value' => date('Y-m-d H:i:s', $task['start_time']),
                    ],
                    [
                        'key' => '任务周期',
                        'value' => Task::TASK_CYCLE[$task['task_cycle']],
                    ],
                    [
                        'key' => '考核指标',
                        'value' => $task['task_index']
                    ]

                ];
        }
        // 开始并且是生日祝福
        if ($type == Task::REMIND_TYPE_START && $task['task_type'] == Task::BIRTHDAY_WISHES) {
            $content_item =
                [
                    [
                        'key' => '任务名称',
                        'value' => $task['task_name'],
                    ],
                    [
                        'key' => '开始日期',
                        'value' => date('Y-m-d H:i:s', $task['start_time']),
                    ],
                    [
                        'key' => '任务周期',
                        'value' => Task::TASK_CYCLE[$task['task_cycle']],
                    ]
                ];
        }
        // 结束
        if ($type == Task::REMIND_TYPE_END && $task['task_type'] != Task::BIRTHDAY_WISHES) {
            $taskData = $this->taskService->getTaskData($task, $memUid);
            Log::debug('完成指标参数: ' . json_encode($taskData));
            $content_item =
                [
                    [
                        'key' => '任务名称',
                        'value' => $task['task_name']
                    ],
                    [
                        'key' => '考核指标',
                        'value' => $task['task_index']
                    ],
                    [
                        'key' => '完成进度',
                        'value' => $taskData['percentage'] . '%'
                    ]
                ];
        }
        // 结束并且是生日祝福
        if ($type == Task::REMIND_TYPE_END && $task['task_type'] == Task::BIRTHDAY_WISHES) {
            $content_item =
                [
                    [
                        'key' => '任务名称',
                        'value' => $task['task_name']
                    ]
                ];
        }

        return $content_item;
    }
}