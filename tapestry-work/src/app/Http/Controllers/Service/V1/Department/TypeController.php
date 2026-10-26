<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 16:46
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use Illuminate\Http\Request;

class TypeController extends BaseController
{
    /**
     * 部门类型配置列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        $result = (new Department())->typeList($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}
