<?php


namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseAdminService;
use Illuminate\Http\Request;

class AccountEnterpriseListController extends BaseController
{
    public function index(Request $request)
    {
        $account = trim($request->post('account'));
        $client  = trim($request->post('client', 'console'));
        $version = trim($request->post('version', ''));

        $enterpriseAdminList = (new EnterpriseAdminService())->enterpriseAdminList($account, $client, $version);

        return success(['enterprise_admin_list' => $enterpriseAdminList]);
    }
}
