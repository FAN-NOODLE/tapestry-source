<?php

namespace App\Http\Controllers\V1\WithdrawaSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class AddController extends BaseController
{


    public function index(Request $request)
    {

        $requests = $request->only(['transfer_no', 'type']);

        $rules = [
            'transfer_no' => 'required',
        ];
        $messages = [
            'transfer_no.required' => '转账项目编号不能为空',
        ];

        Validator::make($requests, $rules, $messages)->validate();
        $params = [
            'enterprise_id'  => $this->enterpriseId,
            'transfer_no'        => $requests['transfer_no'],
            'type'           => $requests['type'] ?? [],
        ];
        (new Withdrawal())->add($params);

        return success();

    }
}