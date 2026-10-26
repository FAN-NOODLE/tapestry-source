<?php


namespace App\Http\Controllers\Service\V1\Partner;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Services\PartnerMemberService;

class SwitchController extends BaseController
{
    public function index()
    {
        // 接值
        $input = $this->request->only([
            'id',
            'enterprise_id',
            'switch',
            'forbidden_reason'
        ]);
        $rules = [
            'id' => 'required',
            'enterprise_id' => 'required',
            'switch' => 'required|in:1,2'
        ];
        $messages = [
            'id.required' => '合伙人id不能为空',
            'switch.required' => '启用禁用状态不能为空',
            'switch.in' => '[参数]启用禁用状态值不合法',
            'enterprise_id.required' => '企业id不能为空',
        ];

        $this->validate($this->request, $rules, $messages);

        $partnerInfo = PartnerEnterpriseMap::query()
            ->where(['partner_id' => $input['id'], 'enterprise_id' => $input['enterprise_id']])
            ->first();

        if (empty($partnerInfo)) {
            throw new Error(700000, '合伙人不存在');
        }

        $input['switch'] = $input['switch'] == 2 ? 0 : 1;
        if ($partnerInfo->partner_status == $input['switch']) {
            return success();
        }
        $input['forbidden_code'] = $input['forbidden_code'] ?? 0;
        $input['partner_enterprise_id'] = $partnerInfo->id;

        (new PartnerMemberService())->Switch($input);

        return success();
    }

}