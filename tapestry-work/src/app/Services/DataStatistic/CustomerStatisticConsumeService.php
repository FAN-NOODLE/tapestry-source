<?php


namespace App\Services\DataStatistic;


use App\Common\DateUtils;
use App\Common\MessageConstant;
use App\Jobs\Message\ExternalContactMessage;
use App\Jobs\Message\ShareLinkMessage;
use App\Library\DataLake\UserRequestService;
use App\Models\MemberUserRelation;
use App\Models\StaffShare;
use App\Models\StoreWorker;
use App\Services\BaseService;
use App\Services\StaffShareService;
use App\Services\UserRemoteService;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Throwable;

class CustomerStatisticConsumeService extends BaseService
{

    private $workerStatisticService;

    private $memberStatisticService;

    private $dataLakeUserService;

    /**
     * CustomerStatisticConsumeService constructor.
     */
    public function __construct()
    {
        $this->workerStatisticService = new DailyWorkerCustomerStatisticService();
        $this->memberStatisticService = new DailyMemberCustomerStatisticService();
        $this->dataLakeUserService    = new UserRequestService();
    }

    /**
     * 外部联系人变更-更新导购顾客数据
     * @param ExternalContactMessage $message
     * @param array $workerList
     * @throws Exception
     */
    public function updateCustomerStatisticByExternalContactMessage(ExternalContactMessage $message, array $workerList)
    {
        Log::info("updateCustomerStatisticByExternalContactMessage start");

        $thenStatTime     = DateUtils::timeToStatTime($message->getCreatedAt()); // 好友添加时间当日0点时间戳
        $todayStatTime    = DateUtils::todayStatTime(); // 今日0点时间戳
        $tomorrowStatTime = $thenStatTime + DateUtils::DAY_SECONDS; // 明日0点时间戳
        $isCreateMessage  = MessageConstant::EXTERNAL_CONTACT_CREATE == $message->getType();
        $statNum          = $isCreateMessage ? 1 : -1;

        // 更新导购绑定顾客统计数据
        DB::connection()->beginTransaction();
        try {
            // 更新员工维度统计数据
            $memberSetData           = [
                'cumulative_customer_num' => $statNum,
                'today_customer_num'      => $statNum
            ];
            $memberCumulativeSetDate = [
                'cumulative_customer_num' => $statNum
            ];

            // 调取DL接口获取会员等级信息
            $this->dataLakeUserService->setEnterpriseId($message->getEnterpriseId());
            try {
                $grade = $this->dataLakeUserService->queryGradeInfo('unionId', $message->getUnionId());
            } catch (Throwable $e) {
                Log::info('当前顾客不是会员');
                $grade = [];
            }

            $isMember = isset($grade['grade']) && $grade['grade'] ? true : false;

            if ($isMember) {
                $memberSetData['cumulative_member_num']           = $statNum;
                $memberSetData['today_member_num']                = $statNum;
                $memberCumulativeSetDate['cumulative_member_num'] = $statNum;
            } else {
                $memberSetData['cumulative_non_member_num']           = $statNum;
                $memberSetData['today_non_member_num']                = $statNum;
                $memberCumulativeSetDate['cumulative_non_member_num'] = $statNum;
            }

            // 更新会员相关信息
            $this->updateMemberUserRelationData($grade, $message->getUnionId());

            // 更新员工当日的累计客户数、当日客户数（添加好友时间）
            $incrRes = $this->memberStatisticService->incrMemberStatisticField($message->getEnterpriseId(),
                $message->getMemUid(), $thenStatTime, $memberSetData, $isCreateMessage);

            // 如果添加好友时间是历史日期，更新这段时间内的累计客户数
            if ($incrRes && $todayStatTime > $thenStatTime) {
                $this->memberStatisticService->incrMemberStatisticFieldByStatTimeRange($message->getEnterpriseId(),
                    $message->getMemUid(), $tomorrowStatTime, $todayStatTime, $memberCumulativeSetDate);
            }

            // 更新导购维度统计数据
            foreach ($workerList as $worker) {
                if (StoreWorker::IN_SERVICE != $worker['worker_status']) {
                    continue;
                }
                $isContactCardWay = !empty($worker['contact_qr_code_id'])
                    && $worker['contact_qr_code_id'] == $message->getSourceType();

                $workerSetData = $workerCumulativeSetDate = [];
                if ($isContactCardWay) {
                    // 名片添加数
                    $workerSetData['cumulative_contact_card_way_customer_num']           = $statNum;
                    $workerSetData['today_contact_card_way_customer_num']                = $statNum;
                    $workerCumulativeSetDate['cumulative_contact_card_way_customer_num'] = $statNum;
                } else {
                    // 其他方式添加数
                    $workerSetData['cumulative_other_way_customer_num']           = $statNum;
                    $workerSetData['today_other_way_customer_num']                = $statNum;
                    $workerCumulativeSetDate['cumulative_other_way_customer_num'] = $statNum;
                }

                $incrRes = $this->workerStatisticService->incrWorkerStatisticField($message->getEnterpriseId(),
                    $worker['id'], $worker['mem_uid'], $thenStatTime, $workerSetData, $isCreateMessage);

                if ($incrRes && $todayStatTime > $thenStatTime) {
                    $this->workerStatisticService->incrWorkerStatisticFieldByStatTimeRange($message->getEnterpriseId(),
                        $worker['id'], $tomorrowStatTime, $todayStatTime, $workerCumulativeSetDate);
                }
            }

            DB::commit();
            Log::info("updateCustomerStatisticByExternalContactMessage succeed");
        } catch (Exception $e) {
            DB::rollBack();
            Log::info("updateCustomerStatisticByExternalContactMessage exception", [
                'code'      => $e->getCode(),
                'message'   => $e->getMessage(),
                'exception' => $e->getTraceAsString()
            ]);
            throw $e;
        }
    }

