<?php

namespace App\Http\Controllers\V1\Member;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Customer;

class MemberGiveTagController extends BaseController
{
    /**
     * 单独给用户打标签
     * @return mixed
     * @throws \Exception
     */
    public function memberGiveTag()
    {
        $params                   = $this->request->all();
        $params['qy_domain']      = $this->qyDomain;
        $params['enterprise_id']  = $this->enterpriseId;
        $params['app_identifier'] = $this->appIdentifier;

        (new Customer())->giveTag($params);

        return success();
    }
}