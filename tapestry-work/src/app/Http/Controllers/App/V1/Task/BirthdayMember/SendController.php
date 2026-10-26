<?php

namespace App\Http\Controllers\App\V1\Task\BirthdayMember;

use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\Task;
use Illuminate\Http\Request;

class SendController extends TaskBaseController
{
    /**
     * 生日提醒记录
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

        // 用户ID
        if (empty($customerId)) {
            error(200011);
        }

        // 写入或更新日志记录
        $this->saveLog($taskId, $customerId, Task::BIRTHDAY_WISHES, $displayState);

        return success(['data' => []]);
    }
}