    /**
     * 更新会员相关信息
     * @param $grade
     * @param $unionId
     * @return bool|int
     * @Date: 2021/4/9
     * @Time: 9:54
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public function updateMemberUserRelationData($grade, $unionId)
    {
        try {
            $memberInfo = $this->dataLakeUserService->queryMemberInfo('unionId', $unionId);
            $point      = $this->dataLakeUserService->getCurrentPoint('unionId', $unionId);

        } catch (Throwable $e) {
            // 非会员
            return false;
        }

        $update['union_code'] = $memberInfo['union_code'] ?? '';
        $update['mobile']     = $memberInfo['mobile'];
        $update['birthday']   = $memberInfo['birthday'];
        $update['point']      = $point['points'][0]['evalidPoints'] ?? 0;
        $update['grade']      = $grade['grade'];

        $result = MemberUserRelation::query()
            ->where(['union_id' => $unionId])
            ->update($update);

        if ($result) {
            Log::info('更新成功');
        }

        return $result;
    }

    /**
     * 外部联系人变更-更新导购绑客数据
     * @param ExternalContactMessage $message
     * @param array $workerList
     */
    public function updateBindingCustomerStatisticByExternalContactMessage(ExternalContactMessage $message,
                                                                           array                  $workerList)
    {
        Log::info("updateBindingCustomerStatisticByExternalContactMessage start");
        $thenStatTime = DateUtils::timeToStatTime($message->getCreatedAt());
        if (MessageConstant::EXTERNAL_CONTACT_CREATE == $message->getType()) {
            // 添加好友关系，判断客户和导购是否绑定了业绩关系，并且绑定时间是今天，如果是，绑客数+1
            $workerMap = array_column($workerList, null, 'id');
            $userInfo  = UserRemoteService::getUserInfoByUnionId($message->getEnterpriseId(), $message->getUnionId());
            $binding   = !empty($userInfo['union_id']) ?
                StaffShareService::getUserBinding($message->getEnterpriseId(), $userInfo['id']) : [];

            // 绑定顾客数据
            if (!empty($userInfo) && !empty($binding) && isset($workerMap[$binding['worker_id']])
                && $binding['created_at'] > $thenStatTime) {
                $bindingWorker = $workerMap[$binding['worker_id']];

                $this->workerStatisticService->incrWorkerStatisticField(
                    $message->getEnterpriseId(),
                    $bindingWorker['id'],
                    $bindingWorker['mem_uid'],
                    $thenStatTime,
                    [
                        'cumulative_binding_customer_num' => 1,
                        'today_binding_customer_num'      => 1
                    ]
                );
            } else {
                Log::info("updateBindingCustomerStatisticByExternalContactMessage conditions are not met", [
                    'userInfo' => $userInfo,
                    'binding'  => $binding
                ]);
            }
        } else {
            // 删除好友关系，刷新导购累计绑客数
            $workerIdList = array_column($workerList, 'id');
            $statNumList  = $this->getMemberCumulativeBindingCustomerNum($message->getEnterpriseId(),
                $message->getMemUid(), $workerIdList);
            Log::info("updateBindingCustomerStatisticByExternalContactMessage statNumList", $statNumList);

            foreach ($workerList as $worker) {
                $this->workerStatisticService->updateWorkerStatisticField(
                    $message->getEnterpriseId(),
                    $worker['mem_uid'],
                    $worker['id'],
                    $thenStatTime,
                    [
                        'cumulative_binding_customer_num' => $statNumList[$worker['id']] ?? 0
                    ],
                    true
                );
            }
        }
        Log::info("updateBindingCustomerStatisticByExternalContactMessage succeed");
    }

