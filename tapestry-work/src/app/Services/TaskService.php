<?php

namespace App\Services;

use App\Common\Constant;
use App\Models\MemberUserRelation;
use App\Models\Task;
use App\Models\TaskMemberRecord;
use App\Models\TaskMemberRelation;
use Carbon\Carbon;
use CrsSdk\Shop\CouponUser;
use CrsSdk\Shop\Member;
use CrsSdk\Shop\User;
use Illuminate\Support\Facades\DB;

class TaskService
{
    /**
     * 任务表
     * @var Task
     */
    protected $task;
    /**
     * 任务员工关系表
     * @var TaskMemberRelation
     */
    protected $taskMemberRelation;

    public function __construct()
    {
        $this->task = new Task();
        $this->taskMemberRelation = new TaskMemberRelation();
    }

    /**
     * 获取本次任务开始时间以及结束时间
     * @param int $startTime 开始时间
     * @param int $taskType 任务周期：(1:一天，2:一周，3:一个月，4:三个月，5:六个月，6:一年)
     * @param int $endTime 任务结束时间
     * @return array
     */
    public function getTaskTimeSection($startTime = 0, $taskType = 0, $endTime = 0)
    {

        $dateStr = date('Y-m-d', time());
        $todayStartTime = strtotime($dateStr);
        // 如果循环的是一天
        if (Task::ONE_DAY_TYPE == $taskType) {

            // 开始时间
            $taskStartTime = $todayStartTime;
            // 结束时间
            $taskEndTime = $todayStartTime + 86400;
            if ($taskEndTime > $endTime) {
                $taskEndTime = $endTime;
            }
        } else {
            // 判断间间隔几周
            $data = $this->getTime($startTime, $taskType, $endTime);
            // 开始时间
            $taskStartTime = $data['start_time'];
            $taskEndTime = $data['end_time'];
        }

        return ['start_time' => $taskStartTime, 'end_time' => $taskEndTime];
    }

    /**
     * 获取当前的时间区间
     * @param int $taskStartTime 开始时间
     * @param int $taskType 任务周期：(2:一周，3:一个月，4:三个月，5:六个月，6:一年)
     * @param int $taskEndTime 结束时间
     * @return array
     */
    private function getTime($taskStartTime = 0, $taskType = 0, $taskEndTime = 0)
    {
        // 剩余时间戳
        $surplusTime = $taskEndTime - $taskStartTime;
        if (empty($surplusTime)) {
            return [];
        }

        $second = 0;
        switch ($taskType) {
            // 一周
            case Task::ONE_WEEK_TYPE:
                $date = "1week";
                $second = 604800;
                break;
            // 一个月
            case Task::ONE_MONTHS_TYPE:
                $date = "1month";
                $second = 86400 * 28;
                break;
            // 三个月
            case Task::THREE_MONTHS_TYPE:
                $date = "3month";
                $second = 86400 * 28 * 3;
                break;
            // 六个月
            case Task::SIX_MONTHS_TYPE:
                $date = "6month";
                $second = 86400 * 28 * 6;
                break;
            // 一年
            case Task::ONE_YEAR_TYPE:
                $date = "1year";
                $second = 31536000;
                break;
            default:
                break;
        }


        // 如果类型不存在
        if (empty($second)) {
            return ['start_time' => 0, 'end_time' => 0];
        }

        $loopWeek = ceil($surplusTime / $second);
        // 判断区间
        for ($i = 0; $i < $loopWeek; $i++) {

            $data_list = date('Y-m-d', $taskStartTime);
            $data_list = explode('-', $data_list);
            $CarbonDate = Carbon::createFromDate($data_list[0], $data_list[1], $data_list[2]);

            // 处理添加几个月
            if (in_array($taskType, [Task::ONE_MONTHS_TYPE, Task::THREE_MONTHS_TYPE, Task::SIX_MONTHS_TYPE])) {
                // 转换时间
                $MonthDate = $CarbonDate->addMonthNoOverflow(intval($date));
            }

            // 如果是加一周
            if (Task::ONE_WEEK_TYPE == $taskType) {
                $MonthDate = $CarbonDate->addWeek(intval($date));
            }

            // 如果是加一年
            if (Task::ONE_YEAR_TYPE == $taskType) {
                $MonthDate = $CarbonDate->addYears(intval($date));
            }

            $endTime = strtotime($MonthDate->toDateTimeString());

            // 如果结束时间小于等于周期时间
            if ($taskEndTime <= $endTime) {
                $endTime = $taskEndTime;
                break;
            }

            // 如果当前时间在开始和结束时间之间
            if (time() < $endTime) {
                continue;
            }

            $taskStartTime = $endTime;
        }

        return ['start_time' => $taskStartTime, 'end_time' => $endTime];
    }

