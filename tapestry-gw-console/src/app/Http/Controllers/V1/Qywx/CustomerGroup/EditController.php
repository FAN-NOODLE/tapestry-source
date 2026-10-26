<?php

namespace App\Http\Controllers\V1\Qywx\CustomerGroup;

use CrsSdk\Work\CustomerGroup;
use App\Http\Controllers\V1\BaseController;

class EditController extends BaseController
{
    /**
     * 更新用户分组
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input = $this->request->input();
        $input['enterprise_id'] = $this->enterpriseId;

        $this->writeOperatorLog('更新用户分组');

        (new CustomerGroup())->update($input);

        return success();
    }
}