<?php

namespace App\Http\Controllers\V1\Qywx\Channel;

use CrsSdk\Work\Channel;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class EditController extends BaseController
{
    /**
     * 编辑渠道
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'id'      => 'required|integer',
            'channel' => 'required|max:20'
        ], [
            'id.required'      => '渠道id必填',
            'id.integer'       => '渠道id需为integer',
            'channel.required' => '请输入渠道名称',
            'channel.max'      => '渠道名称最多输入20字符'
        ]);
        $input = $this->request->only(['id', 'channel']);

        $input['enterprise_id'] = $this->enterpriseId;

        (new Channel())->editChannel($input);

        return success();
    }
}
