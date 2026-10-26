<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;

class StatusController extends BaseController
{
    /**
     * 启用和禁用管理员账号
     * @api {domain}/gw-console/v1/setting/account/status
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'status'
        ]);

        $input['enterprise_id'] = $this->enterpriseId;

        // 验证
        AccountRequest::status($input);

        (new Account())->statusAccount($input);

        return success([]);
    }

}