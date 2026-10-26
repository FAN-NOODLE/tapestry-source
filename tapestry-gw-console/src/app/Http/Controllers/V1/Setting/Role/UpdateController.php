<?php
/*
 * Created by PhpStorm.
 * User: houyingcai
 * Date Time: 2020/4/15 15:12
 */

namespace App\Http\Controllers\V1\Setting\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Role;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class UpdateController extends BaseController
{
    /**
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params = $request->only([
            'role_id',
            'name',
            'powers',
            'type'
        ]);

        $this->validate($this->request, [
            'role_id' => 'required|integer',
            'name'    => 'required',
            'type'    => ['required', Rule::in([1, 2])]
        ], [
            'role_id.required' => '角色ID不能为空',
            'role_id.integer'  => '角色ID类型有误',
            'name.required'    => '角色名称不能为空',
            'type.required'    => '权限类型不能为空',
            'type.in'          => '权限类型有误',
        ]);

        // 更新角色
        $role = new Role();
        $role->editRole([
            'role_id'       => intval($params['role_id']),
            'name'          => trim($params['name']),
            'type'          => $params['type'],
            'powers'        => $params['powers'] ?? '',
            'enterprise_id' => $this->enterpriseId,
            'operator_id'   => $this->enterpriseAdminId
        ]);

        // 写入日志
        $this->writeOperatorLog('编辑角色:' . $params['name']);

        return success();
    }
}
