<?php
/**
 * Created by PhpStorm.
 * User: 陈朔
 * Date: 2018/11/9
 * Time: 9:54 AM
 */

namespace App\Http\Controllers\App\V1\Task;

use App\Common\Constant;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use App\Models\TaskMemberRecord;

class TaskBaseController extends BaseController
{
    /**
     * 新增或更新用户邀请记录日志
     * @param int $taskId 任务ID
     * @param int $customerId 外部人员或会员 ID
     * @param int $taskType 任务类型 2:客户管理，3:生日祝福，4:礼券到期提醒
     * @param int $displayState 是否在列表显示 1=显示 2=隐藏
     * @return bool|void
     */
    protected function saveLog($taskId = 0, $customerId = 0, $taskType = 0, $displayState = 0)
    {
        $data = [
            'enterprise_id' => $this->enterpriseId,
            'mem_uid' => $this->memUid,
            'task_id' => $taskId,
            'customer_id' => $customerId,
            'task_type' => $taskType,
            'display_state' => TaskMemberRecord::TASK_SHOW_STATUS,
            'created_at' => time()
        ];

        // 获取是否邀请过当前指定用户
        $total = TaskMemberRecord::where('enterprise_id', $this->enterpriseId)
            ->where('mem_uid', $this->memUid)
            ->where('customer_id', $customerId)->where('task_type', $taskType)
            ->where('task_id', $taskId)->count();

        $data['display_state'] = $displayState;

        // 如果是隐藏状态
        if (TaskMemberRecord::TASK_HIDE_STATUS == $displayState) {

            // 如果邀请过用户则更新
            if (!empty($total)) {
                unset($data['created_at']);
                $data['updated_at'] = time();
                try {
                    TaskMemberRecord::where('enterprise_id', $this->enterpriseId)
                        ->where('mem_uid', $this->memUid)->where('customer_id', $customerId)
                        ->where('task_type', $taskType)->where('task_id', $taskId)
                        ->update($data);
                } catch (\Exception $e) {
                    error(200009);
                }

                return true;
            }
        }

        // 如果没有邀请过当前指定用户则新增
        if (empty($total)) {

            try {
                TaskMemberRecord::insert($data);
            } catch (\Exception $e) {
                error(200009);
            }
        }

        return true;
    }

    /**
     * 获取邀请状态以及用户信息
     * @param array $data 邀请数据
     * @param int $taskId 任务ID
     * @param int $taskType 任务类型
     * @param bool $isShowName 是否显示备注名
     * @return array
     */
    protected function getUserData($data = [], $taskId = 0, $taskType = 0, $isShowName = false)
    {
        $customerDataIds = [];
        $unionList       = [];
        // 获取用户UID
        $customerIds = array_column($data, 'ec_userid');

        if (empty($customerIds)) {
            $customerIds = array_column($data, 'union_id');
            $unionList   = MemberUserRelation::whereIn('union_id', $customerIds)
                ->where('enterprise_id', $this->enterpriseId)
                ->where('mem_uid', $this->memUid)
                ->get(['union_id', 'ec_userid', 'user_remark_name'])->toArray();
            if (!empty($unionList)) {
                $unionList   = array_column($unionList, null, 'union_id');
                $customerIds = array_column($unionList, 'ec_userid');
            }

        }

        // 获取是否邀请过用户
        $list = TaskMemberRecord::where('enterprise_id', $this->enterpriseId)
            ->where('mem_uid', $this->memUid)->whereIn('customer_id', $customerIds)
            ->where('task_type', $taskType)->where('task_id', $taskId)
            ->get(['customer_id'])->toArray();

        // 已邀请人员IDS
        if (!empty($list)) {
            $customerDataIds = array_column($list, 'customer_id');
        }

        // 遍历数据
        foreach ($data as &$v) {
            // 如果需要显示客户名则客户手机号不存在则不是会员
            if (!empty($isShowName) && empty($v['mobile'] ?? '')
                && !empty($v['union_id'] ?? '')
                && !empty($unionList[$v['union_id']]['user_remark_name'] ?? '')
            ) {
                $v['user_name'] = $unionList[$v['union_id']]['user_remark_name'] ?? '';
            }

            // 如果用户UCID为空切存在union_id
            if (empty($v['ec_userid']) && !empty($v['union_id'] ?? '')) {
                $v['ec_userid'] = $unionList[$v['union_id']]['ec_userid'] ?? '';
            }

            $v['status'] = in_array($v['ec_userid'], $customerDataIds) ? Constant::TASK_REMIND_STATUS : Constant::TASK_Not_REMIND_STATUS;
        }

        return $data;
    }
}

