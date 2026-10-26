<?php
/**
 * 创建导出对账数据列表
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\EnterprisePayment;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\EnterprisePayment;
use Illuminate\Http\Request;

class CreateExportBudgetList extends BaseController
{

    /**
     * @param Request $request
     *
     * @return mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $input                             = $request->only(['search', 'tag', 'route_url']);
        $input['search']['enterprise_id']  = $this->enterpriseId;
        $input['search']['qy_domain']      = $this->qyDomain;
        $input['search']['app_identifier'] = $this->appIdentifier;
        $input['param_json']               = $input['search'];
        $input['enterprise_id']            = $this->enterpriseId;
        $input['operator_id']              = $this->enterpriseAdmin['id'];
        $input['operator_name']            = $this->enterpriseAdmin['username'];

        (new EnterprisePayment())->createExportBudgetList($input);

        return success();
    }
}
