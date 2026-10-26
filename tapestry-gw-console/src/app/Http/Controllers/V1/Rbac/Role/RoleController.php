<?php

namespace App\Http\Controllers\V1\Rbac\Role;

use App\Common\Constant;
use CrsSdk\MicroAuthority\Role;
use App\Http\Controllers\V1\BaseController;

class RoleController extends BaseController
{
    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function lists()
    {
        $params = $this->request->all();

        $data = (new Role())->lists([
            'page'          => $params['page'] ?? Constant::PAGE,
            'size'          => $params['size'] ?? Constant::LIMIT,
            'admin_id'      => $this->enterpriseAdminId,
            'role_name'     => $params['role_name'] ?? '',
            'sort_type'     => $params['sort_type'] ?? 'DESC',
            'admin_name'    => $params['admin_name'] ?? '',
            'add_account'   => $params['add_account'] ?? false,
            'enterprise_id' => $this->enterpriseId,
        ]);

        return success($data);
    }

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function add()
    {
        $params = $this->request->all();

        $data = (new Role())->add([
            'name'               => $params['name'],
            'role_pid'           => $params['role_pid'] ?? 0,
            'weight'             => $params['weight'] ?? 0,
            'admin_id'           => $this->enterpriseAdminId,
            'permission_id_list' => $params['permission_id_list'],
            'enterprise_id'      => $this->enterpriseId,
        ]);

        return success($data);
    }

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function edit()
    {
        $params = $this->request->all();

        $data = (new Role())->edit([
            'name'               => $params['name'],
            'role_id'            => $params['role_id'],
            'role_pid'           => $params['role_pid'] ?? 0,
            'weight'             => $params['weight'] ?? 0,
            'admin_id'           => $this->enterpriseAdminId,
            'permission_id_list' => $params['permission_id_list'],
            'enterprise_id'      => $this->enterpriseId,
        ]);

        return success($data);
    }

    public function detail()
    {
        $params = $this->request->all();

        $data = (new Role())->detail([
            'role_id'        => $params['role_id'],
            'admin_id'       => $this->enterpriseAdminId,
            'version_config' => 0,
            'enterprise_id'  => $this->enterpriseId,
        ]);

        return success($data);
    }

    /**
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function delete()
    {
        $params = $this->request->all();

        $data = (new Role())->delete([
            'role_id'       => $params['role_id'],
            'admin_id'      => $this->enterpriseAdminId,
            'enterprise_id' => $this->enterpriseId,
        ]);

        return success($data);
    }
}