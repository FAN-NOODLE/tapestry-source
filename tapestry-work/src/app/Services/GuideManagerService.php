<?php

namespace App\Services;

use App\Common\Constant;
use App\Library\Uc\Member;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\Models\StoreWorker;
use CrsSdk\MicroUser\User;
use CrsSdk\Shop\Store;
use Illuminate\Support\Facades\Log;
use CrsSdk\MicroAuthority\Account;

class GuideManagerService extends BaseService
{

    public function autoSyncData($enterpriseId, $domain, $identifier)
    {

        // 获取全部的用户memUid
        $allMember = $this->getAllMemUids($domain, $identifier);
        // 全部用户UIDs集合
        $allMemUids = array_column($allMember, 'memUid');
        $allMember = array_column($allMember, null, 'memUid');

        // 获取支付导购
        $memberGuideLIst = MemberGuide::query()
            ->where('enterprise_id', $enterpriseId)
            ->get(['id', 'mem_uid'])
            ->toArray();
        $memUids         = $guideIds = [];
        // 删除支付导购
        foreach ($memberGuideLIst as $value) {
            if (!in_array($value['mem_uid'], $allMemUids)) {
                $memUids[]  = $value['mem_uid'];
                $guideIds[] = $value['id'];
            }
        }

        // 如果存在支付即导购需要删除
        if (!empty($guideIds)) {
            MemberGuide::query()
                ->whereIn('id', $guideIds)
                ->delete();
        }

        switchAppHost(['enterprise_id'=>$enterpriseId]);
        // 导购搜索条件
        $search = [
            'enterprise_id' => $enterpriseId,
            'worker_status' => Constant::IN_SERVICE,
        ];

        // 导购列表
        $guideList  = StoreWorkerService::list($search, false);

        $params['enterprise_id'] = $enterpriseId;
        $params['domain']        = $domain;
        $params['operator_id']   = '';

        $updateWorkerInfo = [
            'domain' => $domain,
            'member_list' => []
        ];
        // 处理导购
        foreach ($guideList['list'] as $value) {
            // 如果不在范围内
            if (!in_array($value['mem_uid'], $allMemUids)) {
                $params['id']      = $value['id'];
                $params['mem_uid'] = $value['mem_uid'];
                $params['operator_id'] = 0;
                StoreWorkerService::deleteGuide($params);
                $memUids[] = $value['mem_uid'];

                continue;
            }

            if ($value['name'] != $allMember[$value['mem_uid']]['memUsername'] || (!empty($value['mobile']) && $value['mobile'] != $allMember[$value['mem_uid']]['memMobile']) || $value['avatar'] != $allMember[$value['mem_uid']]['memFace'] || $value['qr_code'] == '' || $value['contact_qr_code_id'] == 0) {
                $updateWorkerInfo['member_list'][] = [
                    'id' => $value['id'],
                    'enterprise_id' => $enterpriseId,
                    'store_id' => $value['store_id'],
                    'name' => $allMember[$value['mem_uid']]['memUsername'],
                    'mobile' => $allMember[$value['mem_uid']]['memMobile'],
                    'before_mobile' => $value['mobile'],
                    'avatar' => $allMember[$value['mem_uid']]['memFace'],
                ];
            }
        }
        if ($updateWorkerInfo['member_list']) {
            StoreWorkerService::updateWorkerInfo($updateWorkerInfo);
        }

        // 如果存在支付即导购,以及导购需要解除关系
        if (!empty($memUids)) {
            (new User())->staffShareUnbindingLeave([
                'enterprise_id' => $enterpriseId,
                'mem_uids'      => array_unique($memUids),
                'type'          => Constant::GUIDE_TYPE
            ]);
        }

    }


    /**
     * 获取全部的用户UID
     *
     * @param $domain
     * @param $identifier
     *
     * @return array
     * @throws \Exception
     */
    public function getAllMemUids($domain, $identifier)
    {
//        // 获取所有一级部门
//        $dpList= (new Department())->list($domain,$identifier,['pageNum'=>1,'pageSize'=>500,'dpLevel' => 1]);
//
//        // 获取所有部门ID
//        $dpIds = array_column($dpList['list'], 'dpId');

        $members=[];
//        if(!empty($dpIds)){
            //同步通讯录以后对应部门下的新员工
            $members = MemberService::memberList($domain, $identifier, [
//                'dpIdList'               => $dpIds,
//                'departmentChildrenFlag' => 1,
                'memSubscribeStatus'     => 1,
            ]);
//        }

        // 无部门人员列表
        $noDpMember = MemberService::noDpMemberList($domain, $identifier, ['memSubscribeStatus' => 1]);

        $allMember = array_merge($noDpMember, $members);

        return $allMember;
    }

