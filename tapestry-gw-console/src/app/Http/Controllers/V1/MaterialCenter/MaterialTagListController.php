<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialTagListController extends BaseController
{
    /**
     * 素材标签列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->materialTagList($params));
    }
}
