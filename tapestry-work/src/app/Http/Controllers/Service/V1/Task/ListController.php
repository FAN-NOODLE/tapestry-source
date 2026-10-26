<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 10:42
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Task;
use App\Services\TaskService;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    /**
     * 任务列表
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 接受查询参数
        $page = $request->post('page', 1);
        $size = $request->post('size', 10);
        $search = $request->post('search', []);
        $storeRole = $request->post('store', []);
        $enterpriseId = $request->post('enterprise_id');

        // 通过权限获取任务列表
        [$total, $data] = (new TaskService())
                            ->getTaskListByPermission($page, $size, $search, $enterpriseId, $storeRole);
        // 格式化数据
        foreach ($data as &$val) {
            $val['task_type_name'] = Task::TASK_TYPE[$val['task_type']];
            $val['task_cycle'] = Task::TASK_CYCLE[$val['task_cycle']];
            // 任务状态
            if (time() < $val['start_time']) {
                $val['status'] = Task::TASK_STATUS_START;
            }
            if (time() >= $val['start_time'] && time() <= $val['end_time']) {
                $val['status'] = Task::TASK_STATUS_HAND;
            }
            if (time() > $val['end_time']) {
                $val['status'] = Task::TASK_STATUS_END;
            }
        }

        return success([
            'data' => $data,
            'meta' => [
                'total' => (int)$total,
                'page' => $page,
                'size' => $size
            ]
        ]);
    }
}
