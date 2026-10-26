<?php


namespace App\Services;


use App\Models\StaffShare;
use App\Models\StoreWorker;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Shop\User;
use Illuminate\Support\Facades\Redis;

class StaffShareService extends BaseService
{

    /**
     * 绑定链路列表
     *
     * @param array $input
     * @return array
     */
    public function staffList($input = [])
    {
        $ids = $input['ids'] ?? [];
        $workerId = $input['worker_id'] ?? 0;
        $workerIds = $input['worker_ids'] ?? [];
        $user_ids = $input['user_ids'] ?? [];
        $mem_uids = $input['mem_uids'] ?? [];
        $mem_uid = $input['mem_uid'] ?? '';
        $user_id = $input['user_id'] ?? '';
        $page = $input['page'] ?? 1;
        $size = $input['size'] ?? 10;
        $isPage = $input['is_page'] ?? 1;
        $excludeUserIdList = $input['exclude_user_ids'] ?? [];
        $startTime = $input['start_time'] ?? 0;
        $endTime = $input['end_time'] ?? 0;

        $queryBuilder = StaffShare::query()
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->when(!empty($user_id), function ($query) use ($user_id) {
                $query->where('user_id', $user_id);
            })
            ->when(!empty($mem_uid), function ($query) use ($mem_uid) {
                $query->where('mem_uid', $mem_uid);
            })
            ->when(!empty($workerId), function ($query) use ($workerId) {
                $query->where('worker_id', $workerId);
            })
            ->when(!empty($user_ids), function ($query) use ($user_ids) {
                $query->whereIn('user_id', $user_ids);
            })
            ->when(!empty($ids), function ($query) use ($ids) {
                $query->whereIn('id', $ids);
            })
            ->when(!empty($mem_uids), function ($query) use ($mem_uids) {
                $query->whereIn('mem_uid', $mem_uids);
            })
            ->when(!empty($workerIds), function ($query) use ($workerIds) {
                $query->whereIn('worker_id', $workerIds);
            })
            ->when(!empty($excludeUserIdList), function ($query) use ($excludeUserIdList) {
                $query->whereNotIn('user_id', $excludeUserIdList);
            })
            ->when(!empty($startTime) && !empty($endTime), function ($query) use ($startTime, $endTime) {
                $query->whereBetween('created_at', [$startTime, $endTime]);
            })
            ->orderByDesc('created_at');
        $count = $queryBuilder->count();
        if ($isPage) {
            $queryBuilder->forPage($page, $size);
        }
        $data = $queryBuilder->get()->toArray();

        return [
            'total' => $count,
            'data' => $data
        ];
    }

    public function createStaffShareData($memUid, $parentUserId, $userId, $enterpriseId, $shareParentUserInfo, $shareUserInfo, $isStaff)
    {
        // 经由多少人传递 n
        $level = 0;

        // 二维码(单独调用接口了)
        $qrCode = '';

        // 解绑关系（后续用户分享进来不绑定）
        if (empty($memUid) && empty($shareParentUserInfo)) {
            return $memUid;
        }

        if (!empty($shareUserInfo) && !$isStaff && $shareUserInfo['mem_uid'] == $memUid
            && $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_YES) {
            return $memUid;
        }

        // 上一级数据
        if (!empty($shareParentUserInfo)) {
            // 如果解绑，后面人不绑定
            if ($shareParentUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO) {
                return $memUid;
            }
            $level = $shareParentUserInfo['level'] + 1;
        }

        $insertData = [
            'id' => idGenerate(),
            'enterprise_id' => $enterpriseId,
            'user_id' => $userId,
            'mem_uid' => empty($memUid) ? $shareParentUserInfo['mem_uid'] : $memUid,
            'parent_user_id' => $parentUserId ?? 0,
            'level' => $level,
            'qr_code' => $qrCode
        ];

        StaffShare::query()->create($insertData);

        return $memUid;
    }


    public function clearCounter($userId)
    {
        $redis = Redis::connection();

        $key = 'crs:shop:staff:member:counter:' . $userId;

        if ($redis->get($key)) {
            $redis->del([$key]);
        }

        return true;
    }

