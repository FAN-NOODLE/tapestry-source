<?php
/**
 * 获取企业微信红包对账记录员工列表
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;

class ReconciliationMemberController extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params                   = $request->all();
        $params['enterprise_id']  = $this->enterpriseId;
        $params['qy_domain']      = $this->qyDomain;
        $params['app_identifier'] = $this->appIdentifier;

        $data = (new EnterprisePayment())->reconciliationMember($params);

        return success([
            'list' => $data['data']
        ]);
    }
}
