<?php

namespace App\Http\Controllers\Service\V1\Channel;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;

class AddController extends BaseController
{
    /**
     * 新增渠道
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @author: DLong
     */
    public function index()
    {
        $this->validate($this->request, [
            'channel'       => 'required|max:20',
            'enterprise_id' => 'required|integer',
            'admin_id'      => 'required|integer',
            'admin_name'    => 'required',
        ], [
            'channel.required'       => '[参数]channel必填',
            'channel.max'            => '[参数]channel最多20个字符',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
            'admin_id.required'      => '[参数]admin_id必填',
            'admin_id.integer'       => '[参数]admin_id需为integer',
            'admin_name.required'    => '[参数]admin_name必填'
        ]);

        $input        = $this->request->only(['enterprise_id', 'source_type', 'channel', 'admin_id', 'admin_name', 'return_id']);
        $existChannel = Channel::query()
            ->where([
                'enterprise_id' => $input['enterprise_id'],
                'channel'       => $input['channel']
            ])->first(['id']);

        if ($existChannel) {
            //if (!empty($input['return_id'])) {
                return success(['data' => ['id' => $existChannel->id]]);
            //} else {
              //  return notice(610000);
            //}
        }

        unset($input['return_id']);

        $data = Channel::query()->create($input);

        return success(['data' => ['id' => $data->id ?? 0]]);
    }
}
