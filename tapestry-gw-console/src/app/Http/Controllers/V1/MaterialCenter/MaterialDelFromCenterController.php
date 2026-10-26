<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialDelFromCenterController extends BaseController
{
    /**
     * 从素材中心删除素材
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->delFromCenter($params));
    }
}
