<?php
/*
 * Created by PhpStorm.
 * User: houyingcai
 * Date Time: 2020/4/16 17:50
 */

namespace App\Http\Controllers\V1\Setting\Role;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroAuthority\Role;
use Exception;
use Illuminate\Http\JsonResponse;

class ListController extends BaseController
{
    /**
     * 角色列表接口（分页）
     *
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $params                  = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;

        $role = new Role();
        $data = $role->listRole($params);

        if (!empty($data['list'])) {
            return success([
                'page'  => $data['page'],
                'size'  => $data['size'],
                'total' => $data['total'],
                'list'  => $this->formatData($data['list'])
            ]);
        }

        return success(['list' => []]);
    }

    /**
     * 格式化数据
     *
     * @param $data
     * @return array
     */
    public function formatData($data)
    {
        $resultData = [];
        foreach ($data as $key => $val) {
            $adminName        = !empty($val['enterprise_admin']) ? $val['enterprise_admin']['username'] : '';
            $resultData[$key] = [
                'id'                => $val['id'],
                'role_name'         => $val['role_name'],
                'type'              => $val['type'],
                'role_admins_count' => $val['role_admins_count'],
                'identity_type'     => $val['identity_type'],
                'admin_name'        => $adminName,
                'created_at'        => $val['created_at'],
                'updated_at'        => $val['updated_at'],
            ];
        }

        return $resultData;
    }
}
