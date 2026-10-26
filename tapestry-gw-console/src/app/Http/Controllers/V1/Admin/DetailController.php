<?php

namespace App\Http\Controllers\V1\Admin;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\JsonResponse;

class DetailController extends BaseController
{
    /**
     * 企业管理员详情接口
     *
     * @return JsonResponse
     */
    public function index()
    {
        $admin = $this->enterpriseAdmin;
        unset($admin['password'], $admin['salt']);

        return success($admin);
    }
}