    /**
     * 获取员工下所有导购累计绑客数（导购和客户添加了好友并且绑定了业绩）
     * @param $enterpriseId
     * @param $memUid
     * @param $workerIdList
     * @return array
     */
    private function getMemberCumulativeBindingCustomerNum($enterpriseId, $memUid, $workerIdList)
    {
        $result = [];
        foreach ($workerIdList as $workerId) {
            $result[$workerId] = 0;
        }

        // 获取员工所有外部联系人的union_id
        $externalContactList        = MemberUserRelation::query()
            ->select('union_id')
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('union_id', '!=', '')
            ->where('delete_type', MemberUserRelation::DEL_DEFAULT)
            ->get()
            ->toArray();
        $externalContactUnionIdList = array_unique(array_filter(array_column($externalContactList, 'union_id')));
        if (empty($externalContactUnionIdList)) {
            return $result;
        }
        unset($externalContactList);

        // 分页查询用户信息
        $limit     = 1000;
        $pageTotal = (int)ceil(count($externalContactUnionIdList) / $limit);
        for ($page = 1; $page <= $pageTotal; $page++) {
            $pageUnionIdList  = array_slice($externalContactUnionIdList, ($page - 1) * $limit, $limit);
            $pageUserInfoList = UserRemoteService::listUserInfoByUnionId($enterpriseId, $pageUnionIdList);
            if (!empty($pageUserInfoList)) {
                $bindingStatisticList = StaffShare::query()
                    ->select(['worker_id', DB::raw('count(*) as cnt')])
                    ->where('enterprise_id', $enterpriseId)
                    ->when(count($workerIdList) == 1, function ($query) use ($workerIdList) {
                        return $query->where('worker_id', $workerIdList[0]);
                    }, function ($query) use ($workerIdList) {
                        return $query->whereIn('worker_id', $workerIdList);
                    })
                    ->whereIn('user_id', array_column($pageUserInfoList, 'id'))
                    ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                    ->groupBy('worker_id')
                    ->get()
                    ->toArray();
                $bindingStatisticMap  = array_column($bindingStatisticList, 'cnt', 'worker_id');
                foreach ($workerIdList as $workerId) {
                    $result[$workerId] += $bindingStatisticMap[$workerId] ?? 0;
                }
            }
            unset($pageUnionIdList, $pageUserInfoList);
        }

        return $result;
    }

    /**
     * 分享链路变更-更新导购绑客数
     * @param ShareLinkMessage $message
     */
    public function updateBindingCustomerStatisticByShareLinkMessage(ShareLinkMessage $message)
    {
        Log::info("updateBindingCustomerStatisticByShareLinkMessage start");
        // 当前绑定的导购，和用户同时添加了好友，今日新增绑客数+1
        if (in_array($message->getType(), [MessageConstant::SHARE_LINK_CREATE, MessageConstant::SHARE_LINK_UPDATE])) {
            $thenStatTime = DateUtils::timeToStatTime($message->getCreatedAt());

            $userInfo = UserRemoteService::getUserInfoById($message->getEnterpriseId(), $message->getUserId());
            if (!empty($userInfo) && !empty($userInfo['union_id'])) {

                // 今日和导购绑定的次数
                $todayBindingCount = StaffShare::query()
                    ->where('enterprise_id', $message->getEnterpriseId())
                    ->where('user_id', $message->getUserId())
                    ->where('worker_id', $message->getCurrWorkerId())
                    ->whereBetween('created_at', [$thenStatTime, $thenStatTime + DateUtils::DAY_SECONDS - 1])
                    ->count();

                // 查询是否和导购添加了好友
                $externalContactExists = MemberUserRelation::query()
                    ->where('enterprise_id', $message->getEnterpriseId())
                    ->where('mem_uid', $message->getCurrWorkerMemUid())
                    ->where('union_id', $userInfo['union_id'])
                    ->where('delete_type', MemberUserRelation::DEL_DEFAULT)
                    ->exists();
                if ($todayBindingCount == 1 && $externalContactExists) {
                    $this->workerStatisticService->incrWorkerStatisticField(
                        $message->getEnterpriseId(),
                        $message->getCurrWorkerId(),
                        $message->getCurrWorkerMemUid(),
                        $thenStatTime,
                        [
                            'cumulative_binding_customer_num' => 1,
                            'today_binding_customer_num'      => 1
                        ]
                    );
                } else {
                    Log::info("updateBindingCustomerStatisticByShareLinkMessage conditions are not met", [
                        'todayBindingCount'     => $todayBindingCount,
                        'externalContactExists' => $externalContactExists
                    ]);
                }
            } else {
                Log::info("updateBindingCustomerStatisticByShareLinkMessage userInfo or unionId is empty", [
                    'userInfo' => $userInfo
                ]);
            }
        }

        // 上次绑定的导购，重新计算累计绑客数
        if (in_array($message->getType(), [MessageConstant::SHARE_LINK_DELETE, MessageConstant::SHARE_LINK_UPDATE])) {
            $thenStatTime = DateUtils::timeToStatTime($message->getUpdatedAt());
            $this->refreshWorkerCumulativeBindingCustomerNum($message->getEnterpriseId(),
                $message->getOldWorkerMemUid(), $message->getOldWorkerId(), $thenStatTime);
        }
        Log::info("updateBindingCustomerStatisticByShareLinkMessage succeed");
    }

