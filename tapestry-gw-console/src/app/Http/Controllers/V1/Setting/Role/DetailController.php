<?php
/*
 * Created by PhpStorm.
 * User: houyingcai
 * Date Time: 2020/4/15 22:13
 */

namespace App\Http\Controllers\V1\Setting\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Role;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;


class DetailController extends BaseController
{
    /**
     * 角色详情接口
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $roleId = $request->input('role_id');
        $role = new Role();
        $roleInfo = $role->detailRole(['role_id' => $roleId]);

        if (empty($roleInfo)) {
            notice(450100);
        }

        return success([
            'id'                => $roleInfo['id'],
            'role_name'         => $roleInfo['role_name'] ?? '',
            'type'              => $roleInfo['type'],
            'powers'            => $roleInfo['powers'] ?? '',
            'role_admins_count' => $roleInfo['role_admins_count'],
            'admin_name'        => $roleInfo['admin_name'] ?? '',
            'created_at'        => $roleInfo['created_at'],
            'updated_at'        => $roleInfo['updated_at']
        ]);
    }
}
