<?php

namespace App\Http\Controllers\App\V1\Task\Outsiders;

use App\Common\Constant;
use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\MemberUserRelation;
use App\Models\Task;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;

class ListController extends TaskBaseController
{
    /**
     * 非会员列表
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 分页参数
        $page   = $request->get('page', Constant::DEFAULT_PAGE);
        $size   = $request->get('size', Constant::DEFAULT_SIZE);
        $taskId = $request->get('task_id', '');

        // 判断任务ID
        if (empty($taskId)) {
            error(200010);
        }

        $userList      = [];
        $differenceIds = [];

        // 获取客户unionIds集合
        $unionList = MemberUserRelation::where('mem_uid', $this->memUid)
            ->where('enterprise_id', $this->enterpriseId)
            ->where('union_id', '<>', '')
            ->get(['union_id'])->toArray();

        if (!empty($unionList)) {
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'union_ids' => array_column($unionList, 'union_id')
            ];

            // 获取会员用户IDS
            $userList = (new User())->userUnionId($params);
        }

        // 如果有已返回会员用户
        if (!empty($userList['data'] ?? [])) {
            $differenceIds = array_column($userList['data'], 'union_id');
        }

        // 统计非会员数
        $total = MemberUserRelation::where('mem_uid', $this->memUid)
            ->where('enterprise_id', $this->enterpriseId)
            ->when(!empty($differenceIds), function ($query) use ($differenceIds) {
                return $query->WhereNotIn('union_id', $differenceIds);
            })->count();

        // 如果存在非会员
        if ($total) {
            // 查询非会员列表
            $list = MemberUserRelation::where('mem_uid', $this->memUid)
                ->where('enterprise_id', $this->enterpriseId)
                ->when(!empty($differenceIds), function ($query) use ($differenceIds) {
                    return $query->WhereNotIn('union_id', $differenceIds);
                })->forpage($page, $size)
                ->get(['ec_userid', 'ec_name', 'ec_avatar', 'user_remark_name'])->toArray();
        }

        // 遍历数据
        foreach ($list ?? [] as $v) {

            // 默认名称
            $userName = $v['ec_name'];
            if (!empty($v['user_remark_name'])) {
                $userName = $v['user_remark_name'];
            }

            $data[] = [
                'ec_userid' => $v['ec_userid'] ?? '',
                'user_name' => $userName ?? '',
                'avatar' => $v['ec_avatar'] ?? '',
            ];
        }

        // 如果存在数据
        if (!empty($data)) {
            $data = $this->getUserData($data, $taskId, Task::CUSTOMER_MANAGEMENT);
        }

        return success([
            'data' => $data ?? [],
            'meta' => [
                'total' => $total ?? 0,
                'page' => $page ?? 0,
                'size' => $size ?? 0
            ]
        ]);
    }
}

