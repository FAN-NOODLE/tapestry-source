<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Enterprise;

class AuthAdminInfoController extends BaseController
{
    /**
     * 用token查询登录的企业信息
     *
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'token' => 'required|string'
        ]);
        $input = $this->request->input();

        return (new Enterprise())->authAdminInfo([
            'token' => $input['token']
        ]);
    }
}
