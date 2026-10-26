<?php
/*
 * Created by PhpStorm.
 * User: houyingcai
 * Date Time: 2020/4/15 11:23
 */


namespace App\Http\Controllers\V1\Setting\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Role;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Validation\Rule;

class AddController extends BaseController
{
    /**
     * 添加角色
     *
     * @param Request $request
     * @return JsonResponse|void
     * @throws Exception
     */
    public function index(Request $request)
    {
        $params = $this->request->only([
            'name',
            'type',
            'powers',
        ]);
        $this->validate($this->request, [
            'name' => 'required',
            'type' => ['required', Rule::in([1, 2])]
        ], [
            'name.required' => '角色名称不能为空',
            'type.required' => '权限类型不能为空',
            'type.in'       => '权限类型有误',
        ]);

        $role   = new Role();
        $result = $role->addRole([
            'name'          => trim($params['name']),
            'type'          => $params['type'],
            'powers'        => $params['powers'] ?? '',
            'enterprise_id' => $this->enterpriseId,
            'operator_id'   => $this->enterpriseAdminId
        ]);

        if (!$result['role_id']) {
            notice(450101);
        }

        // 写入日志
        $this->writeOperatorLog('添加角色:' . trim($params['name']));

        return success();
    }
}
