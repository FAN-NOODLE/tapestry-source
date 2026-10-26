<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-31
 * Time: 15:55
 */

namespace App\Http\Controllers\V1\Department;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Sdk;
use Illuminate\Http\Request;

class MemberListByScrmController extends BaseController
{
    /**
     * 根据部门列表
     * @param Request $request
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = Sdk::work()
            ->setApi('service/v1/department/member-list-by-dp-ids')
            ->post($params);
        return success($result);
    }
}
