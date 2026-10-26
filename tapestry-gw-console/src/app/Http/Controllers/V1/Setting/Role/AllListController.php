<?php

namespace App\Http\Controllers\V1\Setting\Role;

use Exception;
use CrsSdk\MicroAuthority\Role;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class AllListController extends BaseController
{
    /**
     * 角色列表 - 不分页
     * @return JsonResponse
     * @throws Exception
     */
    public function index()
    {
        $params                  = $this->request->all();
        $params['enterprise_id'] = $this->enterpriseId;

        $data = (new Role())->allList($params);

        return success($data);
    }
}