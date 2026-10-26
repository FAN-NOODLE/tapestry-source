<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;
use Illuminate\Http\Request;

class MaterialFolderListTreeController extends BaseController
{
    /**
     * 素材类型下所有文件夹列表
     *
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->only(['type', 'folder_id']);

        $params['enterprise_id'] = $this->enterpriseId;

        return success((new MaterialCenter())->folderListTree($params));
    }
}
