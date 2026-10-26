<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;
use function GuzzleHttp\Psr7\try_fopen;

/**
 * 编辑管理员账号
 * Class EditController
 * @package App\Http\Controllers\V1\account
 */
class EditController extends BaseController
{
    /**
     * 账号编辑
     * @api {domain}/gw-console/v1/setting/account/edit
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'account',
            'username',
            'password',
            'mobile',
            'email',
            'range_type',
            'range',
            'admin_available_system',
            'role_ids',
        ]);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id'] = $this->enterpriseAdminId;

        // AccountRequest::edit($input);

        (new Account())->editAccount($input);

        return success([]);
    }

}