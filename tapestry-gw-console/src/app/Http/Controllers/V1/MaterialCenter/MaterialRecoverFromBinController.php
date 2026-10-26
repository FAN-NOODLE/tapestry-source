<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialRecoverFromBinController extends BaseController
{
    /**
     * 从回收站恢复素材
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->recoverFromBin($params));
    }
}