    /**
     * getGuideListForStatistic
     * @desc 根据筛选条件获取数据总览所有导购列表
     * @param $params
     * @param $qyDomain
     * @param $mobile
     * @param $adminId
     * @param $enterpriseId
     * @param $storeId
     *
     * @return array
     * @date 2020-08-17
     * @throws \Exception
     */
    public static function getGuideListForStatistic($params, $mobile, $qyDomain, $enterpriseId, $adminId, $storeId)
    {
        $guideInfo = StoreWorkerService::detail([
            'enterprise_id' => $enterpriseId,
            'store_id' => $storeId,
            'mobile' => $mobile,
            'worker_status' => 1
        ]);
        if (empty($guideInfo)) {
            return [
                'guide_info' => [],
                'guide_list' => [],
            ];
        }
        // 店员/社群导购/企业导购 只查询个人名下数据  1、门店导购，2、社群导购，3、企业导购, 4.企业管理
        if (in_array($guideInfo['type'], [Constant::STORE_WORKER_TYPE_COMMUNITY, Constant::STORE_WORKER_TYPE_ENTERPRISE])
            || ($guideInfo['type'] == Constant::STORE_WORKER_TYPE_STORE && $guideInfo['is_boss'] == Constant::NO)) {
            return [
                'guide_info' => $guideInfo,
                'guide_list' => [$guideInfo],
            ];
        }

        // 店长和企业管理身份下，如果选择了导购，则直接返回选择的导购列表
        if (isset($params['worker_ids']) && !empty($params['worker_ids'])) {
            $guideList = StoreWorkerService::list(['ids' => $params['worker_ids'], 'enterprise_id' => $enterpriseId], false);
            return [
                'guide_info' => $guideInfo,
                'guide_list' => $guideList['list'] ?? [],
            ];
        }

        // 店长身份，查询店长所在门店下所有导购列表
        if ($guideInfo['type'] == Constant::STORE_WORKER_TYPE_STORE && $guideInfo['is_boss'] == Constant::YES) {
            $guideList = StoreWorkerService::list(['store_id'=> $guideInfo['store_id'],'enterprise_id' => $enterpriseId], false);
            return [
                'guide_info' => $guideInfo,
                'guide_list' => $guideList['list'] ?? [],
            ];
        }

        // 企业管理者，选择了门店，查询所选门店下所有导购列表
        if (isset($params['store_ids']) && !empty($params['store_ids'])) {
            $guideList = StoreWorkerService::list(['store_ids' => $params['store_ids'], 'enterprise_id' => $enterpriseId], false);
            return [
                'guide_info' => $guideInfo,
                'guide_list' => $guideList['list'] ?? [],
            ];
        }

        // 企业管理者，选择了部门，查询所选部门下所有导购列表
        if (isset($params['department_ids']) && !empty($params['department_ids'])) {
            // todo 递归查询
            $result = (new Member())->ignorePermissionsList($qyDomain, '', ['dpIdList' => $params['department_ids']]);
            $memUidList = array_column($result['list'], 'memUid');
            $guideList = StoreWorkerService::list(['enterprise_id' => $enterpriseId, 'mem_uids' => $memUidList], false);
            return [
                'guide_info' => $guideInfo,
                'guide_list' => $guideList['list'] ?? [],
            ];
        }

        // 企业管理者，没有选择部门或者门店，查询所管辖的所有部门或者门店下导购列表
        if ($guideInfo['type'] == Constant::STORE_WORKER_TYPE_MANAGER) {
            $guideList = self::getManagerGuideList($qyDomain, $enterpriseId, $adminId);
            return [
                'guide_info' => $guideInfo,
                'guide_list' => $guideList,
            ];
        }
    }

