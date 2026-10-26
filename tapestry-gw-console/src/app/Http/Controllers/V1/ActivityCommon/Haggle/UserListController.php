<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/19 0019
 * Time: 20:28
 */

namespace App\Http\Controllers\V1\ActivityCommon\Haggle;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\MicroActivity\Haggle;

class UserListController extends BaseController
{

    /**
     * 发起砍价用户列表
     * @return \Illuminate\Http\JsonResponse
     * @throws \Exception
     */
    function index()
    {
        $params = $this->request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $result = (new Haggle())->initList($params);
        return success($result);
    }


}
