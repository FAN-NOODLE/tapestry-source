<?php

namespace App\Http\Controllers\V1\Postage;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroGoodsV1\Postage;

class AddController extends BaseController
{
    /**
     * 新增邮费模板接口
     *
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->all();
        $input['enterprise_id'] = $this->enterpriseId;
        $input['store_id'] = $this->request->input('store_id', $this->enterpriseId);
        $data = (new Postage())->add($input);

        return success($data);
    }
}
