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

class DetailController extends BaseController
{
    /**
     * 角色详情
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

        $result = (new Role())->detail($qyDomain, $condition);

        return success($result);
    }
}
