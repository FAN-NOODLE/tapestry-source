<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Exceptions\Error;
use App\Http\Controllers\App\V1\BaseController;
use App\Models\GuideUserRelation;
use Illuminate\Http\Request;

class MarkUserController extends BaseController
{
    /**
     * 标记是否复制手机号
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $userId = $request->input('user_id', 0);

        if (empty($userId)) {
            throw new Error(10000, '参数user_id不能为空');
        }

        // 是否已复制手机号（默认1）:0、否, 1、是
        $isCopyMobile = $request->input('is_copy_mobile', 1);

        if (!in_array($isCopyMobile, [0, 1])) {
            throw new Error(20000, '参数is_copy_mobile不在范围');
        }

        $num = GuideUserRelation::query()
            ->where([
                'mem_uid' => $this->memUid,
                'user_id' => $userId,
            ])
            ->update([
                'is_copy_mobile' => $isCopyMobile
            ]);

        if (0 == $num) {
            throw new Error(30000, '标记员工复制手机号失败');
        }

        return success(['data' => []]);
    }
}
