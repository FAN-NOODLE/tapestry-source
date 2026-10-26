<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/12/24
 * Time: 11:34
 */

namespace App\Listeners;

use App\Events\AddTaskEvent;
use App\Exceptions\Error;
use CrsSdk\Work\Task;
use Illuminate\Contracts\Queue\ShouldQueue;

class AddTaskListener implements ShouldQueue
{
    /**
     * 页码最大条数
     */
    const SIZE_MAX = 9999;

    /**
     * 当前数组
     */
    public $data;

    /**
     * Create the event listener.
     *
     * @return void
     */
    public function __construct()
    {
        //
    }

    /**
     * Handle the event.
     *
     * @param  AddTaskEvent $event
     * @return array
     */
    public function handle(AddTaskEvent $event)
    {
        $data = $event->getData();
        $this->data = $data;
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
            throw new Error(600004, '保存失败，所选导购在该时期内已有同类型任务。', 500, $repeatMember);
        }

        $result = (new Task())->save($data);

        return $result;
    }
}