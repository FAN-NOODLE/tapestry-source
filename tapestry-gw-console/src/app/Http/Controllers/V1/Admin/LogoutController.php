<?php

namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\JsonResponse;

class LogoutController extends BaseController
{
    protected $requireLogin = false;
    /**
     * 企业管理员退出登录接口
     *
     * @return JsonResponse
     */
    public function index()
    {
        return success();
    }
}
