<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:28
 */

namespace App\Http\Controllers\Service\V1\Department;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Department;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 部门详情
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
        $dpId          = $request->post('dp_id', '');
        $filterFields  = $request->post('filter_fields', []);

        $result = (new Department())->detail($qyDomain, $appIdentifier, $dpId, $filterFields);

        return success($result);
    }
}
