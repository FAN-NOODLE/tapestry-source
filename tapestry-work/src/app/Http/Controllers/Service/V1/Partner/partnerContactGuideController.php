<?php


namespace App\Http\Controllers\Service\V1\Partner;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerMember;
use App\Services\Message\MessageService;
use CrsSdk\Boss\Enterprise;

class partnerContactGuideController extends BaseController
{

    const DISABLE = 1;//禁用联系

    const CONTACT = 2;//我的导购

    public function index()
    {
        // 接值
        $input = $this->request->only([
            'enterprise_id',
            'partner_id',
            'title',
            'content',
            'page',
            'contactType'
        ]);
        $rules = [
            'enterprise_id' => 'required',
            'partner_id' => 'required',
            'title' => 'required',
            'contactType' => 'required',
            'page' => 'required'
        ];
        $messages = [
            'enterprise_id.required' => '企业id不能为空',
            'partner_id.required' => '合伙人id不能为空',
            'title.required' => '消息title不能为空',
            'contactType.required' => '消息类型不能为空',
            'page.required' => '跳转地址不能为空'
        ];

        $this->validate($this->request, $rules, $messages);

        $enterpriseData = (new Enterprise())->detail(
            ['enterprise_id' => $input['enterprise_id']]
        );

        $qyDomain = $enterpriseData['data']['domain'] ?? '';
        $input['qy_domain'] = $qyDomain;
        $sql = PartnerMember::query()
            ->select(['partner_member.partner_name', 'partner_member.partner_mobile', 'r.*'])
            ->leftJoin('partner_enterprise_map as r', 'r.partner_id', '=', 'partner_member.id')
            ->where(['partner_member.id' => $input['partner_id'], 'r.enterprise_id' => $input['enterprise_id']]);
        $info = $sql->first()->toArray();

        $content = '';
        if ($input['contactType'] == self::DISABLE) {
            $content = [
                [
                    'key' => '合伙人微信昵称',
                    'value' => $info['partner_name'] ?? '',
                ],
                [
                    'key' => '合伙人电话',
                    'value' => $info['mobile_show'] ?? '',
                ],
                [
                    'key' => '禁用原因',
                    'value' => $info['forbidden_reason'],
                ]
            ];
        }

        if ($input['contactType'] == self::CONTACT) {
            $content = [
                [
                    'key' => '合伙人微信昵称',
                    'value' => $info['partner_name'] ?? '',
                ],
                [
                    'key' => '合伙人电话',
                    'value' => $info['mobile_show'] ?? '',
                ],
                [
                    'key' => '备注',
                    'value' => '请回拨电话或添加合伙人为您的联系人进行企业微信会话。',
                ]
            ];
        }

        $title = $input['title'];
        $page = $input['page'];
        (new MessageService())->sendMiniMessage(
            $input['enterprise_id'],
            $content,
            $page,
            $title,
            [$info['worker_uid']],
            $input['qy_domain']
        );

        return success([]);
    }

}