<?php

namespace App\Http\Controllers\App\V1\Task\Outsiders;

use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\Task;
use Illuminate\Http\Request;

class SendController extends TaskBaseController
{
    /**
     * 给外部人员提醒
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $taskId       = $request->get('task_id', '');
        $customerId   = $request->get('ec_userid', '');
        $displayState = $request->get('display_state', '1');

        // 任务ID
        if (empty($taskId)) {
            error(200010);
        }

        // 外部人员ID不能为空
        if (empty($customerId)) {
            error(200008);
        }

        // 写入或更新日志记录
        $this->saveLog($taskId, $customerId, Task::CUSTOMER_MANAGEMENT, $displayState);

        return success(['data' => []]);
    }
}

