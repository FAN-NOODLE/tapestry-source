<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/3/25 16:39
 */

namespace App\Http\Controllers\Service\V1\Task;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\TaskMemberRelation;
use App\Services\TaskService;
use Illuminate\Http\Request;
use App\Models\Task;
use App\Models\TaskMemberRecord;

class TargetController extends BaseController
{

    const STATUS_COMPLETE = 2;

    /**
     * 指标达成率(任务)
     *
     * @param Request $request
     * @return array|void
     */
    public function index(Request $request)
    {
        $input = $request->input();

        if (empty($input['enterprise_id'])) {
            error(600000);
        }

        try {
            // 分门店查询
            if (isset($input['stores'])) {
                $result = Task::query()->join('task_member_relation', 'task_member_relation.task_id', '=', 'task.id')
                    ->where('task.enterprise_id', '=', $input['enterprise_id'])
                    ->when($input['type'] === 1, function ($query) {
                        // 自然流量
                        $query->where('task_member_relation.store_id', '=', '');
                    })->when($input['type'] === 2, function ($query) use ($input) {
                        // 门店
                        $query->whereIn('task_member_relation.store_id', $input['stores']);
                    })->when($input['type'] === 3, function ($query) use ($input) {
                        // 自然流量+门店
                        $query->whereIn('task_member_relation.store_id', $input['stores'])->orWhere('task_member_relation.store_id', '=', '');
                    })->get(['task.*', 'task_member_relation.mem_uid']);
            } else {
                $result = Task::query()->join('task_member_relation', 'task_member_relation.task_id', '=', 'task.id')
                    ->where('task.enterprise_id', '=', $input['enterprise_id'])
                    ->get(['task.*', 'task_member_relation.mem_uid']);
            }

            // 初始化数据结构
            $data = [
                'sales' => ['total' => 0, 'complete_number' => 0],
                'customer' => ['total' => 0, 'complete_number' => 0],
                'birthday' => ['total' => 0, 'complete_number' => 0],
                'coupon' => ['total' => 0, 'complete_number' => 0],
                'member' => ['total' => 0, 'complete_number' => 0]
            ];

            if (!empty($result)) {
                $server = new TaskService();
                foreach ($result as $key => $value) {

                    $info = $server->getTaskData($value, $value['mem_uid']);

                    switch ($value['task_type']) {
                        case Task::SALES_PERFORMANCE:
                            $data['sales']['total']++;
                            if (isset($info['status']) && $info['status'] === self::STATUS_COMPLETE) {
                                $data['sales']['complete_number']++;
                            }
                            break;
                        case Task::CUSTOMER_MANAGEMENT:
                            $data['customer']['total']++;
                            if (isset($info['status']) && $info['status'] === self::STATUS_COMPLETE) {
                                $data['customer']['complete_number']++;
                            }
                            break;
                        case Task::BIRTHDAY_WISHES:
                            $data['birthday']['total']++;
                            if (isset($info['status']) && $info['status'] === self::STATUS_COMPLETE) {
                                $data['birthday']['complete_number']++;
                            }
                            break;
                        case Task::COUPONS_EXPIRATION:
                            $data['coupon']['total']++;
                            if (isset($info['status']) && $info['status'] === self::STATUS_COMPLETE) {
                                $data['coupon']['complete_number']++;
                            }
                            break;
                        case Task::MEMBER_RECRUITMENT:
                            $data['member']['total']++;
                            if (isset($info['status']) && $info['status'] === self::STATUS_COMPLETE) {
                                $data['member']['complete_number']++;
                            }
                            break;
                    }

                }
            }
        } catch (\Exception $e) {
            error(110000);
        }

        return success(['data' => $data]);
    }

}