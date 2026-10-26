<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 11:11
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Task;
use App\Models\TaskMemberRelation;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 删除任务
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
        // 删除
        try {
            Task::where('id', $taskId)->delete();
            TaskMemberRelation::where('task_id', $taskId)->delete();
        } catch (\Exception $e) {
            error(200006);
        }

        return success(['data' => []]);
    }
}