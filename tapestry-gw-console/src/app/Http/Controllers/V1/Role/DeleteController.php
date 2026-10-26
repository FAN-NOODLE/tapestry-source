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

class DeleteController extends BaseController
{
    /**
     * 删除角色
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

        // 调用sdk
        $params = [
            'qy_domain' => $this->qyDomain,
            'condition' => ['roleId' => $roleId]
        ];

        // 查询角色，为日志服务
        $name = (new Role())->detail($params);

        (new Role())->delete($params);

        // 日志记录
        $this->writeOperatorLog('删除了 ' . $name['roleName'] . ' 角色');

        return success(['data' => []]);
    }
}