    /**
     * getGuideListForRank
     * @desc 根据筛选条件获取排行榜所有导购列表
     * @param $mobile
     * @param $qyDomain
     * @param $enterpriseId
     * @param $adminId
     * @param $storeId
     *
     * @return array
     * @date 2020-08-17
     * @throws \Exception
     */
    public static function getGuideListForRank($mobile, $qyDomain, $enterpriseId, $adminId, $storeId)
    {
        $guideInfo = StoreWorkerService::detail([
            'enterprise_id' => $enterpriseId,
            'store_id' => $storeId,
            'mobile' => $mobile,
            'worker_status' => StoreWorker::IN_SERVICE
        ]);
        $guideList = [];
        if (empty($guideInfo)) {
            $guideList = [];
        } else {
            switch ($guideInfo['type']) {
                //1、门店导购
                case Constant::STORE_WORKER_TYPE_STORE:
                    $guideList = (new Store())->areaWorkers([
                        'enterprise_id' => $enterpriseId,
                        'store_id' => $storeId,
                        'is_page' => false,
                    ]);
                    $guideList = $guideList['list'] ?? [];
                    break;
                //2、社群导购
                case Constant::STORE_WORKER_TYPE_COMMUNITY:
                //3、企业导购
                case Constant::STORE_WORKER_TYPE_ENTERPRISE:
                    $result = StoreWorkerService::list([
                        'enterprise_id' => $enterpriseId,
                        'type_list' => [Constant::STORE_WORKER_TYPE_COMMUNITY, Constant::STORE_WORKER_TYPE_ENTERPRISE],
                        'worker_status' => StoreWorker::IN_SERVICE
                    ], false);
                    $guideList = $result['list'] ?? [];
                    break;
                //4.企业管理
                case Constant::STORE_WORKER_TYPE_MANAGER:
                    // 获取 3、企业导购  4.企业管理 列表
                    $result = StoreWorkerService::list([
                        'enterprise_id' => $enterpriseId,
                        'type_list' => [Constant::STORE_WORKER_TYPE_ENTERPRISE, Constant::STORE_WORKER_TYPE_MANAGER],
                        'worker_status' => StoreWorker::IN_SERVICE
                    ], false);
                    $guideList = $result['list'] ?? [];

                    // $guideList = self::getManagerGuideList($qyDomain, $enterpriseId, $adminId);
                    break;
            }
        }
        return [
            'guide_info' => $guideInfo,
            'guide_list' => $guideList,
        ];
    }

    /**
     * getManagerGuideList
     * @desc 获取企业管理者下属的所有导购列表
     * @param $qyDomain
     * @param $enterpriseId
     * @param $adminId
     *
     * @return array
     * @date 2020-08-17
     * @throws \Exception
     */
    public static function getManagerGuideList($qyDomain, $enterpriseId, $adminId)
    {
        $dataRange = (new Account())->dataRangeList(['admin_id' => $adminId]);
        if (empty($dataRange)) {
            return [];
        }
        $storeIdList = [];
        $departmentList = [];
        foreach ($dataRange as $item) {
            switch ($item['type']) {
                case Constant::DATA_RANGE_STORE:
                    $storeIdList[] = $item['object_id'];
                    break;
                case Constant::DATA_RANGE_DEPARTMENT:
                    $departmentList[] = $item['object_id'];
                    break;
            }
        }
        $guideList = [];
        if (!empty($storeIdList)) {
            $result = StoreWorkerService::list(['enterprise_id' => $enterpriseId, 'store_ids' => $storeIdList, 'worker_status' => StoreWorker::IN_SERVICE], false);
            $guideList = $result['list'] ?? [];
        }
        if (!empty($departmentList)) {
            $memberList = self::getDepartmentMemberList($qyDomain, $departmentList);
            $memUidList = array_column($memberList, 'memUid');
            $result = StoreWorkerService::list(['enterprise_id' => $enterpriseId, 'mem_uids' => $memUidList, 'worker_status' => StoreWorker::IN_SERVICE], false);
            $guideList = array_merge($guideList, $result['list'] ?? []);
        }
        return $guideList;
    }

