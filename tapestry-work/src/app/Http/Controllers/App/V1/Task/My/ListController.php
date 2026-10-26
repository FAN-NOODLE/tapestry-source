<?php

namespace App\Http\Controllers\App\V1\Task\My;

use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\Task;
use App\Models\TaskMemberRelation;
use App\Services\TaskService;
use App\Common\Constant;

class ListController extends TaskBaseController
{
    /**
     * 我的任务列表
     *
     * @return array
     */
    public function index()
    {

        // 获取当前用户UID
        $memUid = $this->memUid;
        // 获取企业ID
        $enterpriseId = $this->enterpriseId;

        // 查询用户所在任务的数据
        $taskMemberRelationData = TaskMemberRelation::query()->where('mem_uid', $memUid)
            ->where('enterprise_id', $enterpriseId)
            ->get(['task_id'])->toArray();

        // 初始化进度条
        $progress_bar = [
            'task_name' => '',
            'percentage' => ''
        ];

        $data         = [];
        // 客户管理
        $custom = Task::CUSTOMER_MANAGEMENT;
        // 如果当前用户存在任务
        if (!empty($taskMemberRelationData)) {
            // 获取任务ID集合
            $taskIds = array_column($taskMemberRelationData, 'task_id');

            // 获取我可以参与的任务数据
            $taskData = Task::query()->whereIn('id', $taskIds)->where('enterprise_id', $enterpriseId)
                ->when($this->storeStatus == Constant::STORE_CLOSE, function ($query) use ($custom) {
                    return $query->where('task_type', $custom);
                })
                ->get()->toArray();

            $serviceTask = new TaskService();

            // 遍历数据排除不是今天需要完成的任务
            foreach ($taskData as $v) {
                $taskList = $serviceTask->getTaskData($v, $memUid);
                if (!empty($taskList)) {
                    // 判断返回的数据
                    if (Task::BIRTHDAY_WISHES != $taskList['task_type'] && empty($progress_bar['task_name'])) {
                        $progress_bar = [
                            'task_name' => $taskList['task_name'],
                            'percentage' => $taskList['percentage']
                        ];
                    }

                    $data[] = $taskList;
                }
            }
        }

        return success(['data' => $data, 'progress_bar' => $progress_bar]);
    }
}
