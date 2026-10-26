<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialCheckValidController extends BaseController
{
    /**
     * 判断素材是否有效
     * 1、综合素材、商品素材列表：必须是正常状态；
     * 2、回收站素材列表：必须是在回收站，删除状态；
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id', 'type']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->checkMaterialValid($params));
    }
}
