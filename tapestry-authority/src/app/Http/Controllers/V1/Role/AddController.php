<?php


namespace App\Http\Controllers\V1\Role;


use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseRoleService;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class AddController extends BaseController
{
    /**
     * 新增角色
     *
     * @param EnterpriseRoleService $enterpriseRoleService
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(EnterpriseRoleService $enterpriseRoleService)
    {
        $params = $this->request->only([
            'name',
            'type',
            'powers',
            'operator_id',
            'enterprise_id',
            'identity_type'
        ]);
        $this->validate($this->request, [
            'name'          => 'required',
            'enterprise_id' => 'required|integer',
            'operator_id'   => 'required|integer',
            'type'          => ['required', Rule::in(EnterpriseRoleService::ROLE_TYPE)]
        ], [
            'name.required'          => '角色名称不能为空',
            'type.required'          => '权限类型不能为空',
            'type.in'                => '权限类型有误',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID类型错误',
            'operator_id.required'   => '操作人ID不能为空',
            'operator_id.integer'    => '操作人ID类型错误',
        ]);

        /*
         * 添加非常规管理员
         */
        if (!empty($params['identity_type']) && $params['identity_type'] != 1) {
            $isStoreRole = $enterpriseRoleService->getRoleDetailByWhere([
                'enterprise_id' => $params['enterprise_id'],
                'identity_type' => $params['identity_type']
            ]);
            if (!empty($isStoreRole)) {
                return success(['role_id' => $isStoreRole->id]);
            }
        }

        $auth_rule = $params['powers'];
        // 数据权限格式为json，功能权限格式为前端加密后的字符串
        if ($params['type'] == EnterpriseRoleService::ROLE_TYPE[1]) {
            $auth_rule = json_encode($params['powers']);
        }

        $roleData = [
            'enterprise_id' => $params['enterprise_id'],
            'type'          => $params['type'],
            'role_name'     => trim($params['name']),
            'auth_rule'     => $auth_rule,
            'operator_id'   => $params['operator_id'],
            'identity_type' => $params['identity_type'] ?? 1
        ];

        $roleId = $enterpriseRoleService->saveRole($roleData);

        return success(['role_id' => $roleId]);
    }
}
