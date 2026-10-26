<?php

namespace App\Http\Controllers\V1\WithdrawaSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;
use Illuminate\Http\Request;

class DetailController extends BaseController
{
    /**
     *
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $params['enterprise_id'] = $this->enterpriseId;

        $setInfo = (new Withdrawal())->detail($params);

        return success($setInfo);
    }
}
