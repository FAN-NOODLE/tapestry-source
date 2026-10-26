<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 11:37
 */

namespace App\Http\Controllers\Service\V1\Task;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Task;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 任务详情
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $taskId = $request->post('id');
        // 必填字段验证
        if (!isset($taskId) && empty($taskId)) {
            error(600001);
        }
        try {
            // 任务详情
            $task = Task::where('id', $taskId)->first(['id', 'task_name', 'task_type', 'start_time', 'end_time', 'task_cycle', 'is_loop', 'task_index', 'is_start_task', 'is_end_task', 'task_start_time', 'task_end_time'])->toArray();
        } catch (\Exception $e) {
            error(200007);
        }
        // 转化名称
        $task['task_type_name'] = Task::TASK_TYPE[$task['task_type']];
        $task['task_cycle'] = Task::TASK_CYCLE[$task['task_cycle']];
        $task['task_index'] = ($task['task_type'] == Task::SALES_PERFORMANCE) ? $task['task_index'] / 100 : $task['task_index'];

        return success([
            'data' => $task
        ]);
    }
}