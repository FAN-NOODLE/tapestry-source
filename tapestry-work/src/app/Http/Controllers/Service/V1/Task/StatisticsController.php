<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/12/4
 * Time: 9:58
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Task;
use App\Models\TaskMemberRelation;
use App\Services\TaskService;
use Illuminate\Http\Request;

class StatisticsController extends BaseController
{
    /**
     * 任务完成度的统计接口
     *
     * @param Request $request
     * @param TaskService $taskService
     * @return \App\Http\Resources\Success|array|mixed|null
     */
    public function index(Request $request, TaskService $taskService)
    {
        // 初始化
        $data = [];
        $where = [];
        $result = [];
        // 接受查询参数
        $data['enterprise_id'] = $request->post('enterprise_id');
        $data['task_type'] = $request->post('task_type');
        // 拼凑条件
        array_push($where, ['enterprise_id', '=', $data['enterprise_id']]);

        if (isset($data['task_type']) && !empty($data['task_type'])) {
            array_push($where, ['task_type', '=', $data['task_type']]);
        }
        // 当前时间,正在执行的任务
        $nowTime = time();
        array_push($where, ['start_time', '<=', $nowTime]);
        array_push($where, ['end_time', '>=', $nowTime]);

        $tasks = Task::where($where)->get(['id', 'start_time', 'end_time', 'task_cycle', 'task_index', 'is_loop', 'task_type']);

        // 人员统计
        foreach ($tasks as $task) {
            // 循环执行，获取当前时间的任务开始时间以及结束时间
            if ($task['is_loop'] == \App\Common\Task::IS_LOOP) {
                $taskTime = $taskService->getTaskTimeSection($task['start_time'], $task['task_cycle'], $task['end_time']);
                $task['start_time'] = $taskTime['start_time'];
                // 任务结束都是23：59：59
                $task['end_time'] = $taskTime['end_time'] - 1;
            }
            // 任务执行员工
            $taskMember = TaskMemberRelation::query()->select(['mem_uid', 'store_id'])->where('task_id', $task['id'])->get()->toArray();
            $memUid = array_column($taskMember, 'mem_uid');
            $result[] = [
                'id' => $task['id'],
                'start_time' => $task['start_time'],
                'end_time' => $task['end_time'],
                'task_index' => $task['task_index'],
                'task_type' => $task['task_type'],
                'task_member' => $memUid,
                'store_id' => $taskMember['store_id'] ?? '',
                'task_count' => count($memUid) * $task['task_index'],
            ];
        }

        return success(['data' => $result]);
    }
}