<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/29
 * Time: 18:01
 */

namespace App\Http\Controllers\V1\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Role;

class ListController extends BaseController
{
    /**
     * 角色列表
     *
     * @return array
     */
    public function index()
    {
        // 调用sdk
        $list = (new Role())->list(['qy_domain' => $this->qyDomain]);

        // 初始化待返数据
        $data = [];
        if ($list['list']) {
            // 数据格式化
            foreach ($list['list'] as $role) {
                $data[] = [
                    'role_id'   => $role['roleId'],
                    'role_name' => $role['roleName']
                ];
            }
        }

        return success(['data' => $data]);
    }
}
