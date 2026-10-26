<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class RecycleBinListController extends BaseController
{
    /**
     * 回收站素材列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['search', 'sort', 'page', 'page_size']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->recycleBinList($params));
    }
}
