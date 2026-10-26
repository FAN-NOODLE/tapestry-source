<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 19:34
 */

namespace App\Http\Controllers\Service\V1\Role;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Role;
use Illuminate\Http\Request;

class DeleteController extends BaseController
{

    /**
     * 删除角色
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

        (new Role())->delete($qyDomain, $condition);

        return success(['data' => []]);
    }
}