    /**
     * getDepartmentMemberList
     * @desc 循环查询获取部门成员
     * @param $qyDomain
     * @param $departmentList
     *
     * @return array
     * @date 2020-08-26
     * @throws \Exception
     */
    public static function getDepartmentMemberList($qyDomain, $departmentList)
    {
        $isTrue  = true;
        $pageNum = 1;
        $sizeNum = 50;
        $memberList = [];
        do {
            Log::info('ignorePermissionsList');
            $result = (new Member())->ignorePermissionsList($qyDomain, '', [
                'dpIdList' => $departmentList,
                'pageNum'  => $pageNum,
                'pageSize' => $sizeNum,
            ]);

            $memberList = array_merge($memberList, $result['list']);

            if ($result['total'] < $pageNum * $sizeNum) {
                $isTrue = false;
            } else {
                $pageNum = $pageNum + 1;
            }
        } while ($isTrue);
        Log::info('getDepartmentMemberList', $memberList);
        return $memberList;
    }

    /**
     * guideRankOfCustomerData
     * @param $enterpriseId
     * @param $qyDomain
     * @param $guideList
     * @param $type
     * @param $startTime
     * @param $endTime
     *
     * @return mixed
     * @throws \Exception
     */
    public static function guideRankOfCustomerData($enterpriseId, $qyDomain, $guideList, $type, $startTime, $endTime)
    {
        //返回数组
        $result = [];
        //导购二维码id
        $codeIdList = array_column($guideList, 'contact_qr_code_id');
        $codeIdList = array_filter($codeIdList);
        $customerService = new CustomerService();
        $userService = new \CrsSdk\Shop\User();
        $query = MemberUserRelation::query();
        $query = $query->where('member_user_relation.enterprise_id', $enterpriseId)
            ->whereIn('source_type', $codeIdList)
            ->where('delete_type', 0);
        //当排序为浏览顾客总数或成交顾客总数时，按数据组装成导购二维码ID=》【顾客union_id数组】
        if (in_array($type, [Constant::GUIDE_CUSTOMER_RANK_ACCESS, Constant::GUIDE_CUSTOMER_RANK_DEAL])) {
            //查询时间区间的顾客union_id
            $memberUserRelationList = $query->whereBetween('relation_at', [$startTime, $endTime])
                ->selectRaw('source_type,union_id')
                ->get()
                ->toArray();
            foreach ($memberUserRelationList as $key => $val) {
                $source_type = $val['source_type'];
                if (!isset($result[$source_type])) {
                    $result[$source_type] = [
                        'source_type' => $val['source_type'],
                        'union_id'    => []
                    ];
                }
                $result[$source_type]['union_id'] = array_merge([$val['union_id']],$result[$source_type]['union_id']);
            }
            //去掉key
            sort($result);
        }
        switch ($type){
            case Constant::GUIDE_CUSTOMER_RANK_CUSTOMER:
                $result = $query->groupBy('source_type')
                    ->selectRaw('source_type, count(ec_id) as num')
                    ->orderByDesc('num')
                    ->limit(50)
                    ->get()
                    ->toArray();
                break;
            case Constant::GUIDE_CUSTOMER_RANK_NEW_CUSTOMER:
                $result = $query->whereBetween('relation_at', [$startTime, $endTime])
                    ->groupBy('source_type')
                    ->selectRaw('source_type, count(ec_id) as num')
                    ->orderByDesc('num')
                    ->limit(50)
                    ->get()
                    ->toArray();
                break;
            case Constant::GUIDE_CUSTOMER_RANK_ACCESS:
                //按照导购去查询浏览顾客总数
                foreach ($result as &$item) {
                    $params['enterprise_id'] = $enterpriseId;
                    $params['start_time']    = $startTime;
                    $params['end_time']      = $endTime;
                    $params['union_ids']     = $item['union_id'];
                    $accessCount             = $userService->dailyAccessRecords($params)['list'];
                    $item['num']             = !empty($accessCount) ? array_sum(array_column($accessCount,'num')) : 0;
                    unset($item['union_id']);
                }
                $sortField = array_column($result,'num');
                array_multisort($sortField,SORT_DESC,$result);
                $result = array_slice($result,0,50);
                break;
            case Constant::GUIDE_CUSTOMER_RANK_DEAL:
                //按照导购去查询成交顾客总数
                foreach ($result as &$item) {
                    $dealCount  = $customerService->orderOutStatistics($item['union_id'],$startTime,$endTime);
                    $item['num']             = !empty($dealCount) ? count(array_unique($dealCount)) : 0;
                    unset($item['union_id']);
                }
                $sortField = array_column($result,'num');
                array_multisort($sortField,SORT_DESC,$result);
                $result = array_slice($result,0,50);
                break;
        }
        return self::formatCustomerRankData($guideList, $result, $qyDomain);
    }