    /**
     * 获取指定导购任务列表
     * @param array $data 单条任务表数据
     * @param string $memUid 导购ID
     * @return array
     */
    public function getTaskData($data = [], $memUid = '')
    {
        if (empty($memUid)) {
            return [];
        }

        $list = [];
        $taskCount = 0;
        $data['status'] = Constant::TASK_EXECUTING_STATUS;

        $time = time();

        // 当前时间大于等于开始时间小于登录结束时间
        if ($time >= $data['start_time'] && $time <= $data['end_time']) {

            $taskStartTime = $data['start_time'];
            $taskEndTime = $data['end_time'];
            // 如果是循环
            if (Task::TASK_LOOP_STATUS == $data['is_loop']) {

                $timeData = $this->getTaskTimeSection($taskStartTime, $data['task_type'], $taskEndTime);
                $taskStartTime = $timeData['start_time'];
                $taskEndTime = $timeData['end_time'] - 1;
            }

            // 销售业绩类型
            if (Task::SALES_PERFORMANCE == $data['task_type']) {
                $countParams = [
                    'enterprise_id' => $data['enterprise_id'],
                    'mem_uid' => $memUid,
                    'start_time' => $taskStartTime,
                    'end_time' => $taskEndTime,
                ];
                // 获取统计数据
                $totalData = (new Member())->timeSlotSale($countParams);

                // 销售额
                $taskCount = $totalData['data']['money'] ?? 0;
                if (!empty($taskCount)) {
                    // 金额分转换成元
                    $taskCount = $taskCount / 100;
                }

                // 分转换元
                $data['task_index'] = $data['task_index'] / 100;
            }

            // 如果是会员朝暮
            if (Task::MEMBER_RECRUITMENT == $data['task_type']) {

                // 调用商城会员sdk获取会员数据
                $result = (new User())->userCountByMemUids([
                    'enterprise_id' => $data['enterprise_id'],
                    'mem_uids' => [$memUid],
                    'start_time' => $taskStartTime,
                    'end_time' => $taskEndTime
                ]);
                $taskCount = $result['data'][0]['total'] ?? 0;
            }

            // 如果是礼劵到期类型
            if (Task::COUPONS_EXPIRATION == $data['task_type']) {
                $couponParams = [
                    'enterprise_id' => $data['enterprise_id'],
                    'mem_uid' => $memUid,
                    'start_time' => $taskStartTime,
                    'end_time' => $taskEndTime
                ];
                // 获取统计数据
                $couponData = (new CouponUser())->couponCount($couponParams);

                $taskCount = $couponData['data']['total'] ?? 0;
            }

            // 如果是客户管理类型
            if (Task::CUSTOMER_MANAGEMENT == $data['task_type']) {
                $taskCount = MemberUserRelation::where('mem_uid', $memUid)->where('enterprise_id', $data['enterprise_id'])
                    ->where('created_at', '>=', $taskStartTime)->where('created_at', '<=', $taskEndTime)->count();
            }

            // 如果是生日祝福
            if (Task::BIRTHDAY_WISHES == $data['task_type']) {

                // 查询当前导购所有小程序进入的客户列表
                $memberRelationList = MemberUserRelation::where('mem_uid', $memUid)->where('enterprise_id', $data['enterprise_id'])
                    ->where('union_id', '<>', '')->get(['union_id'])->toArray();
                if (!empty($memberRelationList)) {
                    $paramsData = [
                        'page' => 1, // 查询统计默认页码
                        'size' => 1, // 查询统计默认每页数
                        'union_ids' => array_column($memberRelationList, 'union_id'),
                        'enterprise_id' => $data['enterprise_id'],
                    ];
                    $RelationData = (new User())->birthdayList($paramsData);

                    // 获取今日过生日人员数
                    $total = $RelationData['meta']['total'] ?? 0;

                    // 如果存在今日会员
                    if (!empty($total)) {
                        // 当前已完成任务数
                        $taskCount = TaskMemberRecord::where('task_type', $data['task_type'])->where('mem_uid', $memUid)
                            ->where('created_at', '>=', $taskStartTime)->where('enterprise_id', $data['enterprise_id'])
                            ->where('created_at', '<=', $taskEndTime)->count();

                        // 如果今日过生日人员数和已送祝福人员数字一致
                        if ($total == $taskCount) {
                            $data['status'] = Constant::TASK_EXE_SUCCESS_STATUS;
                        }
                    }
                }

            }

            // 状态更改为以完成
            if ($data['task_index'] <= $taskCount && Task::BIRTHDAY_WISHES != $data['task_type']) {
                $data['status'] = Constant::TASK_EXE_SUCCESS_STATUS;
            }

            // 百分比初始化
            $percentage = 0;
            if (!empty($data['task_index']) && !empty($taskCount)) {
                $percentage = floor(round($taskCount / $data['task_index'] * 10000, 2) / 100);
            }

            $list = [
                'id' => $data['id'], // 任务ID
                'enterprise_id' => $data['enterprise_id'], // 任务对应企业ID
                'task_name' => $data['task_name'], // 任务名称
                'task_type' => $data['task_type'], // 任务类型
                'start_time' => $taskStartTime, // 开始时间
                'end_time' => $taskEndTime, // 结束时间
                'task_index' => $data['task_index'], // 任务邀请完成数
                'task_count' => $taskCount, // 当前导购已完成数
                'percentage' => $percentage,
                'status' => $data['status'] // 是否完成任务 1=未完成 2=已完成
            ];
        }

        return $list;
    }

