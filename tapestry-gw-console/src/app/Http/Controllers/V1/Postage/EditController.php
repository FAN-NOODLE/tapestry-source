<?php

namespace App\Http\Controllers\V1\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Postage;

class EditController extends BaseController
{
    /**
     * 编辑邮费模板接口
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $data = (new Postage())->edit($input);

        return success($data);
    }
}
