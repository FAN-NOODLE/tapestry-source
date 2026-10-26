<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/29
 * Time: 18:03
 */

namespace App\Http\Controllers\V1\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Role;
use Illuminate\Http\Request;

class EditController extends BaseController
{
    /**
     * 编辑角色
     *
     * @param Request $request
     *
     * @return array|mixed
     */
    public function index(Request $request)
    {
        // 角色id
        $roleId = $request->post('role_id', '');
        // 角色id不能为空
        if (!$roleId) {
            error(320201);
        }

        // 角色名称
        $roleName = $request->post('role_name', '');
        // 角色名称不能为空
        if (!$roleName) {
            error(320202);
        }

        // 调用sdk
        $params = [
            'qy_domain' => $this->qyDomain,
            'condition' => [
                'roleId'   => $roleId,
                'roleName' => $roleName
            ]
        ];
        $detail = (new Role())->add($params);

        // 初始化待返数据
        $data = [];
        if ($detail) {
            // 数据格式化
            $data = [
                'role_id'   => $detail['roleId'],
                'role_name' => $detail['roleName']
            ];
        }

        // 日志记录
        $this->writeOperatorLog('编辑了 ' . $roleName . ' 角色');

        return success(['data' => $data]);
    }
}
