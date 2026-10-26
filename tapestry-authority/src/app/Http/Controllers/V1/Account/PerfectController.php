<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Services\EnterpriseAdminService;
use Illuminate\Support\Facades\Log;

/**
 * 完善账户信息
 * @api {domain}/micro-authority/v1/account/perfect
 * Class PerfectController
 * @package App\Http\Controllers\V1\Account
 */
class PerfectController extends BaseController
{
    public function index()
    {
        $input = $this->request->only([
            'id',
            'operator_id',
            'enterprise_id',
            'mobile',
            'password',
        ]);

        // 验证
        AccountRequest::perfect($input);

        // mobile拆解
        // 验证手机号
        $mobileArr = explode('-', $input['mobile']);
        $mobile = $mobileArr[1] ?? $input['mobile'];
        $preg = "/^1[3-9]\d{9}$/";
        if (!preg_match($preg, $mobile)) {
            notice(701001, [], [], 200); // 请输入正确的手机号码
        }

        if (count($mobileArr) > 1) {
            $input['global_roaming'] = $mobileArr[0];
        }

        $input['account'] = count($mobileArr) == 1 ? '+86-' . $input['mobile'] : $input['mobile'];
        $input['mobile'] = $mobile;

        try {
            (new EnterpriseAdminService())->perfectAdmin($input);
        } catch (\Exception $exception) {
            Log::info('完善账户失败' . $exception->getMessage(), $input);
            notice(701008, [], [], 200);
        }

        return success([]);
    }
}