<?php

namespace App\Http\Controllers\App\V1\Task\Memorandum;

use App\Http\Controllers\App\V1\Task\TaskBaseController;
use App\Library\Uc\External;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\User;
use Illuminate\Http\Request;
use App\Services\TaskService;

class ListController extends TaskBaseController
{
    /**
     * 指定客户任务备忘录列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */

    public function index(Request $request)
    {
        $ecUserId = $request->get('ec_userid', '');

        if (empty($ecUserId)) {
            error(200012);
        }

        // 获取客户详情
        $info = $this->getUserInfo($ecUserId);

        // 初始化taskService
        $taskService = new TaskService();

        // 组装数据
        $data = [
            'data' => $taskService->taskList($this->enterpriseId, $this->memUid, $ecUserId, $info['birthday'], $info['union_id']),
            'user_name' => $info['user_name'] ?? '',
            'avatar' => $info['avatar'] ?? '',
            'birthday' => $info['birthday'] ?? '',
            'mobile' => $info['mobile'] ?? '',
        ];

        return success($data);
    }

    /**
     * 获取客户详情
     * @param string $ecUserId 当前用户ID
     * @return array
     * @throws \App\Exceptions\BaseException
     */
    protected function getUserInfo($ecUserId = '')
    {
        // 用户关系数据
        $userRelationData = MemberUserRelation::where('ec_userid', $ecUserId)->where('mem_uid',$this->memUid)
            ->where('enterprise_id', $this->enterpriseId)->first([
                    'ec_name',
                    'ec_avatar',
                    'user_remark_name',
                    'union_id',
                    'ec_id',
                    'user_remark_name'
                ]
            );

        $data = [];

        $data['union_id'] = $userRelationData->union_id ?? '';

        // 如果是会员则查询会员详情
        if (!empty($data['union_id'])) {
            $params = [
                'enterprise_id' => $this->enterpriseId,
                'union_id' => $data['union_id']
            ];

            // 用户详情
            $userInfo = (new User())->detail($params);
    
            $data['user_name'] = $userRelationData->user_remark_name ?? ($userRelationData->ec_name ?? '');
            $data['avatar']    = $userRelationData->ec_avatar ?? '';
            $data['birthday']  = '';
            $data['mobile']    = '';
            // 如果是会员
            if (!empty($userInfo['data'])) {
                $data['user_name'] = $userRelationData->user_remark_name ?? ($userInfo['data']['real_name'] ?? $userRelationData->ec_name);
                $data['avatar']    = $userInfo['data']['avatar'] ?? '';
                $data['birthday']  = $userInfo['data']['birthday'] ?? '';
                $data['mobile']    = $userInfo['data']['mobile'] ?? '';
            }

        } else {
            // 如果是外部联系人获取外部联系人
            $condition = [
                'ecId' => $userRelationData->ec_id,
                'app' => $this->appIdentifier
            ];

            // 获取外部人员数据
            $userInfo = (new External())->detail($this->qyDomain, $condition);

            $data['user_name'] = $userInfo['ecName'] ?? ($userInfo['memUsername'] ?? '');
            $data['avatar']    = $userInfo['ecAvatar'] ?? '';
            $data['birthday']  = $userInfo['birthday'] ?? '';
            $data['mobile']    = $userInfo['ecMobile'] ?? '';
        }

        // 如果没有手机号切有备注则显示
        if (empty($data['mobile']) && !empty($userRelationData->user_remark_name)) {
            $data['user_name'] = $userRelationData->user_remark_name ?? '';
        }

        return $data;
    }
}


