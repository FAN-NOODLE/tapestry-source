<?php


namespace App\Http\Controllers\V1\Admin;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Boss\EnterpriseConfig;
use CrsSdk\MicroAuthority\Account;

class ChangeEnterpriseListController extends BaseController
{
    public function index()
    {
        $adminId = $this->enterpriseAdminId;

        $accountDetail = (new Account())->checkAccount(['admin_id' => $adminId]);

        $account = $accountDetail['account'];

        $accountEnterpriseList = (new Account())->accountEnterpriseList([
            'account' => $account,
            'client'  => 'console',
            'version' => 'v2'
        ]);
        return success($accountEnterpriseList);
    }
}