    /**
     * 刷新导购累计绑客数
     * @param $enterpriseId
     * @param $memUid
     * @param $workerId
     * @param $statTime
     */
    private function refreshWorkerCumulativeBindingCustomerNum($enterpriseId, $memUid, $workerId, $statTime)
    {
        Log::info("refreshWorkerCumulativeBindingCustomerNum params", func_get_args());
        $statNum = $this->getMemberCumulativeBindingCustomerNum($enterpriseId, $memUid, [$workerId])[$workerId] ?? 0;

        $this->workerStatisticService->updateWorkerStatisticField(
            $enterpriseId,
            $memUid,
            $workerId,
            $statTime,
            [
                'cumulative_binding_customer_num' => $statNum
            ]
        );
    }

    /**
     * 分享链路批量变更-更新导购绑客数
     * @param int $type
     * @param array $messageList
     */
    public function updateBindingCustomerStatisticByShareLinkBatchMessage(int $type, array $messageList)
    {

        Log::info("updateBindingCustomerStatisticByShareLinkBatchMessage start");
        $refreshWorkerList = [];

        if (MessageConstant::SHARE_LINK_BATCH_INHERIT == $type) {
            $userIdList = [];
            foreach ($messageList as $message) {
                $userIdList[] = $message->getUserId();
            }
            foreach ($messageList as $message) {
                $thenStatTime = DateUtils::timeToStatTime($message->getUpdatedAt());

                $userInfoList = UserRemoteService::listUserInfo($message->getEnterpriseId(), $userIdList);
                $unionIdList  = array_unique(array_filter(array_column($userInfoList, 'union_id')));

                // 查询继承的所有用户中添加好友数量
                if (!empty($unionIdList)) {
                    $externalContactNum = MemberUserRelation::query()
                        ->where('enterprise_id', $message->getEnterpriseId())
                        ->where('mem_uid', $message->getCurrWorkerMemUid())
                        ->whereIn('union_id', $unionIdList)
                        ->where('delete_type', MemberUserRelation::DEL_DEFAULT)
                        ->count();
                    Log::info("updateBindingCustomerStatisticByShareLinkBatchMessage externalContactNum: " . $externalContactNum);
                    if ($externalContactNum > 0) {
                        $this->workerStatisticService->incrWorkerStatisticField(
                            $message->getEnterpriseId(),
                            $message->getCurrWorkerId(),
                            $message->getCurrWorkerMemUid(),
                            $thenStatTime,
                            [
                                'cumulative_binding_customer_num' => $externalContactNum,
                                'today_binding_customer_num'      => $externalContactNum
                            ]
                        );
                    }
                } else {
                    Log::info("updateBindingCustomerStatisticByShareLinkBatchMessage unionIdList is empty");
                }

                $refreshWorkerList[$message->getOldWorkerId()] = [
                    'enterprise_id' => $message->getEnterpriseId(),
                    'worker_id'     => $message->getOldWorkerId(),
                    'mem_uid'       => $message->getOldWorkerMemUid(),
                    'stat_time'     => DateUtils::timeToStatTime($message->getUpdatedAt())
                ];
                break;
            }
        } else {
            foreach ($messageList as $message) {
                if (!isset($refreshWorkerList[$message->getOldWorkerId()])) {
                    $refreshWorkerList[$message->getOldWorkerId()] = [
                        'enterprise_id' => $message->getEnterpriseId(),
                        'worker_id'     => $message->getOldWorkerId(),
                        'mem_uid'       => $message->getOldWorkerMemUid(),
                        'stat_time'     => DateUtils::timeToStatTime($message->getUpdatedAt())
                    ];
                }
            }
        }
        Log::info("updateBindingCustomerStatisticByShareLinkBatchMessage refreshWorkerList", $refreshWorkerList);

        // 刷新所有变更导购的累计绑客数
        foreach ($refreshWorkerList as $worker) {
            $this->refreshWorkerCumulativeBindingCustomerNum($worker['enterprise_id'], $worker['mem_uid'],
                $worker['worker_id'], $worker['stat_time']);
        }
        Log::info("updateBindingCustomerStatisticByShareLinkBatchMessage succeed");
    }

}