<?php


namespace App\Http\Controllers\V1\Role;


use App\Http\Controllers\V1\BaseController;
use App\Services\EnterpriseRoleService;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\Rule;
use Illuminate\Validation\ValidationException;

class EditController extends BaseController
{
    /**
     * 编辑角色
     *
     * @param EnterpriseRoleService $enterpriseRoleService
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(EnterpriseRoleService $enterpriseRoleService)
    {
        $params = $this->request->only([
            'role_id',
            'name',
            'type',
            'powers',
            'operator_id',
            'enterprise_id'
        ]);
        $this->validate($this->request, [
            'role_id'       => 'required|integer',
            'name'          => 'required',
            'enterprise_id' => 'required|integer',
            'operator_id'   => 'required|integer',
            'type'          => ['required', Rule::in(EnterpriseRoleService::ROLE_TYPE)]
        ], [
            'role_id.required'       => '角色ID不能为空',
            'role_id.integer'        => '角色ID类型错误',
            'name.required'          => '角色名称不能为空',
            'type.required'          => '权限类型不能为空',
            'type.in'                => '权限类型有误',
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer'  => '企业ID类型错误',
            'operator_id.required'   => '操作人ID不能为空',
            'operator_id.integer'    => '操作人ID类型错误',
        ]);

        $auth_rule = $params['powers'];
        // 数据权限格式为json，功能权限格式为前端加密后的字符串
        if ($params['type'] == EnterpriseRoleService::ROLE_TYPE[1]) {
            $auth_rule = json_encode($params['powers']);
        }

        // 查询角色是否存在
        $role = $enterpriseRoleService->getRoleById($params['role_id']);
        if (empty($role)) {
            notice(702021);
        }

        $where    = ['id' => $params['role_id']];
        $roleData = [
            'enterprise_id' => $params['enterprise_id'],
            'type'          => $params['type'],
            'role_name'     => trim($params['name']),
            'auth_rule'     => $auth_rule,
            'operator_id'   => $params['operator_id']
        ];

        $enterpriseRoleService->updateRole($where, $roleData);

        return success();
    }
}
