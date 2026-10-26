<?php

namespace App\Http\Controllers\V1\WithdrawaSetting;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Withdrawal;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class AutoAddController extends BaseController
{


    public function index(Request $request)
    {

        $requestParams = $request->only([
            'auto_status',
            'auto_at','auto_status',
            'mem_type'
        ]);

        $rules = [
            'auto_status' => 'required|in:1,0',
            'auto_at' => 'required',
        ];
        $messages = [
            'auto_status.required' => '是否开启状态不能为空',
            'auto_at.required' => '自动批准时间不能为空',
            'auto_status.in'   => '是否开启状态值必须为0或1',
        ];

        Validator::make($request->input(), $rules, $messages)->validate();

        //mem_type:1、默认导购  2、默认合伙人
        if(!isset($requestParams['mem_type']) || !in_array($requestParams['mem_type'], [1, 2])){
            $requestParams['mem_type'] = 1;
        }

        $params = [
            'enterprise_id'  => $this->enterpriseId,
            'auto_status'        => $requestParams['auto_status'],
            'auto_at'           => $requestParams['auto_at'],
            'mem_type'  => $requestParams['mem_type'],
        ];
        (new Withdrawal())->autoAdd($params);

        return success();
    }
}