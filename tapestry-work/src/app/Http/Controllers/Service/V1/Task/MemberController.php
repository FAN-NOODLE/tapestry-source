<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 12:15
 */

namespace App\Http\Controllers\Service\V1\Task;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\TaskMemberRelation;
use Illuminate\Http\Request;

class MemberController extends BaseController
{
    /**
     * 任务员工列表
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $taskId = $request->post('ids');
        // 限制个数
        $limit = $request->post('limit');

        // 必填字段验证
        if (!isset($taskId) && empty($taskId)) {
            error(600001);
        }

        // 获取任务相关员工
        $taskMember = TaskMemberRelation::whereIn('task_id', $taskId)
            ->when(!empty($limit),
                function ($query) use ($limit) {
                    return $query->limit($limit);
                })
            ->get(['task_id', 'mem_uid'])->toArray();

        return success(['data' => $taskMember]);
    }
}