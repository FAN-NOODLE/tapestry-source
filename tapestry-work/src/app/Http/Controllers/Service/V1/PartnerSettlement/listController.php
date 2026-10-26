<?php


namespace App\Http\Controllers\Service\V1\PartnerSettlement;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\PartnerSettlementService;

class listController extends BaseController
{
    public function index(){
        // 接值
        $input = $this->request->only([
            'page',
            'size',
            'enterprise_id',
            'search_type',
            'search_value',
            'status_type',
            'status_value',
            'worker_ids',
            'order_by_type',
            'orderby'
        ]);
        $rules = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);

        $data = (new PartnerSettlementService())->SettlementList($input);

        return success($data);
    }

}