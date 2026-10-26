<?php


namespace App\Http\Controllers\V1\Admin;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;

class LoginSelectEnterpriseListController extends BaseController
{

    public function index()
    {
        $account   = trim($this->request->post('account'));

        if (empty($account)) {
            error(451001);
        }

        $accountEnterpriseList = (new Account())->accountEnterpriseList([
            'account' => $account,
            'client'  => 'console',
            'version' => 'v2'
        ]);

        return success($accountEnterpriseList);
    }
}
