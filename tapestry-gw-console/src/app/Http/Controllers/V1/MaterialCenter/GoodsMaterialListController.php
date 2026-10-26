<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class GoodsMaterialListController extends BaseController
{
    /**
     * 商品素材中心
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function goodsMaterialList()
    {
        $params                  = $this->request->all();
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new MaterialCenter)->goodsMaterialList($params);
        return success($list);
    }
}