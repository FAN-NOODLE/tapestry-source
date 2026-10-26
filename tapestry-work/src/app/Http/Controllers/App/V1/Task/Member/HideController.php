<?php

namespace App\Http\Controllers\App\V1\Task\Member;

use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\Task;
use App\Models\TaskMemberRecord;
use Illuminate\Http\Request;

class HideController extends TaskBaseController
{
    /**
     * 隐藏会员招募
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $taskId     = $request->get('task_id', '');
        $customerId = $request->get('ec_userid', '');

        // 任务ID
        if (empty($taskId)) {
            error(200010);
        }

        // 外部人员ID不能为空
        if (empty($customerId)) {
            error(200008);
        }

        // 写入或更新日志记录
        $this->saveLog($taskId, $customerId, Task::MEMBER_RECRUITMENT, TaskMemberRecord::TASK_HIDE_STATUS);

        return success(['data' => []]);
    }
}

