<?php

namespace App\Services\StoreWorker\AddGuideInMember;

use App\Common\Constant;
use App\Services\EnterpriseService;
use App\Library\Uc\{Department, Member};
use App\Models\{StoreDepartment, StoreWorker};
use App\Services\StoreWorkerService;
use CrsSdk\Boss\EnterpriseConfig;
use Exception;
use Illuminate\Support\Facades\Log;

class AddGuideService
{
    public static function addGuideNotInStoreWorker($epId, $qyDomain)
    {
        try {
            Log::info("自动添加导购-开始同步", [$epId, $qyDomain]);
            // 获取企业配置要同步的部门
            $enterpriseConfig = EnterpriseService::getEnterpriseConfig($epId, StoreDepartment::STORE_PARENT_DEPT_ID_KEY);

            Log::info("自动添加导购-获取企业配置要同步的部门", [$enterpriseConfig]);

            if (!$enterpriseConfig['data']) {
                return;
            }

            $enterpriseConfigData = current($enterpriseConfig['data']);
            $thirdId              = $enterpriseConfigData['value'];

            // 获取要同步部门的Uc部门详情
            $department = (new Department())->list($qyDomain, '', ['thirds' => [$thirdId]]);
            Log::info("自动添加导购-获取要同步部门的Uc部门详情", [$department]);
            $departmentDetails = current($department['list']);

            // 获取要同步部门下的所有部门Id集合
            $departmentList = (new Department())->list($qyDomain, '', [
                'dpId'                   => $departmentDetails['dpId'],
                'departmentChildrenFlag' => 1
            ]);
            Log::info("自动添加导购-获取要同步部门下的所有部门Id集合", [$departmentList]);
            $departments   = $departmentList['list'];
            $departmentIds = array_column($departments, 'dpId');

            // 查询已经是导购的人员Id集合
            $guideMemUids = StoreWorker::query()
                ->where('enterprise_id', $epId)
                ->pluck("mem_uid")
                ->toArray();

            $appIdentifier = null;
            $pageNum       = 1;

            while (true) {
                $params = [
                    'memSubscribeStatus' => 1,
                    'dpIdList'           => $departmentIds,
                    'pageNum'            => $pageNum,
                    'pageSize'           => 500,
                ];
                $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);
                Log::info("自动添加导购-忽略权限人员列表", [$result]);
                if (!$result['list']) {
                    break;
                }

                $membersData = $result['list'];

                $addMembersData = [];
                foreach ($membersData as $info) {
                    if (!in_array($info['memUid'], $guideMemUids)) {
                        $addMembersData[] =
                            [
                                "mem_uid"   => $info['memUid'],
                                "user_name" => $info["memUsername"],
                                "mobile"    => $info["memMobile"],
                                "user_id"   => $info["memUserid"],
                                "avatar"    => $info["memFace"],
                            ];
                    }
                }
                if ($addMembersData) {
                    $guideData = [
                        'enterprise_id' => $epId,
                        'member'        => $addMembersData,
                        'type'          => Constant::STORE_WORKER_TYPE_ENTERPRISE
                    ];

                    StoreWorkerService::addGuide($guideData);
                }
                $pageNum++;
            }
        } catch (Exception $exception) {
            Log::error("自动添加导购-添加导购出现错误", [$exception]);
        }
    }
}