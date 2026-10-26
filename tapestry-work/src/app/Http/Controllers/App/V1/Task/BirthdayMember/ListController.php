<?php

namespace App\Http\Controllers\App\V1\Task\BirthdayMember;

use App\Common\Constant;
use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Models\MemberUserRelation;
use App\Models\Task;
use Illuminate\Http\Request;
use CrsSdk\Shop\User;

class ListController extends TaskBaseController
{
    /**
     * 今日过生日人员列表
     *
     * @param Request $request
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

        // 查询当前导购所有小程序进入的客户列表
        $memberRelationList = MemberUserRelation::where('mem_uid', $this->memUid)
            ->where('enterprise_id', $this->enterpriseId)
            ->where('union_id', '<>', '')->get(['union_id'])
            ->toArray();

        $defaultData = [
            'data' => [],
            'meta' => [
                'total' => 0,
                'page' => $page ?? 0,
                'size' => $size ?? 0
            ]
        ];

        // 如果存在导购
        if (!empty($memberRelationList)) {
            $paramsData = [
                'page' => $page,
                'size' => $size,
                'enterprise_id' => $this->enterpriseId,
                'union_ids' => array_column($memberRelationList, 'union_id')
            ];

            $data = (new User())->birthdayList($paramsData);

            // 如果有返回数据
            if (!empty($data['data'])) {
                $data['data'] = $this->getUserData($data['data'], $taskId, Task::BIRTHDAY_WISHES);
            }
            return success($data);

        } else {
            return success($defaultData);
        }
    }
}

