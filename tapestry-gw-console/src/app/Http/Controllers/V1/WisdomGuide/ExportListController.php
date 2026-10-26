<?php

namespace App\Http\Controllers\V1\WisdomGuide;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\WisdomGuide;

class ExportListController extends BaseController
{
    /**
     * 导出导购列表
     */
    public function index()
    {
        $input                             = $this->request->only(['search', 'tag', 'route_url']);
        $input['search']['enterprise_id']  = $this->enterpriseId;
        $input['search']['qy_domain']      = $this->qyDomain;
        $input['search']['app_identifier'] = $this->appIdentifier;
        $input['param_json']               = $input['search'];
        $input['enterprise_id']            = $this->enterpriseId;
        $input['operator_id']              = $this->enterpriseAdmin['id'];
        $input['operator_name']            = $this->enterpriseAdmin['username'];
        $data                              = (new WisdomGuide())->createExportList($input);

        return $data;
    }
}
