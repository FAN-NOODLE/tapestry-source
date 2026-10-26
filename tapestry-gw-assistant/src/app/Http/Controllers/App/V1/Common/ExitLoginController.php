<?php

namespace App\Http\Controllers\App\V1\Common;


use App\Http\Controllers\App\V1\BaseController as Controller;
use CrsSdk\MicroAuthority\Account;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cache;

class ExitLoginController extends Controller
{
    protected $isRequireLogin = false;

    /**
     * 退出登录
     *
     * @var int
     */
    const LOGIN_TYPE_CODE = 2;

    /**
     * 登录返回
     *
     * @param Request $request
     * @param Account $account
     * @return \Illuminate\Http\JsonResponse|void
     */
    public function __invoke(Request $request, Account $account)
    {
        $mobile = $this->loginMobile();
        Cache::forget($mobile);
    }
}
