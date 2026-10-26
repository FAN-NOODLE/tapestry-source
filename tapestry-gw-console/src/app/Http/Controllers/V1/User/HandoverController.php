<?php

namespace App\Http\Controllers\V1\User;

use CrsSdk\Work\Member;
use CrsSdk\Work\Customer;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class HandoverController extends BaseController
{
    /**
     * 交接客户
     * @param Request $request
     * @return JsonResponse
     */
    public function index(Request $request)
    {
        // 离职员工id
        $memUid = $request->post('mem_uid', '');
        if (!$memUid) {
            error(320161);
        }

        // 新导购id
        $newMemUid = $request->post('new_mem_uid', '');
        if (!$newMemUid) {
            error(320162);
        }

        // 待分配外部联系人id集合
        $ecUserIds = $request->post('ec_userids', []);
        if (!$ecUserIds) {
            error(320053);
        }

        // 获取新旧交接人名称
        $paramsOld = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => ['memUid' => $memUid]
        ];
        $paramsNew = [
            'qy_domain'      => $this->qyDomain,
            'app_identifier' => $this->appIdentifier,
            'params'         => ['memUid' => $newMemUid]
        ];

        $detailOld    = (new Member())->detail($paramsOld);
        $detailNew    = (new Member())->detail($paramsNew);

        // 调用sdk，交接客户
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'qy_domain'     => $this->qyDomain,
            'mem_uid'       => $memUid,
            'new_mem_uid'   => $newMemUid,
            'ec_userids'    => $ecUserIds
        ];
        (new Customer())->handover($params);

        // 写入记录日志
        $memUsernameOld = $detailOld['memUsername'] ? $detailOld['memUsername'] : '';
        $memUsernameNew = $detailNew['memUsername'] ? $detailNew['memUsername'] : '';
        $this->writeOperatorLog($memUsernameOld .' 交接员工给 '.$memUsernameNew);

        return success();
    }
}
