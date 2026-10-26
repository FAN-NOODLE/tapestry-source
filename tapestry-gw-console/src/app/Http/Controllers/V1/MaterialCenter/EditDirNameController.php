<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class EditDirNameController extends BaseController
{
    /**
     * 编辑文件夹名称
     * @return mixed
     * @throws \Exception
     */
    public function editDirName()
    {
        $params                  = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $this->writeOperatorLog('编辑文件夹名称');
        return (new MaterialCenter)->editDirName($params);
    }
}