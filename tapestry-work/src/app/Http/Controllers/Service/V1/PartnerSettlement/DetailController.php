<?php


namespace App\Http\Controllers\Service\V1\PartnerSettlement;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\PartnerSettlementService;

class DetailController extends BaseController
{
    public function index()
    {
        // 接值
        $input = $this->request->only([
            'page',
            'size',
            'enterprise_id',
            'id'
        ]);
        $rules = [
            'enterprise_id' => 'required',
            'id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'id.required' => '结算单号不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);

        $data = (new PartnerSettlementService())->Detail($input);

        return success($data);

    }

}