    /**
     * 获取任务列表
     * @param string $enterpriseId 当前企业ID
     * @param string $memUid 当前导购ID
     * @param int $ecUserId 客户ID
     * @param string $birthday 客户生日
     * @param string $union_id 关联会员ID
     * @return array
     */
    public function taskList($enterpriseId = '', $memUid = '', $ecUserId = 0, $birthday = '', $union_id = '')
    {

        // 查询用户所在任务的数据
        $taskMemberRelationData = TaskMemberRelation::where('mem_uid', $memUid)
            ->where('enterprise_id', $enterpriseId)->get(['task_id'])->toArray();

        $data = [];
        // 如果当前用户存在任务
        if (!empty($taskMemberRelationData)) {
            // 获取任务ID集合
            $taskIds = array_column($taskMemberRelationData, 'task_id');

            // 获取我可以参与的任务数据
            $taskData = Task::whereIn('id', $taskIds)->where('task_type', '<>', Task::CUSTOMER_MANAGEMENT)
                ->where('enterprise_id', $enterpriseId)->get()->toArray();

            $time = time();
            // 遍历数据排除不是今天需要完成的任务
            foreach ($taskData as $v) {

                // 当前时间大于等于开始时间小于登录结束时间
                if ($time >= $v['start_time'] && $time <= $v['end_time']) {

                    // 如果是不显示的则剔除掉
                    $total = TaskMemberRecord::where('display_state', TaskMemberRecord::TASK_HIDE_STATUS)
                        ->where('mem_uid', $memUid)->where('customer_id', $ecUserId)
                        ->where('enterprise_id', $enterpriseId)
                        ->where('task_type', $v['task_type'])->where('task_id', $v['id'])->count();

                    // 如果存在就不显示待办任务
                    if (!empty($total)) {
                        continue;
                    }

                    $taskStartTime = $v['start_time'];
                    $taskEndTime = $v['end_time'];

                    // 如果是循环
                    if (Task::TASK_LOOP_STATUS == $v['is_loop']) {
                        $timeData = $this->getTaskTimeSection($taskStartTime, $v['task_type'], $taskEndTime);
                        $taskStartTime = $timeData['start_time'];
                        $taskEndTime = $timeData['end_time'] - 1;
                    }

                    // 如果是招募会员
                    if (Task::MEMBER_RECRUITMENT == $v['task_type']) {
                        if (!empty($birthday)) {
                            continue;
                        }
                    }

                    // 如果是礼劵到期类型
                    if (Task::COUPONS_EXPIRATION == $v['task_type']) {

                        if (empty($union_id)) {
                            continue;
                        }

                        $couponParams = [
                            'enterprise_id' => $enterpriseId,
                            'mem_uid' => $memUid,
                            'union_id' => $union_id,
                            'start_time' => $taskStartTime,
                            'end_time' => $taskEndTime
                        ];

                        // 获取统计数据
                        $couponData = (new CouponUser())->couponCount($couponParams);

                        $taskCount = $couponData['data']['total'] ?? 0;
                        if (empty($taskCount)) {
                            continue;
                        }
                    }

                    // 如果是生日祝福
                    if (Task::BIRTHDAY_WISHES == $v['task_type']) {
                        if (!empty($birthday)) {

                            // 如果今日不过生日则不显示
                            if (date('m-d', strtotime($birthday)) !== date('m-d', time())) {
                                continue;
                            }
                        } else {
                            continue;
                        }
                    }

                    $data[] = [
                        'id' => $v['id'],
                        'enterprise_id' => $v['enterprise_id'],
                        'task_name' => $v['task_name'],
                        'task_type' => $v['task_type'],
                        'start_time' => $taskStartTime,
                        'ec_userid' => $ecUserId,
                        'end_time' => $taskEndTime,
                    ];
                }
            }
        }

        return $data;
    }

