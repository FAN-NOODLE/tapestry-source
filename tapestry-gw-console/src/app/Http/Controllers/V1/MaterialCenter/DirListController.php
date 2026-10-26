<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class DirListController extends BaseController
{
    /**
     * 文件夹列表
     * @return mixed
     * @throws \Exception
     */
    public function dirList()
    {
        $params                  = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new MaterialCenter)->dirList($params);
        return success($list);
    }
}