<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;

class AccountSystemController extends BaseController
{
    /**
     * 获取创建账号可用系统选项
     * @api {domain}/gw-console/v1/setting/account/account-system
     * @return \Illuminate\Http\JsonResponse
     */
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;

        $list = (new Account())->accountSystemCheckbox($input);

        return success([
            'list' => $list['list'] ?? []
        ]);
    }

}