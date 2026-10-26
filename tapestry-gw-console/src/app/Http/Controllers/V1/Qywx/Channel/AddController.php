<?php

namespace App\Http\Controllers\V1\Qywx\Channel;

use CrsSdk\Work\Channel;
use App\Http\Controllers\V1\BaseController;
use Illuminate\Validation\ValidationException;

class AddController extends BaseController
{
    /**
     * 新增渠道
     * @return mixed
     * @throws ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'channel' => 'required|max:20'
        ], [
            'channel.required' => '请输入渠道名称',
            'channel.max'      => '渠道名称最多输入20字符'
        ]);

        $input                  = $this->request->only(['channel']);
        $input['enterprise_id'] = $this->enterpriseId;
        $input['admin_id']      = $this->enterpriseAdminId;
        $input['admin_name']    = $this->enterpriseAdmin['username'];

        (new Channel())->addChannel($input);

        return success();
    }
}
