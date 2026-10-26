<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/29
 * Time: 18:02
 */

namespace App\Http\Controllers\V1\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Role;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 角色详情
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        // 角色id
        $roleId = $request->get('role_id', '');
        // 角色id不能为空
        if (!$roleId) {
            error(320201);
        }

        // 调用sdk
        $params = [
            'qy_domain' => $this->qyDomain,
            'condition' => ['roleId' => $roleId]
        ];
        $detail = (new Role())->detail($params);

        // 初始化待返数据
        $data = [];
        if ($detail) {
            // 数据格式化
            $data = [
                'role_id'   => $detail['roleId'],
                'role_name' => $detail['roleName']
            ];
        }

        return success(['data' => $data]);
    }
}