    /**
     * 根据数据权限及条件 获取任务列表
     * @param $page
     * @param $size
     * @param $search
     * @param $enterpriseId
     * @param $storeRole
     * @return array
     */
    public function getTaskListByPermission($page, $size, $search, $enterpriseId, $storeRole)
    {
        $filed = ['task.id', 'task_name', 'start_time', 'end_time', 'task_cycle', 'task_type', 'serve_type'];
        // 查询条件拼凑
        $where = [];
        if (isset($search['task_name'])) {
            array_push($where, ['task_name', 'LIKE', '%' . $search['task_name'] . '%']);
        }
        if (isset($search['task_type']) && is_numeric($search['task_type'])) {
            array_push($where, ['task_type', '=', $search['task_type']]);
        }
        // 服务类型:crs/shop下存在
        if (isset($search['serve_type']) && in_array($search['serve_type'], [1, 2])) {
            array_push($where, ['serve_type', '=', $search['serve_type']]);
        }

        $orm = Task::query()
            ->whereHas('taskMemberRelations', function ($query) use ($enterpriseId, $storeRole) {
                $query->{'when'}(isset($storeRole['is_super_admin'])
                    && ($storeRole['is_super_admin'] != 1), function ($query) use ($storeRole) {
                            $query->{'when'}(isset($storeRole['type'])
                                && $storeRole['type'] == 1, function ($query) {
                                    $query->{'where'}('store_id', '=', '');
                                })
                            ->{'when'}(isset($storeRole['type'])
                                && $storeRole['type'] == 2
                                && !empty($storeRole['stores']), function ($query) use ($storeRole) {
                                    $query->{'whereIn'}('store_id', $storeRole['stores']);
                                })
                            ->{'when'}(isset($storeRole['type'])
                                && $storeRole['type'] == 3
                                && !empty($storeRole['stores']), function ($query) use ($storeRole) {
                                    $query->{'whereIn'}('store_id', $storeRole['stores'])
                                          ->{'orWhere'}('store_id', '=', '');
                                });
                    })
                    ->{'where'}('enterprise_id', '=', $enterpriseId);
            })
            ->where($where)
            ->where('enterprise_id', '=', $enterpriseId)
            ->when(isset($search['start_time']) && isset($search['end_time'])
                && !empty($search['start_time']) && !empty($search['end_time']), function ($query) use ($search) {
                    $query->{'where'}(function ($query) use ($search) {
                        $query->{'where'}([
                                ['start_time', '>', $search['start_time']],
                                ['start_time', '<', $search['end_time']]
                            ])
                        ->orwhere([
                                ['start_time', '<=', $search['start_time']],
                                ['end_time', '>=', $search['end_time']]
                            ])
                        ->orwhere([
                                ['end_time', '>', $search['start_time']],
                                ['end_time', '<', $search['end_time']]
                            ]);
                    });
                })
            ->orderBy('created_at', 'desc')
            ->paginate($size, $filed, '', $page);
        return [$orm->total() ?? 0, $orm->items() ?? []];
    }
}
