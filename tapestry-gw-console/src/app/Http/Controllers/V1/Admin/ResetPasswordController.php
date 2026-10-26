<?php

namespace App\Http\Controllers\V1\Admin;

use App\Micros\Authority\AdminMicro;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;

class ResetPasswordController extends BaseController
{

    /**
     * 非必须登录
     *
     * @var bool
     */
    protected $requireLogin = false;

    /**
     * 企业ID
     */
    protected $enterpriseId = null;

    /**
     * 重置密码
     *
     */
    public function index(Request $request): JsonResponse
    {
        $input = $request->all();
        if (empty($input['code'])) {
            notice(451016);
        }
        $res = AdminMicro::resetPassword($input);
        return success($res);
    }

    /**
     * 重置密码验证
     *
     */
    public function checkResetPassword(Request $request): JsonResponse
    {
        $input = $request->all();
        $res = AdminMicro::checkResetPassword($input);
        return success($res);
    }


}
