<?php


namespace App\Http\Controllers\V1\Setting\Account;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Account;

class PerfectDetailController extends BaseController
{
    public function index()
    {
        $enterpriseAdminId = $this->request->input('id', $this->enterpriseAdminId);
        $enterpriseId      = $this->request->input('enterprise_id', $this->enterpriseId);
        $input             = [
            'id'            => $enterpriseAdminId,
            'enterprise_id' => $enterpriseId,
            'only'          => 2
        ];

        // 获取管理远账号详情
        $detail = (new Account())->detailAccount($input);

        $backData['mobile']  = $detail['mobile'];
        $backData['account'] = $detail['account'] ?? $detail['mobile'];
        $backData['has_psw'] = 0; // 未设置密码
        if (!empty($detail['password'])) {
            $backData['has_psw'] = 1; // 已设置密码
        }

        // 获取账号关联的企业列表和小程序信息
        $enterpriseList = (new Account())->accountEnterpriseList([
            'account' => $backData['account'],
            'client'  => 'gw-console'
        ]);

        $backData['enterprise_list'] = $enterpriseList['enterprise_admin_list'];

        return success($backData);
    }

}