<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 19:32
 */

namespace App\Http\Controllers\Service\V1\Role;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Role;
use Illuminate\Http\Request;

class ListController extends BaseController
{

    /**
     * 角色列表接口
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain  = $request->post('qy_domain', '');
        $condition = $request->post('condition', []);

        $result = (new Role())->list($qyDomain, $condition);

        return success($result);
    }
}
