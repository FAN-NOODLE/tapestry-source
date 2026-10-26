<?php

namespace App\Http\Controllers\Service\V1\Enterprise;

use App\Library\Uc\Common;
use App\Http\Controllers\Service\V1\BaseController;

class AuthCodeInfoController extends BaseController
{
    /**
     * 扫码登录后根据auth_code获取信息
     *
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'auth_code' => 'required|string'
        ]);
        $input = $this->request->all();

        return (new Common())->authCodeLogin($input['auth_code']);
    }
}
