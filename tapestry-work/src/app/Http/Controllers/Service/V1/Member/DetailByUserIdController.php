<?php

namespace App\Http\Controllers\Service\V1\Member;

use App\Library\Uc\Member;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;

class DetailByUserIdController extends BaseController
{
    /**
     * @param Request $request
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);

        $result = (new Member())->detailByUserId($qyDomain, $appIdentifier, $params);

        return success($result);
    }
}