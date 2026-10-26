<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use CrsSdk\MicroAuthority\Account;

/**
 *
 * Class PerfectController
 * @package App\Http\Controllers\V1\account
 */
class PerfectController extends BaseController
{
    /**
     * 完善账号信息
     * @api {domain}/gw-console/v1/setting/account/perfect
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only([
            'mobile',
            'password',
        ]);

        if (isset($input['mobile'])) {
            $input['mobile'] = trim($input['mobile']);
        }
        $input['enterprise_id'] = $this->enterpriseId;
        $input['operator_id'] = $this->enterpriseAdminId;
        $input['id'] = $this->enterpriseAdminId;

        // 验证
        AccountRequest::perfect($input);

        // 调用微服务进行完善
        (new Account())->perfectAccount($input);

        return success([]);
    }
}