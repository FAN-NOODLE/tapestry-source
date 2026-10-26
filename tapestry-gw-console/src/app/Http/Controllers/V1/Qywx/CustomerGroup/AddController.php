<?php

namespace App\Http\Controllers\V1\Qywx\CustomerGroup;

use CrsSdk\Work\CustomerGroup;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * 添加用户分组
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input                  = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('添加用户分组');

        $data = (new CustomerGroup())->add($input);

        return success(['id' => $data['data']['id']]);
    }
}