    public function setCounter($userId)
    {
        $redis = Redis::connection();

        $key = 'crs:shop:staff:member:counter:'.$userId;

        $counter = $redis->incr($key);

        return $counter;
    }

    public function getShareParentInfo($memUid, $parentUserId, $userId, $enterpriseId, $shareUserInfo, $isStaff)
    {
        if ($isStaff) {
            $memUid = $this->getUserInfoById($enterpriseId, $userId)['mem_uid'];
            return [$memUid, [], 0];
        }

        $shareParentUserInfo = [];
        $bool = true;
        if (!empty($parentUserId) && (empty($shareUserInfo) || $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO)) {
            $pOpenId = $this->getUserInfoById($enterpriseId, $parentUserId);
            $pRight = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, '', '');
            if ($pRight) {
                $bool = false;
                $memUid = $pOpenId['mem_uid'];
                $parentUserId = 0;
            }
        }

        if (!empty($memUid) && $bool && (empty($shareUserInfo) || $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO)) {
            $pUserArr = $this->getUserInfoByMemUid($enterpriseId, $memUid);
            if (!empty($pUserArr)) {
                $right = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, '', '');
                if (!$right) {
                    $memUid = '';
                    $parentUserId = $pUserArr['id'];
                }
            }
        }

        if (!empty($parentUserId)) {
            // 上级解绑也需要展示(获取最近一条数据展示)
            $shareParentUserInfo = (new StaffShare())
                ->getStaffShareInfoByCondition(['user_id' => $parentUserId], false);

            if (empty($shareParentUserInfo)) {
                return [$memUid, $shareParentUserInfo, $parentUserId];
            }

            // 同一个用户 重置父id
            if ($userId == $parentUserId) {

                $parentUserId = $shareParentUserInfo['parent_user_id'] ?? 0;

                // 重置上级关系
                if (!empty($parentUserId)) {
                    $shareParentUserInfo = (new StaffShare())
                        ->getStaffShareInfoByCondition(['user_id' => $parentUserId], false);
                }
            }

            // 用户分享 赋值mem_uid
            $memUid = $shareParentUserInfo['binding_type'] == StaffShare::BINDING_TYPE_YES ? $shareParentUserInfo['mem_uid'] : '';

            // 获取上级mem_uid
            if ($parentUserId != 0) {
                $parentMemUid = $this->getUserInfoById($enterpriseId, $userId);
                $parentMemUidRight = (new StaffManagerService())->staffIsBuy($enterpriseId, $memUid, '', '');
                if ($parentMemUid['mem_uid'] == $memUid && !empty($parentMemUid) && $parentMemUidRight) {
                    $parentUserId = 0;
                }
            }

            // 不存在父级
            if (empty($parentUserId)) {
                $shareParentUserInfo = [];
            }
        }

        return [$memUid, $shareParentUserInfo, $parentUserId];
    }

    /**
     * 查询链路
     * @param array $params
     * @return array
     */
    public function detail(array $params)
    {
        $columns = ['id',
            'enterprise_id',
            'store_id',
            'type',
            'user_id',
            'parent_user_id',
            'mem_uid',
            'worker_id',
            'level',
            'binding_type',
            'created_at',
            'updated_at'];
        if (empty($params['enterprise_id']) && empty($params['user_id']) && empty($params['store_id']) && empty($params['worker_id']) && empty($params['type']) && empty($params['types']) && empty($params['binding_type'])) {
            return [];
        }

        $staffShare = StaffShare::query()
            ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                return $query->where('enterprise_id', '=', $params['enterprise_id']);
            })
            ->when(!empty($params['user_id']), function ($query) use ($params) {
                return $query->where('user_id', '=', $params['user_id']);
            })
            ->when(!empty($params['worker_id']), function ($query) use ($params) {
                return $query->where('id', '=', $params['worker_id']);
            })
            ->when(!empty($params['store_id']), function ($query) use ($params) {
                return $query->where('store_id', '=', $params['store_id']);
            })
            ->when(!empty($params['type']), function ($query) use ($params) {
                if (is_array($params['type'])) {
                    return $query->whereIn('type', $params['type']);
                } else {
                    return $query->where('type', $params['type']);
                }
            })
            ->when(!empty($params['binding_type']), function ($query) use ($params) {
                if (is_array($params['binding_type'])) {
                    return $query->whereIn('binding_type', $params['binding_type']);
                } else {
                    return $query->where('binding_type', $params['binding_type']);
                }
            })
            ->first($columns);

        return !empty($staffShare) ? $staffShare->toArray() : [];
    }

    /**
     * 获取用户Ids
     * @param $params
     * @return array
     */
    public function getUserIds(array $params)
    {
        switch ($params['ascription']) {
            case StaffShare::ACHIEVEMENT_ASCRIPTION_SELF :
                $userIds = StaffShare::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->where('worker_id', $params['worker_id'])
                    ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                    ->pluck('user_id')
                    ->toArray();
//                if ($params['worker_user_id']) {
//                    $userIds[] = $params['worker_user_id'];
//                }
                break;
            case StaffShare::ACHIEVEMENT_ASCRIPTION_OTHER :
                $userInfoIds = array_diff($params['user_ids'], [$params['worker_user_id']]);
                $userIds = StaffShare::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->where('worker_id', '<>', $params['worker_id'])
                    ->whereIn('user_id', $userInfoIds)
                    ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                    ->pluck('user_id')
                    ->toArray();
                $getListParams = [
                    'enterprise_id' => $params['enterprise_id'],
                    'worker_status' => StoreWorker::IN_SERVICE
                ];
                if (!empty($params['user_ids'])) {
                    $getListParams['user_ids'] = $userInfoIds;
                }
                $storeWorkerList = StoreWorkerService::list($getListParams, false)['list'];
                $storeWorkerUserIds = array_column($storeWorkerList, 'user_id');

                $userIds = array_unique(array_merge($userIds, $storeWorkerUserIds));
                if (!empty($params['worker_user_id'])) {
                    $workerStaffShare = StaffShare::query()
                        ->where('enterprise_id', $params['enterprise_id'])
                        ->where('user_id', $params['worker_user_id'])
                        ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                        ->first();
                    if (!empty($workerStaffShare) && $workerStaffShare['worker_id'] != $params['worker_id']) {
                        array_push($userIds, $params['worker_user_id']);
                    }
                }
                break;
            case StaffShare::ACHIEVEMENT_ASCRIPTION_NONE :
                $userIds = StaffShare::query()
                    ->where('enterprise_id', $params['enterprise_id'])
                    ->whereIn('user_id', $params['user_ids'])
                    ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                    ->pluck('user_id')
                    ->toArray();
                $getListParams = [
                    'enterprise_id' => $params['enterprise_id'],
                    'worker_status' => StoreWorker::IN_SERVICE
                ];
                if (!empty($params['user_ids'])) {
                    $getListParams['user_ids'] = $params['user_ids'];
                }
                $storeWorkerList = StoreWorkerService::list($getListParams, false)['list'];
                $storeWorkerUserIds = array_column($storeWorkerList, 'user_id');

                $userIds = array_diff($params['user_ids'], $userIds);
                $userIds = array_diff($userIds, $storeWorkerUserIds);
                break;
            default :
                $userIds = [];
        }

        return $userIds;
    }

    /**
     * 查询绑定的导购信息
     * @param int $user_id
     * @return array
     */
    public function getBindWorkerInfo(int $user_id)
    {
        $storeWorker = StaffShare::query()
            ->join('store_worker AS sw', 'sw.id', '=', 'staff_share.worker_id')
            ->where('staff_share.user_id', '=', $user_id)
            ->where('staff_share.binding_type', '=', StaffShare::BINDING_TYPE_YES)
            ->first(['sw.*']);

        return !empty($storeWorker) ? $storeWorker->toArray() : [];
    }

    public function getStaffShareList($memUid, $parentUserId, $userId, $enterpriseId, $isStaff, $shareUserInfo)
    {
        // 初始化
        $list = $ParentParam = [];

        // 当前用户
        $userInfo  = $this->getUserInfoById($enterpriseId, $userId);
        // 当前用户mem_uid
        $memId = $userInfo['mem_uid'];

        // 自己是员工情况
        if ($isStaff) {
            $memUid = $memId;
        }

        // 解绑展示后台配置logo（产品需求）
        if (!$isStaff && !empty($shareUserInfo) && $shareUserInfo['binding_type'] == StaffShare::BINDING_TYPE_NO) {
            $memUid = '';
        }

        // 获取企业domain
        $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId]);

        // 认证员工
        if (!empty($memUid)) {
            $result = (new StaffManagerService())->getStaffListFromUc($enterpriseInfo['data']['domain'], 'contact', ['memUids' => [$memUid]]);
            $memInfo = (new StaffManagerService())->formatUcMember($result);

            // 打通企业微信
            $memUserId = $this->getUserInfoByMemUid($enterpriseId, $memUid)['id'];
        }

        $appInfo = (new Enterprise())->enterpriseMerchantSettingInfo(['enterprise_id' => $enterpriseId]);

        $memParam = [
            'avatar'    => $memInfo['data'][0]['mem_face'] ?? ($appInfo['data']['logo'] ?? ''),
            'nick_name' => $memInfo['data'][0]['mem_username'] ?? ($appInfo['data']['short_name'] ?? ''),
            'type'      => 1,
            'user_id'   => (int) ($memUserId ?? 0)
        ];
        $memParam['avatar'] = str_replace('http:', 'https:', $memParam['avatar']);
        array_push($list, $memParam);

        // 上一级员工
        if (!empty($parentUserId) && !$isStaff) {
            $userParentInfo = $this->getUserInfoById($enterpriseId, $parentUserId);
            $ParentParam    = [
                'avatar'    => $userParentInfo['avatar'] ?? '',
                'nick_name' => $userParentInfo['nick_name'] ?? '',
                'type'      => 2,
                'user_id'   => (int) $parentUserId
            ];
            $ParentParam['avatar'] = str_replace('http:', 'https:', $ParentParam['avatar']);
            array_push($list, $ParentParam);
        }

        $orderMoney = $count = 0;
        // 是内购员工 （计数器加1）
        if ($isStaff) {
            // 今日业绩
            $search = [
                'type'           => 1,
                'enterprise_id'  => $enterpriseId,
                'qy_domain'      => $enterpriseInfo['data']['domain'],
                'app_identifier' => 'contact',
                'mem_uid'        => $memId
            ];
            $order = (new StaffOrderService())->getMyStaff($search);

            $orderMoney = $order['data']['total_money'] ?? 0;

            $count = $this->setCounter($userId);
        }

        // 不是员工清除计数器
        if (!$isStaff) {
            $this->clearCounter($userId);
        }

        $userParam = [
            'avatar'    => $userInfo['avatar'] ?? '',
            'nick_name' => $userInfo['nick_name'] ?? '',
            'type'      => 3,
            'user_id'   => (int) $userId,
            'mem_uid'   => $memId,
            'is_staff'  => $isStaff,
            'count'     => $count,
            'order_money' => $orderMoney
        ];
        $userParam['avatar'] = str_replace('http:', 'https:', $userParam['avatar']);
        array_push($list, $userParam);

        return [$list, $appInfo['data']['short_name'] ?? ''];
    }

    public function getUserInfoById($enterpriseId, $userId)
    {
        return (new User())->findUser([
            'enterprise_id' => $enterpriseId,
            'id' => $userId
        ]);
    }

    public function getUserInfoByMemUid($enterpriseId, $memUid)
    {
        return (new User())->findUser([
            'enterprise_id' => $enterpriseId,
            'mem_id' => $memUid
        ]);
    }

    public static function getUserBinding($enterpriseId, $userId)
    {
        $binding = StaffShare::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('user_id', $userId)
            ->where('binding_type', StaffShare::BINDING_TYPE_YES)
            ->first();
        return $binding ? $binding->toArray() : [];
    }

}
