<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class CreateDirController extends BaseController
{
    /**
     * @return mixed
     * @throws \Exception
     */
    public function createDir()
    {
        $params                  = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('创建文件夹');
        $list = (new MaterialCenter)->createDir($params);
        return success($list);
    }
}