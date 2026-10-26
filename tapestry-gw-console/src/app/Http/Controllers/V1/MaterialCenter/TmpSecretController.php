<?php

namespace App\Http\Controllers\V1\MaterialCenter;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\MaterialCenter;

class TmpSecretController extends BaseController
{
    /**
     * 获取腾讯cos临时秘钥
     * @return mixed
     * @throws \Exception
     */
    public function tmpSecret()
    {
        $params = $this->request->all();

        $data = (new MaterialCenter)->getTmpSecret($params);

        return success($data);
    }
}