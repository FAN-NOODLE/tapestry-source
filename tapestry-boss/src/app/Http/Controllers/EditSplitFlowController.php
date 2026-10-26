<?php

namespace App\Http\Controllers;

use App\Http\Middleware\SignMiddleware;
use App\Models\Enterprise;

class EditSplitFlowController extends Controller
{
    /**
     * 编辑是否切流配置接口
     *
     * @return array|string
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'enterprise_id' => 'required|integer|min:0',
            'is_split_flow' => 'present|integer|in:0,1',
        ]);
        $sign = $this->request->input('sign');
        $epID = (int)$this->request->input('enterprise_id');
        $isSplitFlow = (int)$this->request->input('is_split_flow');

        if ($sign !== SignMiddleware::DEBUG_SIGN) {
            return '你想干嘛!';
        }

        $enterprise = Enterprise::query()->where('id', '=', $epID)->first();
        if ($enterprise === null) {
            return '记录不存在';
        }
        $enterprise->update([
            'is_split_flow' => $isSplitFlow
        ]);

        return $enterprise->toArray();
    }
}
