<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     * 员工详情
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
        try {
            $result = (new Member())->detail($qyDomain, $appIdentifier, $params);
        } catch (\Exception $e) {
            $result = [];
        }

        return success($result);
    }
}
