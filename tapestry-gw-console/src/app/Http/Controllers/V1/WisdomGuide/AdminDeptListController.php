<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/10/28 0028
 * Time: 15:22
 */

namespace App\Http\Controllers\V1\WisdomGuide;


use App\Exceptions\Error;
use App\Services\AdminDeptService;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\JsonResponse;

class AdminDeptListController extends BaseController
{

    /**
     * 获取某导购部门管理范围
     * @return JsonResponse
     */
    public function index(): JsonResponse
    {
        $params = $this->request->all();
        try{
            return success((new AdminDeptService())->list($this->enterpriseId, $params['guide_mobile'] ?? ''));
        }catch (\Throwable $e){
            throw new Error($e->getCode(),$e->getMessage());
        }

    }
}
