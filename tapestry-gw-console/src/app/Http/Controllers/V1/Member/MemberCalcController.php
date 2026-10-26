<?php

namespace App\Http\Controllers\V1\Member;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\Order;
use CrsSdk\Work\Customer;

class MemberCalcController extends BaseController
{
    /**
     * 客户统计
     */
    public function memberCalc()
    {
        $params                   = $this->request->all();
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;

        $data = (new Customer)->customerCalc($params);

        return success($data['data']);
    }
}