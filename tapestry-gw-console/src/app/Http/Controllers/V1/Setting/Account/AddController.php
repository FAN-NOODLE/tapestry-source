<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Support\Facades\Log;

/**
 * 创建管理员账号-新
 * Class AddController
 * @package App\Http\Controllers\V1\account
 */
class AddController extends BaseController
{
    /**
     * 新增账号
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @api {domain}/gw-console/v1/setting/account/add
     */
    public function index()
    {
        $input = $this->request->only([
            'username',
            'account',
            'mobile',
            'password',
            'email',
            'range_type',
            'range',
            'admin_available_system',
            'role_ids',
        ]);

        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id']   = $this->enterpriseAdminId;

        AccountRequest::add($input);

        // 验证手机号
        $mobileArr = explode('-', $input['mobile']);
        $mobile    = $mobileArr[1] ?? $input['mobile'];
        $preg      = "/^1[3-9]\d{9}$/";
        if (!preg_match($preg, $mobile)) {
            notice(451014); // 请输入正确的手机号码
        }

        (new Account())->addAccount($input);

        return success([]);
    }

}