<?php


namespace App\Http\Controllers\Service\V1\Partner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;
use App\Services\PartnerMemberService;

class checkPartneStatusByMobileController extends BaseController
{

    public function index()
    {
        // 接值
        $input = $this->request->only([
            'mobile',
            'enterprise_id',
        ]);
        $rules = [
            'enterprise_id' => 'required',
            'mobile' => 'required'
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'mobile.required' => '合伙人手机号不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        $sql = PartnerMember::query()
            ->leftJoin('partner_enterprise_map as r', 'r.partner_id', '=', 'partner_member.id')
            ->where('partner_member.partner_mobile', '=', $input['mobile'])
            ->where('r.enterprise_id', '=', $input['enterprise_id'])
            ->select(['partner_member.*', 'r.worker_id', 'r.forbidden_reason', 'r.partner_status', 'r.user_id']);
        $info = $sql->first()->toArray();

        return success($info);
    }

}