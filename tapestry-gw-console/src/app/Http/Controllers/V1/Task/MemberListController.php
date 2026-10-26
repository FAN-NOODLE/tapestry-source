<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 15:06
 */

namespace App\Http\Controllers\V1\Task;

use App\Http\Controllers\V1\BaseController;
use App\Services\TaskService;
use CrsSdk\Work\Member;
use CrsSdk\Work\Role;
use Illuminate\Http\Request;
use CrsSdk\Work\Task;
use Illuminate\Support\Facades\Log;

class MemberListController extends BaseController
{
    /**
     * 员工在职状态：在职
     */
    const MEM_ACTIVE_TRUE = 1;

    /**
     * 页码最大条数
     */
    const SIZE_MAX = 9999;

    /**
     * 获取员工接口
     *
     * @param Request $request
     * @return \App\Http\Resources\Success|array|mixed|null
     */
    public function index(Request $request, TaskService $taskService)
    {
        // 初始化
        $memberIds = [];
        // 门店id，手机号，姓名
        $page = $request->get('page', 1);
        $size = $request->get('size', 0);
        $search = $request->get('search', []);
        $memUid = $request->get('mem_uid', []);

        // 不分页，size重新赋值
        if (!$size) {
            $size = 1;
        }

        // 员工参数列表
        $params = [
            'pageNum' => $page,
            'pageSize' => $size
        ];
        // 员工id存在
        if (empty($memUid)) {
            // 任务条件：类型，开始时间，周期，是否循环
            $data['search']['task_type'] = $search['task_type'];
            $data['enterprise_id'] = $this->enterpriseId;
            // 时间戳毫秒转换成秒
            $data['search']['start_time'] = $search['start_time'] / 1000;
            if (!empty($search['end_time'])) {
                $data['search']['end_time'] = (int)($search['end_time'] / 1000);
            }
            // 开始时间时间格式
            $startTime = date('Y-m-d H:i:s', $data['search']['start_time']);
            // 任务不循环
            if ($search['is_loop'] == \App\Common\Task::NO_IS_LOOP) {
                // 结束时间为前一天的59秒
                $data['search']['end_time'] = strtotime($taskService->getEndTime($startTime, $search['task_cycle'])) - 1;
            }
            // 在当前时间区间内，改类型任务下的员工uid
            $data['size'] = self::SIZE_MAX;

            $taskList = (new Task())->taskList($data);
            $taskId = array_unique(array_column($taskList['data'], 'id'));
            $member = (new Task())->taskMember(['ids' => $taskId]);
            $memberIds = array_unique(array_column($member['data'], 'mem_uid'));
        }

        // 姓名
        if (isset($search['mem_username'])) {
            $params['memUsername'] = $search['mem_username'];
        }

        // 手机号码
        if (isset($search['mem_mobile'])) {
            $params['memMobile'] = $search['mem_mobile'];
        }

        // 门店id
        if (isset($search['dp_id']) && $search['dp_id']) {
            $params['dpId'] = $search['dp_id'];
        }

        // 员工uid集合
        if (!empty($memUid)) {
            $params['memUids'] = $memUid;
        }

        // 员工角色
        $roleCondition = [
            'qy_domain' => $this->qyDomain,
            'params' => []
        ];
        $roleList = (new Role())->list($roleCondition);
        $roleId = [];
        // 获取员工角色为店长和导购的员工
        foreach ($roleList['list'] as $role) {
            if (in_array($role['roleName'], \App\Common\Task::TASK_ROLE_NAME)) {
                array_push($roleId, $role['roleId']);
            }
        }
        $params['roleIdList'] = $roleId;
        // 调用sdk
        $condition = [
            'qy_domain' => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params' => $params
        ];
        Log::debug('员工列表请求uc参数: ' . json_encode($condition));
        $list = (new Member())->memberList($condition);
        Log::debug('员工列表uc返回参数: ' . json_encode($list));

        // 总数
        $total = intval($list['total']);
        // 不分页 && 员工多余一个，查询列表
        if (1 == $size && $total > 1) {
            // size重新赋值
            $condition['params']['pageSize'] = $total;
            // 列表重新赋值
            $list = (new Member())->memberList($condition);
        }

        // 用户列表数据格式化
        $data = $this->format($list['list'], $memberIds);

        // 有分页
        if ($size != 1) {
            return success([
                'data' => $data,
                'meta' => [
                    'total' => intval($list['total']),
                    'page' => intval($page),
                    'size' => intval($size)
                ]
            ]);
        }

        return success(['data' => $data]);
    }

    /**
     * 格式化员工数据
     *
     * @param array $list 当前员工
     * @param array $memUid 已有任务员工
     * @return array
     */
    protected function format($list, $memUid = [])
    {
        // 初始化待返数据
        $data = [];
        if (!$list) {
            return $data;
        }

        foreach ($list as $member) {
            $data[] = [
                'mem_uid' => $member['memUid'],
                'mem_username' => $member['memUsername'],
                'mem_face' => $member['memFace'],
                'dp_id' => $member['dpName'][0]['dpId'] ?? '',
                'dp_name' => $member['dpName'][0]['dpName'] ?? '',
                'mem_role' => $member['memRole'],
                'mem_gender' => $member['memGender'],
                'mem_mobile' => $member['memMobile'],
                'is_user' => in_array($member['memUid'], $memUid) ? 1 : 2,
                'mem_active' => $member['memActive']
            ];
        }

        return $data;
    }
}