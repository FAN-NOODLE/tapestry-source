<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 14:44
 */

namespace App\Http\Controllers\V1\Task;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Task;
use Illuminate\Http\Request;
use App\Services\TaskService;

class AddController extends BaseController
{
    /**
     * 页码最大条数
     */
    const SIZE_MAX = 9999;

    /**
     * 任务添加
     *
     * @param Request $request
     * @param TaskService $taskService
     * @return array|mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request,TaskService $taskService)
    {
        // 并发锁定
        $fp = fopen("lock.txt", "r");
        if(!flock($fp,LOCK_EX | LOCK_NB)) {
            error(600005);
        }
        // 接受查询参数
        $data = $request->post();
        $rules = [
            'task_type' => 'required | integer',
            'task_name' => 'required',
            'start_time' => 'required',
            'mem_uid' => 'required',
            'task_cycle' => 'required',
            'is_loop' => 'required',
            'is_start_task' => 'required',
            'is_end_task' => 'required',
            'task_start_time' => 'required',
            'task_end_time' => 'required'
        ];
        $messages = [
            'task_type.required' => '任务类型不能为空',
            'task_name.required' => '任务名称不能为空',
            'start_time.required' => '开始时间不能为空',
            'mem_uid.required' => '员工id不能为空',
            'task_cycle.required' => '任务周期不能为空',
            'is_loop.required' => '是否循环不能为空',
            'is_start_task.required' => '任务开始之前提醒不能为空',
            'is_end_task.required' => '任务结束之前提醒不能为空',
            'task_start_time.required' => '开始提醒时间不能为空',
            'task_end_time.required' => '结束提醒时间不能为空',
        ];
        $this->validate($request, $rules, $messages);
        // 企业id
        $data['enterprise_id'] = $this->enterpriseId;
        // 时间戳毫秒转换成秒
        $data['start_time'] = $data['start_time'] / 1000;
        if (!empty($data['end_time'])) {
            $data['end_time'] = (int)($data['end_time'] / 1000);
            if ($data['start_time'] >= $data['end_time']) {
                error(420001);
            }
        }
        // 任务指标不能小于零
        if ($data['task_index'] <= 0 && $data['task_type'] != \App\Common\Task::TASK_TYPE_BIRTHDAY) {
            error(420002);
        }
        // 开始提醒时间小时换算成分钟
        if ($data['task_start_time'] < 10) {
            $data['task_start_time'] *= \App\Common\Task::SIX_MIN;
        }
        // 结束提醒时间小时换算成分钟
        if ($data['task_end_time'] < 10) {
            $data['task_end_time'] *= \App\Common\Task::SIX_MIN;
        }
        // 开始时间时间格式
        $startTime = date('Y-m-d H:i:s', $data['start_time']);
        // 任务不循环
        if ($data['is_loop'] == \App\Common\Task::NO_IS_LOOP) {
            // 结束时间为前一天的59秒
            $data['end_time'] = strtotime($taskService->getEndTime($startTime, $data['task_cycle'])) - 1;
        }
        // 检查是否重复保存 任务条件：类型，开始时间，周期，是否循环
        $search['search']['task_type'] = $data['task_type'];
        $search['enterprise_id'] = $data['enterprise_id'];
        // 时间戳毫秒转换成秒
        $search['search']['start_time'] = $data['start_time'];
        $search['search']['end_time'] = $data['end_time'];
        // 在当前时间区间内，改类型任务下的员工uid
        $search['size'] = self::SIZE_MAX;

        $taskList = (new Task())->taskList($search);
        $taskId = array_unique(array_column($taskList['data'], 'id'));
        $member = (new Task())->taskMember(['ids' => $taskId]);
        $memberIds = array_unique(array_column($member['data'], 'mem_uid'));
        $repeatMember = [];
        // 重复员工
        foreach ($data['mem_uid'] as $members) {
            if (in_array($members, $memberIds)) {
                array_push($repeatMember, $members);
            }
        }
        // 重复向外抛错
        if (!empty($repeatMember)) {
            notice(600004, '', $repeatMember);
        }

        $result = (new Task())->save($data);

        flock($fp,LOCK_UN);

        fclose($fp);

        $this->writeOperatorLog('新增任务:'.$data['task_name']);

        return success($result);
    }
}