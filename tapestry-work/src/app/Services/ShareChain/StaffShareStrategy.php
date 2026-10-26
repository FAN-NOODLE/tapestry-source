<?php

namespace App\Services\ShareChain;

use App\Common\Constant;
use App\Common\MessageConstant;
use App\Jobs\Message\ShareLinkMessage;
use App\Jobs\Message\ShareLinkMessageAdapter;
use App\Jobs\ShareLinkBatchUpdateCustomerStatisticJob;
use App\Jobs\ShareLinkUpdateCustomerStatisticJob;
use App\Library\Ddc\DdcMessage;
use App\Models\StaffShare;
use App\Models\StoreWorker;
use App\Services\Events\ShareLinkChangeEventService;
use App\Services\StaffManagerService;
use Illuminate\Support\Facades\Log;

class StaffShareStrategy implements StrategyInterface
{
    /**
     * binding_type
     */
    const SETTING_BINDING_TYPE_FOREVER = 1;
    const SETTING_BINDING_TYPE_EXPIRE = 2;

    /**
     * change_owner
     */
    const SETTING_CHANGE_OWNER_NOT_ALLOW = 1;
    const SETTING_CHANGE_OWNER_ALLOW = 2;

    /**
     * 员工进入初始化数据
     * 此时user_id 和 memUid 和worker_id 指同一个人
     *
     * @param $enterpriseId
     * @param $userId
     * @param $memUid
     * @param $workerId
     * @param $type
     *
     * @return StaffShare|array|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     */
    public function staffInitial($enterpriseId, $userId, $memUid, $workerId, $type)
    {
        if ($userId) {
            // 登录者查询绑定关系
            $currentUserBinding = $this->getUserBinding($enterpriseId, $userId, $type);

            // 用户变员工，解绑原有绑定关系
            if ($currentUserBinding) {
                Log::info('用户分享-用户变员工，解绑原有绑定关系', $currentUserBinding->toArray());
                $this->unbinding($enterpriseId, $userId, $type);
            }
        }

        $currentStaffBinding = $this->getStaffBinding($enterpriseId, $workerId, $type);

        if ($currentStaffBinding) {
            return $currentStaffBinding;
        }

        $staffInitialParams = [
            'id' => idGenerate(),
            'enterprise_id'  => $enterpriseId,
            'user_id'        => $userId,
            'mem_uid'        => $memUid,
            'worker_id'      => $workerId,
            'type'           => $type,
            'parent_user_id' => $userId,
            'level'          => 0,
            'binding_type'   => StaffShare::BINDING_TYPE_YES
        ];

        Log::info('用户分享-员工信息初始化', $staffInitialParams);
        $staffInitial = StaffShare::query()->create($staffInitialParams);

        return $staffInitial ?? [];
    }