    /**
     * formatCustomerRankData
     * @desc 格式化顾客数据排行榜
     * @param $guideList
     * @param $rank
     * @param $qyDomain
     *
     * @return mixed
     * @date 2020-08-19
     * @throws \Exception
     */
    private static function formatCustomerRankData($guideList, $rank, $qyDomain)
    {
        // 补充排行榜列表
        $rank = self::setRankList($guideList, $rank);
        $guideList = array_column($guideList, null, 'id');
        $storeIdList = [];
        $departmentMemUidList = [];
        foreach ($rank as &$item) {
            $guideInfo = $guideList[$item['worker_id']];
            $item['name'] = $guideInfo['name'];
            $item['avatar'] = $guideInfo['avatar'];
            $item['store_id'] = $guideInfo['store_id'];
            $item['mem_uid'] = $guideInfo['mem_uid'];
            if ($guideInfo['type'] == Constant::STORE_WORKER_TYPE_STORE) {
                $storeIdList[] = $item['store_id'];
            } else {
                $departmentMemUidList[] = $item['mem_uid'];
            }
        }
        $departmentList = [];
        $storeList = [];
        if (!empty($departmentMemUidList)) {
            $departmentList = self::getGuideDepartment($departmentMemUidList, $qyDomain);
            $departmentList = array_column($departmentList, null, 'memUid');
        }
        if (!empty($storeIdList)) {
            $storeList = (new Store())->listByIds(['store_ids' => $storeIdList]);
            $storeList = array_column($storeList, null, 'id');
        }
        array_walk($rank, function (&$item) use ($departmentList, $storeList, $guideList) {
            if ($guideList[$item['worker_id']]['type'] == Constant::STORE_WORKER_TYPE_STORE) {
                $item['store_name'] = $storeList[$item['store_id']]['short_name'] ?: $storeList[$item['store_id']]['name'];
            } else {
                $item['store_name'] = $departmentList[$item['mem_uid']]['dpName'] ?? '';
            }
        });
        return $rank;
    }

    /**
     * setRankList
     * @desc 排行榜数量不到50个时，补充没有数据的导购到排行榜列表中
     * @param $guideList
     * @param $rankList
     *
     * @return array
     * @date 2020-08-26
     */
    private static function setRankList($guideList, $rankList)
    {
        $guideIdDict = array_column($guideList, 'id', 'contact_qr_code_id');
        foreach ($rankList as &$value) {
            $value['worker_id'] = $guideIdDict[$value['source_type']];
        }
        $rankCount = count($rankList);
        if ($rankCount >= Constant::GUIDE_RANK_NUM) {
            return $rankList;
        }
        $guideCount = count($guideList);
        if ($guideCount == $rankCount) {
            return $rankList;
        }
        $rankSourceTypeList = array_column($rankList, 'source_type');
        foreach ($guideList as $key => $guide) {
            if (in_array($guide['contact_qr_code_id'], $rankSourceTypeList)) {
                unset($guideList[$key]);
            }
        }
        $guideList = array_slice($guideList, 0, Constant::GUIDE_RANK_NUM - $rankCount);
        foreach ($guideList as $key => $guide) {
            $rankList[] = [
                'source_type' => $guide['contact_qr_code_id'],
                'worker_id' => $guide['id'],
                'num' => 0
            ];
        }
        return $rankList;
    }

    /**
     * getGuideDepartment
     * @desc 获取导购部门信息
     * @param $memUidList
     * @param $qyDomain
     *
     * @return mixed
     * @date 2020-08-19
     * @throws \Exception
     */
    public static function getGuideDepartment($memUidList, $qyDomain)
    {
        $params = [
            'memUids'  => $memUidList,
            'pageNum'  => 1,
            'pageSize' => count($memUidList)
        ];
        $guideInfoList = (new Member())->list($qyDomain, '', $params);
        $guideInfoList = $guideInfoList['list'];
        array_walk($guideInfoList, function (&$item) {
            $item = [
                'memUid' => $item['memUid'],
                'dpId' => $item['dpName'][0]['dpId'] ?? 0,
                'dpName' => $item['dpName'][0]['dpName'] ?? ''
            ];
        });
        return $guideInfoList;
    }

}
