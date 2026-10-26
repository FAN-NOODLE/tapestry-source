<?php

namespace App\Http\Controllers\V1\WithdrawaSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;
use Illuminate\Http\Request;

class CheckController extends BaseController
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
        $input                  = $this->request->input();

        $data = (new Withdrawal())->checkTransferNo($input);

        return success($data);
    }
}
