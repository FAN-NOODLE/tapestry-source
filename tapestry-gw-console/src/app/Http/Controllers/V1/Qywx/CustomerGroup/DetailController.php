<?php

namespace App\Http\Controllers\V1\Qywx\CustomerGroup;

use CrsSdk\Work\CustomerGroup;
use App\Http\Controllers\V1\BaseController;

class DetailController extends BaseController
{
    /**
     * 获取用户分组详情
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new CustomerGroup())->detail($input);

        return success($list['data']);
    }
}