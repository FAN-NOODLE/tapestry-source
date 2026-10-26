<?php

namespace App\Http\Controllers\Service\V1\Channel;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;
use App\Models\ContactQrCode;

class EditController extends BaseController
{
    /**
     * 编辑渠道
     * {domain}/service/v1/channel/edit
     *
     * @return array|void
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'id'            => 'required|integer',
            'enterprise_id' => 'required|integer',
            'channel'       => 'required|max:20'
        ], [
            'id.required'            => '[参数]渠道id必填',
            'id.integer'             => '[参数]渠道id需为integer',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
            'channel.required'       => '[参数]channel必填',
            'channel.max'            => '[参数]channel最多20个字符'
        ]);

        $input = $this->request->only(['id', 'enterprise_id', 'channel']);

        $existChannel = Channel::query()
            ->where('id', '<>', $input['id'])
            ->where([
                'enterprise_id' => $input['enterprise_id'],
                'channel'       => $input['channel']
            ])->exists();

        if ($existChannel) {
            notice(610000);
        }

        // 原渠道名称
        $oldChannel = Channel::query()->where('id', $input['id'])->value('channel');

        Channel::query()
            ->where([
                'id'            => $input['id'],
                'enterprise_id' => $input['enterprise_id']
            ])->update([
                'channel' => $input['channel']
            ]);

        ContactQrCode::query()
            ->where([
                'enterprise_id' => $input['enterprise_id'],
                'channel'       => $oldChannel
            ])->update([
                'channel' => $input['channel']
            ]);

        return success(['data' => []]);
    }
}
