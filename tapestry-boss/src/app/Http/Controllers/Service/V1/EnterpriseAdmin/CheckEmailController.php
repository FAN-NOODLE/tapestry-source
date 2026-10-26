<?php

namespace App\Http\Controllers\Service\V1\EnterpriseAdmin;

use App\Models\EnterpriseAdmin;
use App\Http\Controllers\Service\V1\BaseController;

class CheckEmailController extends BaseController
{
    /**
     * 检查邮箱是否唯一接口
     *
     * @return \Illuminate\Http\JsonResponse
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'email'                => 'required|string|email',// 邮箱
            'filter_enterprise_id' => 'integer|min:0',// 过滤的企业ID 编辑可使用
        ]);

        $input = $this->request->all();

        $exists = EnterpriseAdmin::query()
            ->where('email', '=', $input['email'])
            ->when(isset($input['filter_enterprise_id']), function ($query) use ($input) {
                return $query->where('enterprise_id', '!=', $input['filter_enterprise_id']);
            })
            ->exists();

        return success([
            'status' => $exists
        ]);
    }
}