    /**
     * 内购 员工顾客绑定关系
     * 包含抢客判断
     *
     * @param $enterpriseId
     * @param $userId
     * @param $parentUserId
     * @param $workerId
     * @param $type
     * @return StaffShare|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|object|null
     * @throws \Exception
     */
    public function staffBinding($enterpriseId, $userId, $parentUserId, $workerId, $type)
    {
        // 登陆者查询绑定关系
        $currentUserBinding = $this->getUserBinding($enterpriseId, $userId, $type);

        if (empty($parentUserId) && empty($workerId)) {
            Log::info('用户分享-parentUserId 和 workerId 都为空');
            return $currentUserBinding;
        }

        // 分享者查询绑定关系，分享者可能是员工、客户或合伙人，
        $parentUser = $workerId
            ? $this->getStaffBinding($enterpriseId, $workerId, $type)
            : $this->getUserBinding($enterpriseId, $parentUserId, $type);

        // 登陆者无绑定、分享者无绑定
        // 不绑定，展示 企业logo -> $preUserAvatar -> $currentUserAvatar
        if (empty($currentUserBinding) && empty($parentUser)) {
            $parentUserUnbindingParams = [
                'id' => idGenerate(),
                'enterprise_id'  => $enterpriseId,
                'user_id'        => $userId,
                'mem_uid'        => '',
                'parent_user_id' => $parentUserId,
                'worker_id'      => $workerId,
                'level'          => 0,
                'qr_code'        => '',
                'type'           => $type,
                'binding_type'   => StaffShare::BINDING_TYPE_PARENT_UNBINDING // 客户分享，分享人无绑定
            ];
            Log::info('用户分享-登陆者无绑定,分享者无绑定', $parentUserUnbindingParams);
            return StaffShare::query()->create($parentUserUnbindingParams);
        }


        if ($workerId && $currentUserBinding && ($workerId == $currentUserBinding['worker_id'])) {
            Log::info('用户分享-同一数据请求多次，相同user_id,相同woker_id', $currentUserBinding->toArray());
            return $currentUserBinding;
        }

        if (empty($workerId) && $currentUserBinding && $parentUserId && ($parentUserId == $currentUserBinding['parent_user_id'])) {
            Log::info('用户分享-同一数据请求多次，相同user_id,相同parent_user_id', $currentUserBinding->toArray());
            return $currentUserBinding;
        }

        // 获取是否可以抢客的判断
        $canRobbing = $this->canRobbingCustomer($enterpriseId, $currentUserBinding, $type);

        // 登陆者有绑定关系、但是不可以被抢，直接返回登陆者当前绑定关系
        if ($currentUserBinding && !$canRobbing) {
            Log::info('用户分享-登陆者有绑定关系、但是不可以被抢，直接返回登陆者当前绑定关系', $currentUserBinding->toArray());
            return $currentUserBinding;
        }

        // 登陆者有绑定关系
        // 登陆者可以被抢（已过保护期）
        // 解绑原来的，再绑定当前的员工
        if ($currentUserBinding && $canRobbing) {
            Log::info('用户分享-抢客解绑', $currentUserBinding->toArray());
            $this->unbinding($enterpriseId, $userId, $type);
        }

        if ($parentUser) {
            $createParams = (0 == $parentUser['level'])  // 区分是员工分享还是客户分享
                ? [
                    'id' => idGenerate(),
                    'enterprise_id'  => $enterpriseId,
                    'user_id'        => $userId,
                    'mem_uid'        => $parentUser['mem_uid'] ?? '',
                    'worker_id'      => $workerId ?: ($parentUser['worker_id'] ?? 0),
                    'parent_user_id' => $parentUserId,
                    'level'          => 1,  // 员工分享
                    'type'           => $parentUser['type'] ?? $type,
                    'binding_type'   => StaffShare::BINDING_TYPE_YES   // 绑定  展示 staff -> $currentUserAvatar
                ] : [
                    'id' => idGenerate(),
                    'enterprise_id'  => $enterpriseId,
                    'user_id'        => $userId,
                    'mem_uid'        => $parentUser['mem_uid'] ?? '',  // 绑定 登陆者 至 分享者的所属员工
                    'worker_id'      => $workerId ?: ($parentUser['worker_id'] ?? 0),  // 绑定 登陆者 至 分享者的所属员工
                    'parent_user_id' => $parentUserId,
                    'level'          => $parentUser['level'] + 1, // 客户分享
                    'type'           => $parentUser['type'] ?? $type,
                    'binding_type'   => StaffShare::BINDING_TYPE_YES  //  展示 staff -> $preUserAvatar -> $currentUserAvatar
                ];
            Log::info('用户分享-绑定逻辑', $createParams);
            return StaffShare::query()->create($createParams);
        }

        return $currentUserBinding;
    }


