<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/23
 * Time: 18:09
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\Task;
use App\Models\TaskMemberRelation;
use App\Models\TaskType;
use CrsSdk\Boss\Enterprise;
use Illuminate\Http\Request;

class AddController extends BaseController
{
    /**
     * 添加任务
     *
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        // 初始化字段
        $taskMemberData = [];
        // 接受查询参数
        $data = $request->post();
        // 员工id
        $munUid = $data['mem_uid'];
        unset($data['mem_uid']);
        // 金额单位元变成分
        if ($data['task_type'] == Task::SALES_PERFORMANCE) {
            $data['task_index'] = $data['task_index'] * 100;
        }
        // 组装任务类型 属性serve_type
        $data['serve_type'] = $this->getServeType($data['task_type']);

        // 添加任务
        try {
            $task = Task::query()->create($data);
        } catch (\Exception $e) {
           error(200005);
        }
        // 任务id
        $taskId = $task->id;
        // 获取企业domain
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' =>$data['enterprise_id']]);

        // 过滤当前员工id，标示门店
        $memList = (new Member())->list($enterpriseInfo['data']['domain'], 'contact', ['memUids' => $munUid])['list'];
        $memList = array_combine_by_key($memList, 'memUid');

        // 员工关联表
        foreach ($munUid as $key => $value) {
            // 拼接门店信息入库
            $taskMemberData[$key]['store_id'] = '';
            if (isset($memList[$value]['dpName'][0]['dptName'])
                && $memList[$value]['dpName'][0]['dpId']
                && $memList[$value]['dpName'][0]['dptName'] === '门店') {
                $taskMemberData[$key]['store_id'] = $memList[$value]['dpName'][0]['dpId'];
            }
            $taskMemberData[$key]['task_id'] = $taskId;
            $taskMemberData[$key]['mem_uid'] = $value;
            $taskMemberData[$key]['enterprise_id'] = (int)$data['enterprise_id'];
            $taskMemberData[$key]['created_at'] = time();
        }

        // 员工入库
        try {
            TaskMemberRelation::query()->insert($taskMemberData);
        } catch (\Exception $e) {
            error(200005);
        }

        return success(['data' => ['id' => (int)$taskId], 'meta' => []]);
    }

    /**
     * 获取任务类型 属性serve_type
     * @param $taskTypeId
     * @return int
     */
    public function getServeType($taskTypeId)
    {
        // 初始化服务类型 为shop+crs
        $serveType = 0;
        $taskType = TaskType::query()->find($taskTypeId)->toArray();
        if (!empty($taskType)) {
            $serveType = $taskType['serve_type'];
        }
        return $serveType;
    }
}
