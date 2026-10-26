<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Library\Uc\Member;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class ListByDpIdsController extends BaseController
{


    /**
     * 员工列表（忽略权限）
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params = $request->post('params', []);

        $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}