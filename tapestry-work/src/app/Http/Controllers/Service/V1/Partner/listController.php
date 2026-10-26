<?php


namespace App\Http\Controllers\Service\V1\Partner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\PartnerMemberService;

class listController extends BaseController
{

    public function index()
    {
        // 接值
        $input = $this->request->only([
            'page',
            'size',
            'enterprise_id',
            'search_type',
            'search_value',
            'partner_status',
            'commission',
            'commission_min',
            'commission_max',
            'orderby',
            'order_by_type',
            'worker_ids'
        ]);
        $rules = [
            'enterprise_id' => 'required',
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $data = (new PartnerMemberService())->partnerList($input);

        return success($data);
    }

}