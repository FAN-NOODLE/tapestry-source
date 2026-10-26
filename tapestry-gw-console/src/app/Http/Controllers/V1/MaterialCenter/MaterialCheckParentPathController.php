<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialCheckParentPathController extends BaseController
{
    /**
     * 回收站恢复判断原路径是否完整
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->recoverFromBinCheckParentPath($params));
    }
}
