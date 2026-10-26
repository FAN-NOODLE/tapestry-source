<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use CrsSdk\Shop\MaterialCenter;
use App\Http\Controllers\V1\BaseController;

class GoodsFileController extends BaseController
{
    /**
     * 商品管理上传文件
     * @return mixed
     * @throws \Exception
     */
    public function goodsFile()
    {
        $params                  = $this->request->input();
        $params['admin_id']      = $this->enterpriseAdminId;
        $params['enterprise_id'] = $this->enterpriseId;
        return (new MaterialCenter())->goodsUploadFile($params);
    }
}