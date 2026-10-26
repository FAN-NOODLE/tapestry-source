<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class MaterialListController extends BaseController
{
    /**
     * 素材列表
     * @return mixed
     * @throws \Exception
     */
    public function materialList()
    {
        $params                  = $this->request->all();
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['enterprise_id'] = $this->enterpriseId;

        $list =  (new MaterialCenter)->materialList($params);
        return success($list);
    }
}