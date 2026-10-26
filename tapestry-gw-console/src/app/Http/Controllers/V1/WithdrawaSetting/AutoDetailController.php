<?php

namespace App\Http\Controllers\V1\WithdrawaSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;
use Illuminate\Http\Request;

class AutoDetailController extends BaseController
{
    /**
     *
     *
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(Request $request)
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;

        $setInfo = (new Withdrawal())->autoDetail($params);

        return success($setInfo);
    }
}
