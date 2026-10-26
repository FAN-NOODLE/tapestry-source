<?php


namespace App\Http\Controllers\Service\V1\Partner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Services\PartnerMemberService;

class DetailController extends BaseController
{

    public function index()
    {
        // 接值
        $input = $this->request->only([
            'page',
            'size',
            'id',
            'enterprise_id',
        ]);
        $rules = [
            'enterprise_id' => 'required',
            'id' => 'required'
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'id.required' => '合伙人id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        switchAppHost(['enterprise_id' => $input['enterprise_id']]);
        $data = (new PartnerMemberService())->partnerDetail($input);

        return success($data);
    }

}