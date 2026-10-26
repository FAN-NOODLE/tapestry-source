<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use App\Http\Requests\AccountRequest;
use App\Models\EnterpriseAdmin;
use App\Services\EnterpriseAdminService;
use App\Services\EnterpriseService;
use Illuminate\Support\Facades\Log;
use function Composer\Autoload\includeFile;

class StatusController extends BaseController
{
    /**
     * 修改账号状态
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     * @api {domain}/micro-authority/v1/account/status
     */
    public function index()
    {
        $input = $this->request->only([
            'id',
            'enterprise_id',
            'status',
            'ids',
            'all'
        ]);

        if (isset($input['id']) && $input['id'] === 0) {
            error(701014);
        }

        $this->checkPermissionsKey($input['enterprise_id'],$input['id']);

        $input['ids'] = $input['ids'] ?? [];
        $input['all'] = $input['all'] ?? '';

        $openCount = 0;
        if (isset($input['id']) && $input['id'] > 0 && $input['status'] == EnterpriseAdmin::STATUS_OPEN) {
            $openCount++;
        }

        if (!empty($input['ids']) && $input['status'] == EnterpriseAdmin::STATUS_OPEN) {
            $openCount = $openCount + count($input['ids']);
        }

        // 验证
        AccountRequest::status($input);

        // 如果已经启用的账号+要启用的账号数大于限制数量 则不可启用
//        $list = (new EnterpriseAdminService())->getEnterpriseChildAdminList($input['enterprise_id']);
//        // 当前启用数量
//        $currentOpenCount = count($list);

        $currentOpenCount = EnterpriseAdmin::query()
            ->where('enterprise_id', '=', $input['enterprise_id'])
            ->where('admin_available_system', '!=', json_encode([2])) // 过滤智慧导购
            ->where('admin_available_system', '!=', '[]') // 过滤空的可用系统
            ->where('is_super', EnterpriseAdmin::NO_SUPER)
            ->where('status', EnterpriseAdmin::STATUS_OPEN)
            ->count();

        $detail = (new EnterpriseService())->enterpriseDetail($input['enterprise_id']);
        // 限制数量
        $limitCount = $detail['account_limit'];

        // 可启用数量
        $canCount = $limitCount - $currentOpenCount;

        // 如果要启用的数量大于可启用则提示
//        if ($openCount > $canCount) {
//            notice(701011, [$limitCount]);
//        }

        // 修改状态
        try {
            (new EnterpriseAdminService())->updateEnterpriseAdminStatus($input);
        } catch (\Exception $exception) {
            Log::info('修改状态失败' . $exception->getMessage(), $input);
            notice(701010);
        }

        return success([]);
    }

}