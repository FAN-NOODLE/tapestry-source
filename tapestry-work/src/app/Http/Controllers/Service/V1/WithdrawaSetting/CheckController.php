<?php


namespace App\Http\Controllers\Service\V1\WithdrawaSetting;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\SettlementService;


class CheckController extends BaseController
{
    /**
     *
     * 根据项目编码获取项目信息
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(){
        $rules = [
            'transfer_no' => 'required',
        ];
        $messages = [
            'transfer_no.required' => '转账项目编号不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        $info = (new SettlementService())->getProjectByCode($this->request['transfer_no']);

        return success($info);
    }

}