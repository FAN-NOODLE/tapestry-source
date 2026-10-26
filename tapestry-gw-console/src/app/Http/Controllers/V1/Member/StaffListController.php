<?php

namespace App\Http\Controllers\V1\Member;

use CrsSdk\Work\Customer;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class StaffListController extends BaseController
{
    /**
     * 员工列表
     * @return JsonResponse
     * @throws \Exception
     */
    public function staffList()
    {
        $params                   = $this->request->all();
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;

        $list = (new Customer())->staffList($params);

        return success([
            'list' => $list['data']
        ]);
    }
}