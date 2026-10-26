<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialExchangeParentFolderController extends BaseController
{
    /**
     * 回收站素材变更父级文件夹(回收站恢复时存在父级文件夹被删除)
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['material_id', 'folder_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->exchangeParentFolder($params));
    }
}
