<?php
/**
 * 内购部门业绩订单列表导出
 * User: liurui
 * Date: 2018/11/26
 * Time: 14:10
 */

namespace App\Http\Controllers\V1\StaffOrder;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StaffOrder;
use Illuminate\Http\Request;

class ExportAchievementDpListController extends BaseController
{

    public function index(Request $request)
    {
        $input = $request->only(['search', 'tag', 'route_url']);
        $input['search']['enterprise_id'] = $this->enterpriseId;
        $input['search']['qy_domain'] = $this->qyDomain;
        $input['search']['app_identifier'] = $this->appIdentifier;
        $input['param_json'] = $input['search'];
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id'] = $this->enterpriseAdmin['id'];
        $input['operator_name'] = $this->enterpriseAdmin['username'];
        $data = (new StaffOrder())->createExportAchievementDpList($input);

        return $data;
    }
}
