<?php


namespace App\Http\Controllers\V1\Account;


use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Services\EnterpriseAdminRoleService;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseDataRangeService;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

/**
 * 删除账号
 * Class DeleteController
 * @package App\Http\Controllers\V1\Account
 */
class DeleteController extends BaseController
{
    /**
     * @api {domain}/micro-authority/v1/account/del
     * @author lijian lijian@vchangyi.com 2020/04/13
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'enterprise_id'
        ]);

        $this->checkPermissionsKey($input['enterprise_id'],$input['id']);
        // 验证
        AccountRequest::del($input);

        // 执行删除逻辑
        try {
            DB::beginTransaction();
            
            // 删除管理员账号
            (new EnterpriseAdminService())->delEnterpriseAdmin($input);

            // 删除账号关联范围数据
            (new EnterpriseDataRangeService())->delDataRangeByAdminId($input['id']);

            // 删除账号关联角色数据
            (new EnterpriseAdminRoleService())->delAdminRoleByAdminId($input['id']);

            DB::commit();
        } catch (\Exception $exception) {
            DB::rollBack();
            Log::error('删除账号失败'.$exception->getMessage(), $input);
            error(100000, [$exception->getMessage()]);
        }

        return success([]);
    }

}