    /**
     * 是否可以抢客判断
     *
     * @param $enterpriseId
     * @param $currentUserBinding
     * @param $type
     *
     * @return int
     * @throws \Exception
     */
    public function canRobbingCustomer($enterpriseId, $currentUserBinding, $type)
    {
        $canRobbing = false;

        $guideSetting = (new StaffManagerService())->getStaffSetting($enterpriseId, Constant::LOGIN_IDENTITY_GUIDE);
        $setting = $guideSetting ?? [];

        // 判断是否可以抢客
        if ($setting
            && $currentUserBinding
            && self::SETTING_CHANGE_OWNER_ALLOW == $setting['change_owner']  // 是否允许抢客 1:不允许 2：允许
            && ($setting['protected_time'] == 0 || time() > ($currentUserBinding['created_at'] + $setting['protected_time'] * 60 * 60 * 24)) // 时效为0 或 不在保护期内，可以抢
            && 0 <> $currentUserBinding['level']  // 用户分享等级大于1，即非员工可被抢客
        ) {
            $canRobbing = true;
        }

        return $canRobbing;
    }


    /**
     * 分享链路列表 员工头像、分享者、登陆者，仅返回这3个
     * 业务层可根据 memuid组装员工信息，当前顾客和分享者这里返回了头像昵称
     *
     * @param $enterpriseId
     * @param $userId
     * @param $type
     *
     * @return StaffShare|\Illuminate\Database\Eloquent\Builder|\Illuminate\Database\Eloquent\Model|mixed|object|null
     */
    public function getStaffChain($enterpriseId, $userId, $type)
    {
        return StaffShare::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'user_id'       => $userId,
                'type'          => $type,
            ])->whereIn('binding_type', [StaffShare::BINDING_TYPE_YES, StaffShare::BINDING_TYPE_PARENT_UNBINDING])
            ->with(['user', 'parentUser'])
            ->first();
    }

    protected function getUserBinding($enterpriseId, $userId, $type)
    {
        return StaffShare::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'user_id'       => $userId,
                'binding_type'  => StaffShare::BINDING_TYPE_YES,
            ])
            ->whereIn('type', [Constant::LOGIN_IDENTITY_GUIDE, Constant::LOGIN_IDENTITY_PARTNER])
            ->where('level', '<>', 0)
            ->first();
    }

    protected function getStaffBinding($enterpriseId, $workerId, $type)
    {
        return StaffShare::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'worker_id'     => $workerId,
                'level'         => 0,
                'binding_type'  => StaffShare::BINDING_TYPE_YES,
                'type'          => $type,
            ])
            ->first();
    }

    /**
     * 解绑逻辑，抢客解绑、客户变员工先解绑
     * 该解绑只解绑登陆者自己的绑定关系
     *
     * @param $enterpriseId
     * @param $userId
     * @param $type
     *
     * @return bool|int
     */
    protected function unbinding($enterpriseId, $userId, $type)
    {
        $unbindingParams = [
            'enterprise_id' => $enterpriseId,
            'user_id'       => $userId,
            'binding_type'  => StaffShare::BINDING_TYPE_YES,
        ];
        if (in_array($type, [Constant::LOGIN_IDENTITY_GUIDE, Constant::LOGIN_IDENTITY_PARTNER])) {
            Log::info('用户分享-解绑逻辑', $unbindingParams);
            return StaffShare::query()->where($unbindingParams)->whereIn('type', [Constant::LOGIN_IDENTITY_GUIDE, Constant::LOGIN_IDENTITY_PARTNER])->update(['binding_type' => StaffShare::BINDING_TYPE_NO]);
        } else {
            $unbindingParams['type'] = $type;
            Log::info('用户分享-解绑逻辑', $unbindingParams);
        }
        return StaffShare::query()->where($unbindingParams)->update(['binding_type' => StaffShare::BINDING_TYPE_NO]);
    }

    /**
     * 解绑逻辑，员工离职
     * 该解绑会影响该员工的所有被分享者
     *
     * @param $enterpriseId
     * @param $memUids
     * @param $workerIds
     * @param $type
     *
     * @return bool|int
     */
    public function unbindingLeave(int $enterpriseId, array $memUids, array $workerIds, int $type)
    {
        $unbindingParams = [
            'enterprise_id' => $enterpriseId,
            'binding_type'  => StaffShare::BINDING_TYPE_YES,
            'type'          => $type,
        ];
        Log::info('用户分享-解绑逻辑-员工离职', [$unbindingParams, $memUids]);

        $bindingList = StaffShare::query()
            ->where($unbindingParams)
            ->when($memUids, function ($query) use ($memUids) {
                return $query->whereIn('mem_uid', $memUids);
            }, function ($query) use ($workerIds) {
                return $query->whereIn('worker_id', $workerIds);
            })
            ->get()
            ->toArray();
        if (!empty($bindingList)) {
            $updatedAt = time();
            StaffShare::query()->whereIn('id', array_column($bindingList, 'id'))->update([
                'binding_type' => StaffShare::BINDING_TYPE_NO,
                'updated_at' => $updatedAt
            ]);

            $shareLinkMessageList = [];
            $notifyDataList = [];
            foreach ($bindingList as $binding) {

                $notifyDataList[] = [
                    'enterprise_id' => $enterpriseId,
                    'third_user_id' => $binding['user_id'],
                    'curr_worker_id' => $binding['worker_id'],
                    'curr_worker_mem_uid' => $binding['mem_uid'],
                    'bind_time' => intval($binding['created_at']),
                ];
                array_push($shareLinkMessageList, ShareLinkMessageAdapter::deleteMessage($binding, $updatedAt));
            }

            dispatch(new ShareLinkBatchUpdateCustomerStatisticJob(MessageConstant::SHARE_LINK_BATCH_LEAVE,
                $shareLinkMessageList));

            (new ShareLinkChangeEventService())->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_DELETE, $notifyDataList);
        }
    }

    /**
     * 解绑逻辑，员工变导购的情况
     * 该解绑只解绑登陆者自己的绑定关系
     *
     * @param $enterpriseId
     * @param $memUids
     * @param $workerIds
     * @param $type
     *
     * @return bool|int
     */
    public function unbindingChange(int $enterpriseId, array $memUids, array $workerIds, int $type)
    {
        $userIds = StoreWorker::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->where('user_id', '<>', 0)
            ->when($memUids, function ($query) use ($memUids) {
                return $query->whereIn('mem_uid', $memUids);
            }, function ($query) use ($workerIds) {
                return $query->whereIn('id', $workerIds);
            })
            ->pluck('user_id')->all();

        $unbindingParams = [
            'enterprise_id' => $enterpriseId,
            'binding_type'  => StaffShare::BINDING_TYPE_YES,
            'type'          => $type,
        ];

        Log::info('用户分享-解绑逻辑-员工变导购', [$unbindingParams, $userIds, $memUids]);

        $bindingList = StaffShare::query()
            ->where($unbindingParams)
            ->whereIn('user_id', $userIds)
            ->get()
            ->toArray();
        if (!empty($bindingList)) {
            $updatedAt = time();
            StaffShare::query()->whereIn('id', array_column($bindingList, 'id'))->update([
                'binding_type' => StaffShare::BINDING_TYPE_NO,
                'updated_at' => $updatedAt
            ]);

            $shareLinkMessageList = [];
            $notifyDataList = [];
            foreach ($bindingList as $binding) {
                $notifyDataList[] = [
                    'enterprise_id' => $enterpriseId,
                    'third_user_id' => $binding['user_id'],
                    'curr_worker_id' => $binding['worker_id'],
                    'curr_worker_mem_uid' => $binding['mem_uid'],
                    'bind_time' => intval($binding['created_at']),
                ];
                array_push($shareLinkMessageList, ShareLinkMessageAdapter::deleteMessage($binding, $updatedAt));
            }
            dispatch(new ShareLinkBatchUpdateCustomerStatisticJob(MessageConstant::SHARE_LINK_BATCH_CREATE_GUIDE,
                $shareLinkMessageList));

            (new ShareLinkChangeEventService())->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_DELETE, $notifyDataList);
        }
    }

}
