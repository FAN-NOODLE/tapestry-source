<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/12/3
 * Time: 18:32
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\DepartmentService;
use Illuminate\Http\Request;

class StoreDptIdController extends BaseController
{
    /**
     * 获取门店dptId
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
        $dptId         = (new DepartmentService())->getStoreDptId($qyDomain, $appIdentifier);

        return success([
            'data' => [
                'dptId' => $dptId
            ]
        ]);
    }
}
