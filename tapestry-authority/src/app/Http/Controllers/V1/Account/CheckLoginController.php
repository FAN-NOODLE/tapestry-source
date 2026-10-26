<?php

namespace App\Http\Controllers\V1\Account;

use App\Services\EnterpriseAdminService;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;

class CheckLoginController extends BaseController
{
    public function index(Request $request)
    {
        $account   = trim($request->post('account'));
        $password  = trim($request->post('password'));
        $loginType = trim($request->post('login_type', EnterpriseAdminService::LOGIN_BY_PASSWORD));
        $client    = $this->request->input('client', '');
        $account   = str_replace('+86-', '', $account);
        if (empty($account)) {
            notice(700001);
        }

        if (empty($password)) {
            notice(700002);
        }

        if (empty($client)) {
            notice(700013);
        }

        (new EnterpriseAdminService())->checkLogin($loginType, $account, $password, $client);

        return success();
    }
}
