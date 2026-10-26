<?php
/*
 * Created by PhpStorm.
 * User: houyingcai
 * Date Time: 2020/4/15 16:20
 */

namespace App\Http\Controllers\V1\Setting\Role;

use Exception;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Role;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{
    /**
     * 删除角色
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function index(Request $request)
    {
        $roleId = $request->input('role_id');

        $role = new Role();
        $res  = $role->deleteRole(['role_id' => $roleId]);

        if (empty($res)) {
            notice(450102);
        }

        $this->writeOperatorLog('删除角色:' . $res['name']);
        return success();
    }
}
