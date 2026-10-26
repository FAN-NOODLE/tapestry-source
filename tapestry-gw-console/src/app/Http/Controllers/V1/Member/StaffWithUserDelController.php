<?php

namespace App\Http\Controllers\V1\Member;

use CrsSdk\Work\Customer;
use App\Http\Controllers\V1\BaseController;

class StaffWithUserDelController extends BaseController
{
    /**
     * 客户员工互删
     * @return mixed
     * @throws \Exception
     */
    public function staffWithUserDel()
    {
        $params                   = $this->request->all();
        $params['enterprise_id']  = $this->enterpriseId;
        $params['qy_domain']      = $this->qyDomain;
        $params['app_identifier'] = $this->appIdentifier;

        $list = (new Customer())->staffUserDel($params);

        return success([
            'list'  => $list['data'],
            'page'  => $list['meta']['page'],
            'size'  => $list['meta']['size'],
            'total' => $list['meta']['total'],
        ]);
    }
}