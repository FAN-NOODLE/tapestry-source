<?php

namespace App\Http\Controllers\Service\V1\StaffShare;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\StaffShare;
use App\Services\StaffShareService;

class GetStaffShareUserIdsController extends BaseController
{
    public function index()
    {
        $rules = [
            'enterprise_id' => 'required|integer',
            'worker_id' => 'integer',
            'worker_user_id' => 'integer',
            'user_ids' => 'array',
            'ascription' => 'required|in:' . StaffShare::ACHIEVEMENT_ASCRIPTION_SELF . ',' . StaffShare::ACHIEVEMENT_ASCRIPTION_OTHER . ',' . StaffShare::ACHIEVEMENT_ASCRIPTION_NONE,
        ];
        $message = [
            'enterprise_id.required' => 'enterprise_id不能为空',
            'enterprise_id.integer' => 'enterprise_id必须为int类型',
            'worker_id.integer' => 'worker_id必须为int类型',
            'worker_user_id.integer' => 'worker_id必须为int类型',
            'user_ids.array' => 'user_ids必须为数组',
            'ascription.required' => 'ascription参数必填',
            'ascription.in' => 'ascription参数错误',
        ];
        $input = $this->validate($this->request, $rules, $message);
        $data = (new StaffShareService())->getUserIds($input);

        return success($data);
    }
}
