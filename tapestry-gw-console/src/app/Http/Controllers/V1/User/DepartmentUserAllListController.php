<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use CrsSdk\Work\Department;
use App\Http\Controllers\V1\BaseController;

class DepartmentUserAllListController extends BaseController
{
    /**
     * 员工部门分层列表
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        // 调用sdk获取门店dptId
        $dptIdData = (new Department())->storeDptId([
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier
        ]);

        // 部门门店类型id
        $dptId = $dptIdData['data']['dptId'] ?? '';

        // 所有部门
        $dpList = (new Department())->departmentList([
            'params'         => ['dptId' => $dptId, 'pageSize' => 1000],
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier
        ])['list'];

        $storeList = [];
        array_walk($dpList, function ($value) use (&$storeList, $dptId) {
            if ($value['dptId'] == $dptId) {
                $storeList[$value['dpId']] = [
                    'dpId'     => $value['dpId'],
                    'dpName'   => $value['dpName'],
                    'children' => []
                ];
            }
        });
        if (empty($storeList)) {
            return success(['data' => []]);
        }

        // 获取门店内所有人
        $storeIds = array_column($storeList, 'dpId');
        $user     = (new Member())->memberList([
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => ['dpIdList' => $storeIds, 'pageSize' => 1000],
        ])['list'];

        foreach ($user as $value) {
            if (
                isset($value['dpName'])
                && current($value['dpName'])
                && isset($storeList[current($value['dpName'])['dpId']])
            ) {
                $storeList[current($value['dpName'])['dpId']]['children'][] = [
                    'memUid'               => $value['memUid'],
                    'memUsername'          => $value['memUsername'],
                    'memUsernameAndDpName' => $value['memUsername'] . '@' . $storeList[current($value['dpName'])['dpId']]['dpName']
                ];
            }
        }

        return success(['list' => array_values($storeList)]);
    }
}
