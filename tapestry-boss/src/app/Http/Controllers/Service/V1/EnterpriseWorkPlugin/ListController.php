<?php

namespace App\Http\Controllers\Service\V1\EnterpriseWorkPlugin;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\EnterpriseWorkPlugin;

class ListController extends BaseController
{
    /**
     * 企业微信下的应用列表
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required|integer'
        ]);

        $input = $this->request->all();

        $data = EnterpriseWorkPlugin::query()->where('enterprise_id', '=', $input['enterprise_id'])->get()->toArray();

        return success($data);
    }
}
