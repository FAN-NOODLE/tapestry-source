<?php


namespace App\Http\Controllers\V1\Role;


use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseRoleService;
use Illuminate\Http\JsonResponse;

class DetailController extends BaseController
{
    /**
     * 角色详情
     *
     * @param EnterpriseRoleService $enterpriseRoleService
     * @return JsonResponse
     */
    public function index(EnterpriseRoleService $enterpriseRoleService)
    {
        $roleId = $this->request->input('role_id');
        if (empty($roleId)) {
            notice(702001);
        }

        $returnData = [];
        $roleInfo   = $enterpriseRoleService->getRoleDetail($roleId);
        if (!empty($roleInfo)) {
            if ($roleInfo['type'] == EnterpriseRoleService::ROLE_TYPE[1]) {
                // 数据权限
                $roleInfo['auth_rule'] = json_decode($roleInfo['auth_rule'], true);
            }

            $adminName = !empty($roleInfo['enterprise_admin']) ? $roleInfo['enterprise_admin']['username'] : '';
            $returnData = [
                'id'                => $roleInfo['id'],
                'role_name'         => $roleInfo['role_name'],
                'type'              => $roleInfo['type'],
                'powers'            => $roleInfo['auth_rule'],
                'role_admins_count' => $roleInfo['role_admins_count'],
                'admin_name'        => $adminName,
                'created_at'        => $roleInfo['created_at'],
                'updated_at'        => $roleInfo['updated_at']
            ];
        }

        return success($returnData);
    }
}
