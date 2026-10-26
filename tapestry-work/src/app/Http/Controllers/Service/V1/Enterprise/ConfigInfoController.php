<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;

class ConfigInfoController extends BaseController
{
    /**
     * 获取企业配置信息
     *
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required|string'
        ]);
        $input = $this->request->all();

        return (new Enterprise())->configInfo($input['enterprise_id']);
    }
}
