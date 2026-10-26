<?php

namespace App\Http\Controllers\V1\Account;

use App\Models\EnterpriseAdmin;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;

class CheckAccountController extends BaseController
{
    public function index(Request $request)
    {
        $account = $request->post('account');
        $adminId = trim($request->post('admin_id'));
        if (empty($account) && empty($adminId)) {
            notice(700001);
        }
        $account = str_replace('+86-', '', $account);
        if (!is_numeric($account)) {
            try {
                if (!$trueAccount = $this->redis->get($account)) {
                    notice(700032);
                }
                $account = $trueAccount;
            } catch (\Exception $e) {
            }
        }
        if (!empty($adminId)) {
            $where = ['id', $adminId];
        } else {
            $where = ['account', 'like', '%' . $account];
        }

        $accountDetail = (new EnterpriseAdmin())->accountDetail([$where], ['id', 'account'], $account);

        // 账号尚未开通
        if (empty($accountDetail)) {
            notice(700005);
        }

        return success($accountDetail);
    }
}
