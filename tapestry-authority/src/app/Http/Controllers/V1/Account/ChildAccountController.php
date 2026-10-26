<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Models\EnterpriseAdmin;
use App\Services\EnterpriseAdminService;

class ChildAccountController extends BaseController
{
    /**
     * @is_super 1-只获取超管 2-获取所有非超管
     * @status 1-启用 2-禁用
     */
    public function index()
    {
        $input = $this->request->only([
            'enterprise_id',
            'is_super',
            'status'
        ]);

        $input['is_super'] = $input['is_super'] ?? EnterpriseAdmin::NO_SUPER;
        $input['status'] = $input['status'] ?? EnterpriseAdmin::STATUS_OPEN;

        $list = (new EnterpriseAdminService())->getEnterpriseChildAdminList($input['enterprise_id'], $input['is_super'], $input['status']) ?? [];

        return success([
            'list' => $list
        ]);


    }

}