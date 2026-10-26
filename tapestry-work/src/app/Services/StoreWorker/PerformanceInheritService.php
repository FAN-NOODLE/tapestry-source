<?php


namespace App\Services\StoreWorker;


use App\Common\MessageConstant;
use App\Exceptions\Error;
use App\Jobs\Message\ShareLinkMessage;
use App\Jobs\Message\ShareLinkMessageAdapter;
use App\Jobs\ShareLinkBatchUpdateCustomerStatisticJob;
use App\Jobs\ShareLinkUpdateCustomerStatisticJob;
use App\Library\Ddc\DdcMessage;
use App\Library\Uc\Member;
use App\Library\Util\Logger;
use App\Models\MemberUserRelation;
use App\Models\StaffShare;
use App\Models\StoreWorker;
use App\Services\BaseService;
use App\Services\Events\ShareLinkChangeEventService;
use App\Services\StaffSettingService;
use App\Services\StaffShareService;
use App\Services\StoreWorkerService;
use CrsSdk\MicroUser\User;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class PerformanceInheritService extends BaseService
{

    // 按照顾客昵称搜索
    const SEARCH_TYPE_NICK_NAME = 1;
    // 按照顾客手机号搜索
    const SEARCH_TYPE_MOBILE = 2;

    const SEARCH_TYPE_LIST = [
        self::SEARCH_TYPE_NICK_NAME,
        self::SEARCH_TYPE_MOBILE
    ];

    // 按添加好友时间搜索
    const SEARCH_TIME_TYPE_ADD_CONTACT = 1;
    // 按业绩绑定时间搜索
    const SEARCH_TIME_TYPE_PERFORMANCE_BINDING = 2;

    const SEARCH_TIME_TYPE_LIST = [
        self::SEARCH_TIME_TYPE_ADD_CONTACT,
        self::SEARCH_TIME_TYPE_PERFORMANCE_BINDING
    ];

    // 是否添加企微好友-是
    const ADD_CONTACT_YES = 1;

    // 是否添加企微好友-否
    const ADD_CONTACT_NO = 2;

    /**
     * 导购业绩继承-获取和导购绑定业绩的客户列表
     * @param array $params
     * @return array
     * @throws \Exception
     */
    public function customerList(array $params)
    {
        $worker = StoreWorkerService::detail([
            'enterprise_id' => $params['enterprise_id'],
            'worker_id'     => $params['worker_id'],
        ]);
        if (empty($worker)) {
            throw new Error(110010, "导购信息不存在", 422);
        }

        $searchParams      = [
            'search_type'      => $params['search_type'],
            'search_value'     => $params['search_value'],
            'search_time_type' => $params['search_time_type'],
            'start_time'       => $params['start_time'],
            'end_time'         => $params['end_time'],
            'page'             => $params['page'],
            'size'             => $params['size'],
        ];
        $bindingListResult = $this->getBindingList($params['enterprise_id'], $worker, $searchParams, true);
        if (empty($bindingListResult['data'])) {
            return [
                'total' => $bindingListResult['total'],
                'data'  => $bindingListResult['data']
            ];
        }
        $total       = $bindingListResult['total'];
        $bindingList = $bindingListResult['data'];

        // 查询用户数据
        $queryUserParams = [
            'is_paginate'   => 0,
            'enterprise_id' => $params['enterprise_id'],
            'conditions'    => [
                ['id', 'in', array_column($bindingList, 'user_id')]
            ]
        ];
        $userList        = (new User())->userList($queryUserParams)['list'];
        $userMap         = array_column($userList, null, 'id');


        // 查询好友信息信息
        $contactMap  = [];
        $unionIdList = array_filter(array_column($userMap, 'union_id'));
        if (!empty($unionIdList)) {
            MemberUserRelation::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('mem_uid', $worker['mem_uid'])
                ->whereIn('union_id', $unionIdList)
                ->get()
                ->each(function (MemberUserRelation $v) use (&$contactMap) {
                    $contactMap[$v['union_id']] = $v->toArray();
                });
        }


        // 拼装返回结果
        $customerList = [];
        foreach ($bindingList as $bindingInfo) {
            $userInfo = $userMap[$bindingInfo['user_id']] ?? null;
            if (empty($userInfo)) {
                Log::info("PerformanceInheritService.customerList userInfo is empty", [
                    'user_id' => $bindingInfo['user_id']
                ]);
                continue;
            }

            $customerInfo                             = [];
            $customerInfo['user_id']                  = $userInfo['id'];
            $customerInfo['nick_name']                = $userInfo['nick_name'];
            $customerInfo['avatar']                   = $userInfo['avatar'];
            $customerInfo['sex']                      = $userInfo['sex'];
            $customerInfo['mobile']                   = $userInfo['mobile'];
            $customerInfo['whether']                  = self::ADD_CONTACT_NO;
            $customerInfo['add_contact_time']         = '';
            $customerInfo['performance_binding_time'] = date("Y-m-d H:i:s", $bindingInfo['created_at']);

            $contactInfo = null;
            if (!empty($userInfo['union_id']) && !empty($contactMap[$userInfo['union_id']])) {
                $contactInfo = $contactMap[$userInfo['union_id']];
                if (MemberUserRelation::DEL_DEFAULT == $contactInfo['delete_type']) {
                    $customerInfo['whether'] = self::ADD_CONTACT_YES;
                }
                if (empty($customerInfo['nick_name'])) {
                    $customerInfo['nick_name'] = $contactInfo['user_remark_name'] ?? $contactInfo['ec_name'];
                }
                if (empty($customerInfo['avatar'])) {
                    $customerInfo['avatar'] = $contactInfo['ec_avatar'];
                }
                if (empty($customerInfo['sex']) || !in_array($customerInfo['sex'], [1, 2])) {
                    $customerInfo['sex'] = $contactInfo['ec_gender'];
                }
                $customerInfo['add_contact_time'] = date("Y-m-d H:i:s", $contactInfo['created_at']);
            }
            $customerList[] = $customerInfo;
        }

        return [
            'total' => $total,
            'data'  => $customerList
        ];
    }

    private function getBindingList(int $enterpriseId, array $worker, array $searchParams, bool $isPage = true)
    {

        // 查询联系人列表
        $searchContactMap = [];
        if (!empty($searchParams['search_time_type'])
            && self::SEARCH_TIME_TYPE_ADD_CONTACT == $searchParams['search_time_type']
            && !empty($searchParams['start_time']) && !empty($searchParams['end_time'])) {
            MemberUserRelation::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('mem_uid', $worker['mem_uid'])
                ->where('union_id', '!=', '')
                ->whereBetween('created_at', [$searchParams['start_time'], $searchParams['end_time']])
                ->get()
                ->each(function (MemberUserRelation $v) use (&$searchContactMap) {
                    $searchContactMap[$v['union_id']] = $v->toArray();
                });
            if (empty($searchContactMap)) {
                Log::info("PerformanceInheritService.customerList searchContact is empty");
                return [
                    'total' => 0,
                    'data'  => []
                ];
            }
        }

        // 根据用户筛选条件查询用户列表
        $searchUserMap    = [];
        $searchUserParams = [
            'is_paginate'   => 0,
            'enterprise_id' => $enterpriseId,
            'conditions'    => []
        ];
        if (!empty($searchContactMap)) {
            $searchUserParams['conditions'][] = [
                'union_id',
                'in',
                array_keys($searchContactMap)
            ];
        }
        if (!empty($searchParams['search_type']) && !empty($searchParams['search_value'])) {
            $searchValue = trim($searchParams['search_value']);
            if (self::SEARCH_TYPE_NICK_NAME == $searchParams['search_type']) {
                $searchUserParams['conditions'][] = [
                    'nick_name',
                    'like',
                    $searchValue
                ];
            } else if (self::SEARCH_TYPE_MOBILE == $searchParams['search_type']) {
                $searchUserParams['conditions'][] = [
                    'mobile',
                    'like',
                    $searchValue
                ];
            }

            // 商城用户表数据较大，模糊匹配查询较慢，先查询分享链路id缩小检索范围
            if (!empty($searchParams['search_time_type'])
                && self::SEARCH_TIME_TYPE_PERFORMANCE_BINDING == $searchParams['search_time_type']) {
                $bindingAllQueryParams = [
                    'enterprise_id' => $enterpriseId,
                    'worker_id'     => $worker['id'],
                    'is_page'       => 0,
                ];
                if (!empty($searchParams['start_time']) && !empty($searchParams['end_time'])) {
                    $bindingAllQueryParams['start_time'] = $searchParams['start_time'];
                    $bindingAllQueryParams['end_time']   = $searchParams['end_time'];
                }
                $bindingAllListResult = (new StaffShareService())->staffList($bindingAllQueryParams);
//                    (new StaffRight())->BindList($bindingAllQueryParams);
                if (!isset($bindingAllListResult['data']) || empty($bindingAllListResult['data'])) {
                    Log::info("PerformanceInheritService.customerList bindingAllUser is empty");
                    return [
                        'total' => 0,
                        'data'  => []
                    ];
                } else {
                    $searchUserParams['conditions'][] = [
                        'id',
                        'in',
                        array_column($bindingAllListResult['data'], 'user_id')
                    ];
                }
            }
        }
        if (!empty($searchUserParams['conditions'])) {
            $searchUserResult = (new User())->userList($searchUserParams);
            if (empty($searchUserResult['list'])) {
                Log::info("PerformanceInheritService.customerList searchUser is empty");
                return [
                    'total' => 0,
                    'data'  => []
                ];
            }
            $searchUserMap = array_column($searchUserResult['list'], null, 'id');
        }

        // 查询分享链路数据
        $bindingQueryParams = [
            'enterprise_id' => $enterpriseId,
            'worker_id'     => $worker['id'],
            'is_page'       => $isPage ? 1 : 0,
            'page'          => $isPage ? $searchParams['page'] ?? 1 : 0,
            'size'          => $isPage ? $searchParams['size'] ?? 10 : 0,
        ];
        // 用户搜索条件
        if (!is_null($searchUserMap)) {
            $bindingQueryParams['user_ids'] = array_keys($searchUserMap);
        }
        if (!empty($searchParams['user_ids'])) {
            $bindingQueryParams['user_ids'] = $searchParams['user_ids'];
        }
        // 排查当前导购
        if (!empty($worker['user_id'])) {
            $bindingQueryParams['exclude_user_ids'] = [$worker['user_id']];
        }
        // 业绩绑定时间筛选
        if (!empty($searchParams['search_time_type'])
            && self::SEARCH_TIME_TYPE_PERFORMANCE_BINDING == $searchParams['search_time_type']
            && !empty($searchParams['start_time']) && !empty($searchParams['end_time'])) {
            $bindingQueryParams['start_time'] = $searchParams['start_time'];
            $bindingQueryParams['end_time']   = $searchParams['end_time'];
        }
        $bindingListResult = (new StaffShareService())->staffList($bindingQueryParams);
//            (new StaffRight())->BindList($bindingQueryParams);
        if (!isset($bindingListResult['data']) || empty($bindingListResult['data'])) {
            Log::info("PerformanceInheritService.customerList bindingUser is empty");
            return [
                'total' => $bindingListResult['total'],
                'data'  => []
            ];
        }
        return [
            'total' => $bindingListResult['total'],
            'data'  => $bindingListResult['data']
        ];
    }

    /**
     * 导购业绩继承-部分客户继承
     * @param array $params
     * @throws \Exception
     */
    public function partialInherit(array $params)
    {
        $oldWorker = StoreWorkerService::detail([
            'enterprise_id' => $params['enterprise_id'],
            'worker_id'     => $params['old_worker_id'],
        ]);
        if (empty($oldWorker)) {
            throw new Error(110010, "导购信息不存在", 422);
        }

        $newWorker = StoreWorkerService::detail([
            'enterprise_id' => $params['enterprise_id'],
            'worker_id'     => $params['new_worker_id'],
        ]);
        if (empty($newWorker)) {
            throw new Error(110010, "导购信息不存在", 422);
        }

        $searchParams      = [
            'user_ids' => $params['user_ids'],
        ];
        $bindingListResult = $this->getBindingList($params['enterprise_id'], $oldWorker, $searchParams, false);
        if (!isset($bindingListResult['data']) || empty($bindingListResult['data'])) {
            throw new Error(110010, "继承用户列表为空");
        }

        $this->inherit($params['enterprise_id'], $oldWorker, $newWorker, $bindingListResult['data']);
    }

    /**
     * 导购业绩继承-全部客户继承
     * @param array $params
     * @throws \Exception
     */
    public function allInherit(array $params)
    {

        $oldWorker = StoreWorkerService::detail([
            'enterprise_id' => $params['enterprise_id'],
            'worker_id'     => $params['old_worker_id'],
        ]);
        if (empty($oldWorker)) {
            throw new Error(110010, "导购信息不存在", 422);
        }

        $newWorker = StoreWorkerService::detail([
            'enterprise_id' => $params['enterprise_id'],
            'worker_id'     => $params['new_worker_id'],
        ]);
        if (empty($newWorker)) {
            throw new Error(110010, "导购信息不存在", 422);
        }

        $searchParams      = [
            'search_type'      => $params['search_type'],
            'search_value'     => $params['search_value'],
            'search_time_type' => $params['search_time_type'],
            'start_time'       => $params['start_time'],
            'end_time'         => $params['end_time'],
        ];
        $bindingListResult = $this->getBindingList($params['enterprise_id'], $oldWorker, $searchParams, false);
        if (!isset($bindingListResult['data']) || empty($bindingListResult['data'])) {
            throw new Error(110010, "继承用户列表为空");
        }

        $this->inherit($params['enterprise_id'], $oldWorker, $newWorker, $bindingListResult['data']);
    }

    /**
     * 业绩继承
     * @param $enterpriseId
     * @param $oldWorker
     * @param $newWorker
     * @param $bindingList
     * @throws \Exception
     */
    public function inherit($enterpriseId, $oldWorker, $newWorker, $bindingList)
    {
        $now                  = time();
        $notifyDataList       = [];
        $insertBindingList    = [];
        $shareLinkMessageList = [];
        foreach ($bindingList as $binding) {
            $newBinding = [
                'id'             => idGenerate(),
                'enterprise_id'  => $enterpriseId,
                'store_id'       => $newWorker['store_id'],
                'type'           => $binding['type'],
                'user_id'        => $binding['user_id'],
                'parent_user_id' => $binding['parent_user_id'],
                'mem_uid'        => $newWorker['mem_uid'],
                'worker_id'      => $newWorker['id'],
                'level'          => $binding['level'],
                'qr_code'        => $binding['qr_code'],
                'created_at'     => $now,
                'updated_at'     => $now,
                'binding_type'   => 1
            ];
            array_push($insertBindingList, $newBinding);
            array_push($shareLinkMessageList, ShareLinkMessageAdapter::updateMessage($binding, $newBinding));
            array_push($notifyDataList, [
                'enterprise_id'       => $enterpriseId,
                'third_user_id'       => $binding['user_id'],
                'curr_worker_id'      => $newWorker['id'],
                'curr_worker_mem_uid' => $newWorker['mem_uid'],
                'old_worker_id'       => $oldWorker['id'],
                'old_worker_mem_uid'  => $oldWorker['mem_uid'],
                'bind_time'           => $now
            ]);
        }

        $unbindIdList = array_column($bindingList, 'id');

//        $inheritParams = [
//            'unbind_id_list' => $unbindIdList,
//            'insert_binding_list' => $insertBindingList
//        ];
        $this->InheritStaff($unbindIdList, $insertBindingList);
//        (new StaffRight())->InheritStaff($inheritParams);

        dispatch(new ShareLinkBatchUpdateCustomerStatisticJob(MessageConstant::SHARE_LINK_BATCH_INHERIT, $shareLinkMessageList));

        (new ShareLinkChangeEventService())->pushDdcMessage(DdcMessage::GUIDE_USER_BINDING_UPDATE, $notifyDataList);
    }

    public function InheritStaff($unbindIdList, $insertBindingList)
    {
        DB::beginTransaction();
        try {
            (new StaffSettingService())->unBindingStaffShareById($unbindIdList);
            StaffShare::query()->insert($insertBindingList);
            DB::commit();
        } catch (\Exception $e) {
            DB::rollback();
            throw new Error(100000, '业绩继承失败，请稍后重试', [], 500, $e);
        }
    }

    public function StaffList(array $params)
    {
        $enterpriseId = $params['enterprise_id'];
        $workerIds    = $params['worker_ids'] ?? '';
        $page         = $params['page'] ?? 0;
        $size         = $params['size'] ?? 0;
        $isPage       = $params['is_page'] ?? false;
        $queryBuilder = StaffShare::query()
            ->select(['staff_share.user_id', 'staff_share.worker_id', 'staff_share.created_at', 'w.mem_uid', 'w.worker_no', 'w.name'])
            ->leftJoin('store_worker as w', 'staff_share.worker_id', '=', 'w.id')
            ->where('staff_share.binding_type', StaffShare::BINDING_TYPE_YES)
            ->where('staff_share.enterprise_id', $enterpriseId)
            ->where('staff_share.user_id', '>', 0)
            ->where('staff_share.mem_uid', '<>', '')
            ->when(!empty($workerIds), function ($query) use ($workerIds) {
                $query->whereIn('worker_id', $workerIds);
            })
            ->orderByDesc('created_at');
        $count        = $queryBuilder->count();
        if ($isPage) {
            $queryBuilder->forPage($page, $size);
        }
        $data    = $queryBuilder->get()->toArray();
        $memUids = array_column($data, 'mem_uid');

        // 查询用户数据
        $queryUserParams = [
            'is_paginate'   => 0,
            'enterprise_id' => $enterpriseId,
            'conditions'    => [
                ['id', 'in', array_column($data, 'user_id')]
            ]
        ];
        $userList        = (new User())->userList($queryUserParams)['list'];
        $userMap         = array_column($userList, null, 'id');

        // 查询好友信息信息
        $contactMap  = [];
        $unionIdList = array_filter(array_column($userMap, 'union_id'));
        if (!empty($unionIdList)) {
            MemberUserRelation::query()
                ->where('enterprise_id', $enterpriseId)
                ->whereIn('mem_uid', $memUids)
                ->whereIn('union_id', $unionIdList)
                ->get()
                ->each(function (MemberUserRelation $v) use (&$contactMap) {
                    $contactMap[$v['union_id']] = $v->toArray();
                });
        }

        // 拼装返回结果
        $customerList = [];
        foreach ($data as $bindingInfo) {
            $userInfo                     = $userMap[$bindingInfo['user_id']] ?? null;
            $customerInfo                 = [];
            $customerInfo['open_id']      = $userInfo['open_id'] ?? '';
            $customerInfo['union_id']     = $userInfo['union_id'] ?? '';
            $customerInfo['nick_name']    = $userInfo['nick_name'] ?? '';
            $customerInfo['name']         = $bindingInfo['name'];
            $customerInfo['worker_no']    = $bindingInfo['worker_no'];
            $customerInfo['binding_time'] = date("Y-m-d H:i:s", $bindingInfo['created_at']);
            $customerInfo['whether']      = self::ADD_CONTACT_NO;

            $contactInfo = null;
            if (!empty($userInfo['union_id']) && !empty($contactMap[$userInfo['union_id']])) {
                $contactInfo = $contactMap[$userInfo['union_id']];
                if (MemberUserRelation::DEL_DEFAULT == $contactInfo['delete_type']) {
                    $customerInfo['whether'] = self::ADD_CONTACT_YES;
                }
                if (empty($customerInfo['nick_name'])) {
                    $customerInfo['nick_name'] = $contactInfo['user_remark_name'] ?? $contactInfo['ec_name'];
                }
            }
            $customerList[] = $customerInfo;
        }
        return [
            'total' => $count,
            'data'  => $customerList
        ];
    }

}