<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use CrsSdk\Shop\MaterialCenter;
use App\Http\Controllers\V1\BaseController;

class UploadFileController extends BaseController
{
    /**
     * 上传素材
     * @return mixed
     * @throws \Exception
     */
    public function uploadFile()
    {
        $params                  = $this->request->input();
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['enterprise_id'] = $this->enterpriseId;

        $list = (new MaterialCenter())->uploadFile($params);
        return success($list);
    }
}