<?php


namespace App\Http\Controllers\V1\Common;

use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Common;

class SendCodeController extends BaseController
{
    /**
     * 发送验证码
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index()
    {
        $input                = $this->request->only([
            'mobile',
            'type'
        ]);
        $input['client']      = Constant::CLIENT;

        $input['type'] = $input['type'] ?? 0;
        // 非登录时调用发送验证码接口的规则传输 --type = 1是为修改管理员手机号-完善信息部分
        if (isset($this->enterpriseId) && isset($this->enterpriseAdminId) && $input['type'] == 1) {
            $input['module_rule'] = [
                'check_enterprise_admin' => [
                    'id'            => $this->enterpriseAdminId,
                    'enterprise_id' => $this->enterpriseId
                ]
            ];
        }

        (new Common())->SendCode($input);

        return success();
    }
}
